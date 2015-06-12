package collectors

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/oliveagle/hickwall/collectors/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"time"
)

var (
	_ = fmt.Sprint("")
)

type kafka_subscription struct {
	name string // collector name

	// kafka_subscription specific attributes
	opts      *config.KafkaSubscription
	closing   chan chan error
	updates   chan newcore.MultiDataPoint
	master    sarama.Consumer
	consumers []sarama.PartitionConsumer

	flush_interval time.Duration
	max_batch_size int
}

func newKafkaSubscription(name string, opts config.KafkaSubscription) newcore.Subscription {
	var max_batch_size = 100
	if opts.Max_batch_size > 0 {
		max_batch_size = opts.Max_batch_size
	}
	c := &kafka_subscription{
		name:           name,
		opts:           &opts,
		updates:        make(chan newcore.MultiDataPoint), // for Updates
		closing:        make(chan chan error),             // for Close
		flush_interval: newcore.Interval(opts.Flush_interval).MustDuration(time.Millisecond * 100),
		max_batch_size: max_batch_size,
	}
	go c.loop()
	return c
}

func (c *kafka_subscription) Name() string {
	return c.name
}

func (c *kafka_subscription) Close() error {
	errc := make(chan error)
	c.closing <- errc
	return <-errc
}

func (c *kafka_subscription) connect() error {
	// var consumers []sarama.PartitionConsumer
	logging.Info("connect")

	sconfig := sarama.NewConfig()
	// sconfig.Consumer.Return.Errors = true // Handle errors manually instead of letting Sarama log them.
	logging.Debugf("broker list: %v", c.opts.Broker_list)

	master, err := sarama.NewConsumer(c.opts.Broker_list, sconfig)
	if err != nil {
		return fmt.Errorf("Cannot connect to kafka: %v", err)
	}
	// c.master = master
	c.master = master
	// c.consumers = consumers[:]
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

	partitions, err := c.master.Partitions(c.opts.Topic_id)
	if err != nil {
		return nil, fmt.Errorf("Cannot get partitions: %v", err)
	}
	logging.Infof("partitions: %v", partitions)

	for _, part := range partitions {
		// consumer, err := c.master.ConsumePartition(c.opts.Topic_id, part, sarama.OffsetOldest)
		// consumer, err := c.master.ConsumePartition(c.opts.Topic_id, part, sarama.OffsetNewest)
		consumer, err := c.master.ConsumePartition(c.opts.Topic_id, part, 0)
		// consumer, err := c.master.ConsumePartition(c.opts.Topic_id, part, sarama.OffsetOldest)
		if err != nil {
			logging.Criticalf("Cannot consumer partition: %d, %v", part, err)
			return nil, fmt.Errorf("Cannot consumer partition: %d, %v", part, err)
		}
		logging.Infof("created consumer: %v", consumer)

		consumers = append(consumers, consumer)

		go func(part int32, out chan newcore.MultiDataPoint, consumer sarama.PartitionConsumer) {
			logging.Infof("start goroutine to consume: part: %d,  %v", part, &consumer)

			var items newcore.MultiDataPoint
			// var flush_tick = time.Tick(time.Millisecond * 100)
			var flush_tick = time.Tick(c.flush_interval)
			var _out chan newcore.MultiDataPoint
			// var max_batch_size = 100
			var startConsume <-chan *sarama.ConsumerMessage
			var flushing bool

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
					// logging.Infof("Consumed message with offset %d", message.Offset)
					// logging.Infof("message.Value: ->%s<-", message.Value)
					dp, err := newcore.NewDPFromJson(message.Value)
					if err != nil {
						logging.Tracef("[ERROR]failed to parse datapoint: %v", err)
					}
					// logging.Tracef("kafka dp --> %v", dp)
					items = append(items, dp)
				case <-flush_tick:
					// logging.Debugf("Flush Tick")
					flushing = true
				case _out <- items:
					items = nil                        // clear items
					_out = nil                         // disable output branch
					startConsume = consumer.Messages() // enable consuming branch
					flushing = false                   // disable flusing
				case err := <-consumer.Errors():
					logging.Infof("consumer.Errors: part:%d,  %v", part, err)
				}
			}
		}(part, out, consumer)
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
