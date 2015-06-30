package collectors

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	. "github.com/oliveagle/hickwall/collectors/config"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	_ = fmt.Sprint("")
)

type kafka_sub_state struct {
	mu          sync.RWMutex
	filename    string
	state       []kafka_sub_state_line
	last_flush  time.Time
	last_update time.Time
}

type kafka_sub_state_line struct {
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
}

func newKafkaSubState(filename string) *kafka_sub_state {
	return &kafka_sub_state{
		filename:    filename,
		last_update: time.Now(),
	}
}

func (s *kafka_sub_state) Length() int {
	return len(s.state)
}

func (s *kafka_sub_state) State() []kafka_sub_state_line {
	return s.state[:]
}

func (s *kafka_sub_state) Clear() {
	s.state = nil
}

func (s *kafka_sub_state) Update(topic string, part int32, offset int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Length() > 0 {
		old_line_idx := -1 // should not start from 0
		for idx, line := range s.state {
			if line.Topic == topic && line.Partition == part {
				old_line_idx = idx
				if line.Offset == offset {
					return nil // nothing changed
				}
				break
			}
		}
		if old_line_idx >= 0 && old_line_idx < s.Length() {
			s.state = append(s.state[:old_line_idx], s.state[old_line_idx+1:]...)
		}
	}
	s.state = append(s.state, kafka_sub_state_line{topic, part, offset})
	s.last_update = time.Now()
	return nil
}

func (s *kafka_sub_state) Partitions(topic string) []int32 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var parts []int32
	for _, line := range s.state {
		if line.Topic == topic {
			parts = append(parts, line.Partition)
		}
	}
	return parts
}

func (s *kafka_sub_state) Offset(topic string, part int32) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, line := range s.state {
		if line.Topic == topic && line.Partition == part {
			return line.Offset
		}
	}
	return -1
}

func (s *kafka_sub_state) Changed() bool {
	return s.last_update.After(s.last_flush)
}

func (s *kafka_sub_state) Save() (err error) {
	file, err := os.OpenFile(s.filename, os.O_CREATE|os.O_WRONLY|os.O_SYNC|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	err = enc.Encode(s.state)
	if err != nil {
		return err
	}
	s.last_flush = time.Now()
	return nil
}

func (s *kafka_sub_state) Load() (err error) {
	var state []kafka_sub_state_line
	file, err := os.OpenFile(s.filename, os.O_RDONLY, 0644)
	if err != nil {
		//TODO: doesn't exists is not an error
		return err
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	err = dec.Decode(&state)
	if err != nil {
		return err
	}
	s.state = state[:]
	s.last_update = time.Now()
	return nil
}

type kafka_subscription struct {
	// kafka_subscription specific attributes
	opts      Config_KafkaSubscription
	closing   chan chan error
	updates   chan newcore.MultiDataPoint
	master    sarama.Consumer
	consumers []sarama.PartitionConsumer

	flush_interval time.Duration
	max_batch_size int

	state_file_name string
	state           *kafka_sub_state
}

func NewKafkaSubscription(opts Config_KafkaSubscription) (*kafka_subscription, error) {
	var max_batch_size = 100
	var state_file_name string

	if opts.Max_batch_size > 0 {
		max_batch_size = opts.Max_batch_size
	}

	if opts.Name == "" {
		return nil, fmt.Errorf("name should not be empty")
	}
	if opts.Topic == "" {
		return nil, fmt.Errorf("topic should not be empty")
	}
	state_file_name = filepath.Join(config.SHARED_DIR, fmt.Sprintf("%s_%s.state", opts.Name, opts.Topic))

	c := &kafka_subscription{
		opts:            opts,
		updates:         make(chan newcore.MultiDataPoint), // for Updates
		closing:         make(chan chan error),             // for Close
		flush_interval:  newcore.Interval(opts.Flush_interval).MustDuration(time.Millisecond * 100),
		max_batch_size:  max_batch_size,
		state_file_name: state_file_name,
		state:           newKafkaSubState(state_file_name),
	}

	go c.loop()
	return c, nil
}

func (c *kafka_subscription) Name() string {
	return "kafka_subscription" // TODO: we should remove Name()
}

func (c *kafka_subscription) Close() error {
	errc := make(chan error)
	c.closing <- errc
	return <-errc
}

func (c *kafka_subscription) connect() error {
	logging.Info("connect")

	sconfig := sarama.NewConfig()
	logging.Debugf("broker list: %v", c.opts.Broker_list)

	master, err := sarama.NewConsumer(c.opts.Broker_list, sconfig)
	if err != nil {
		return fmt.Errorf("Cannot connect to kafka: %v", err)
	}
	c.master = master
	return nil
}

func (c *kafka_subscription) consume() (<-chan newcore.MultiDataPoint, error) {
	logging.Info("consume")

	var out = make(chan newcore.MultiDataPoint)
	var err error
	var consumers []sarama.PartitionConsumer
	if c.master == nil {
		err = c.connect()
		if err != nil {
			return nil, err
		}
	}

	for _, c := range c.consumers {
		c.Close()
	}
	c.consumers = nil

	partitions, err := c.master.Partitions(c.opts.Topic)
	if err != nil {
		return nil, fmt.Errorf("Cannot get partitions: %v", err)
	}
	logging.Infof("partitions: %v", partitions)

	err = c.state.Load()
	if err != nil {
		logging.Errorf("failed to load kafka state: %v", err)
	} else {
		logging.Infof("state: %+v", c.state.State())
	}

	flush_offset := true

	for _, part := range partitions {
		offset := int64(0)
		if c.state.Length() > 0 {
			offset = c.state.Offset(c.opts.Topic, part)
			if offset < 0 {
				offset = 0
			}
		}
		consumer, err := c.master.ConsumePartition(c.opts.Topic, part, offset)
		if err != nil {
			logging.Criticalf("Cannot consumer partition: %d, %v", part, err)
			return nil, fmt.Errorf("Cannot consumer partition: %d, %v", part, err)
		}
		logging.Infof("created consumer: %v", consumer)

		consumers = append(consumers, consumer)

		go func(flush_offset bool, topic string, part int32, out chan newcore.MultiDataPoint, consumer sarama.PartitionConsumer) {
			logging.Infof("start goroutine to consume: part: %d,  %v", part, &consumer)

			var items newcore.MultiDataPoint
			var flush_tick = time.Tick(c.flush_interval)
			var _out chan newcore.MultiDataPoint
			var startConsume <-chan *sarama.ConsumerMessage
			var flushing bool
			var offset int64

			for {
				if (flushing == true && len(items) > 0) || len(items) >= c.max_batch_size {
					_out = out         // enable output branch
					startConsume = nil // disable consuming branch
				} else if len(items) < c.max_batch_size {
					startConsume = consumer.Messages() // enable consuming branch
					_out = nil                         // disable output branch
				}

				select {
				case message := <-startConsume:
					offset = message.Offset
					dp, err := newcore.NewDPFromJson(message.Value)
					if err != nil {
						logging.Tracef("[ERROR]failed to parse datapoint: %v", err)
					}
					logging.Tracef("kafka dp --> %v", dp)
					items = append(items, dp)
				case <-flush_tick:
					flushing = true
					// every part consumer will record offset with interval
					c.state.Update(topic, part, offset)

					// only 1 goroutine will save state to disk
					if flush_offset == true && c.state.Changed() == true {
						logging.Tracef("flusing to disk: part: %d, offset: %d", part, offset)
						c.state.Save()
					}
				case _out <- items:
					items = nil                        // clear items
					_out = nil                         // disable output branch
					startConsume = consumer.Messages() // enable consuming branch
					flushing = false                   // disable flusing
				case err := <-consumer.Errors():
					logging.Infof("consumer.Errors: part:%d,  %v", part, err)
				}
			}
		}(flush_offset, c.opts.Topic, part, out, consumer)

		flush_offset = false // only 1st goroutine is responsible for flushing state back into disk
	}
	c.consumers = consumers
	return out, nil
}

func (c *kafka_subscription) loop() {
	logging.Info("loop started")

	var items newcore.MultiDataPoint
	// var tick = time.Tick(c.interval)
	var output chan newcore.MultiDataPoint
	var try_open_consumer_tick <-chan time.Time
	// var tick_reconsume chan time.Time
	var input <-chan newcore.MultiDataPoint
	// var err error

	for {
		if input == nil && try_open_consumer_tick == nil {
			logging.Info("input == nil, try_open_consumer_tick == nil")
			try_open_consumer_tick = time.After(0)
		}

		select {
		case <-try_open_consumer_tick:
			logging.Info("try_open_consumer_tick")
			_input, err := c.consume()
			if _input != nil && err == nil {
				try_open_consumer_tick = nil
			} else {
				logging.Errorf("failed to create consumers: %v", err)
				try_open_consumer_tick = time.After(time.Second)
			}
			input = _input
		case md := <-input:
			items = md[:]
			output = c.updates
			// fmt.Println
		case output <- items:
			items = nil
			output = nil
		case errc := <-c.closing:
			// clean up collector resource.
			output = nil
			close(c.updates)
			errc <- nil
			return
		}
	}
}

func (c *kafka_subscription) Updates() <-chan newcore.MultiDataPoint {
	return c.updates
}
