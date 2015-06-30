package backends

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/oliveagle/hickwall/backends/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"strings"
	"time"
)

var (
	_ = time.Now()
	_ = fmt.Sprintf("")
)

type kafkaBackend struct {
	name    string
	closing chan chan error             // for Close
	updates chan newcore.MultiDataPoint // for receive updates

	// kafka backend specific attributes
	conf     *config.Transport_kafka
	kconf    *sarama.Config
	producer sarama.AsyncProducer
}

func MustNewKafkaBackend(name string, bconf *config.Transport_kafka) *kafkaBackend {
	logging.Info("MustNewKafkaBackend: %+v", bconf)
	_kconf := sarama.NewConfig()

	_kconf.Net.DialTimeout = newcore.Interval(bconf.Dail_timeout).MustDuration(time.Second * 5)
	_kconf.Net.WriteTimeout = newcore.Interval(bconf.Write_timeout).MustDuration(time.Second * 1)
	_kconf.Net.ReadTimeout = time.Second * 10
	_kconf.Net.KeepAlive = newcore.Interval(bconf.Keepalive).MustDuration(time.Second * 30)

	if bconf.Ack_timeout_ms <= 0 {
		_kconf.Producer.Timeout = time.Millisecond * 100
	} else {
		_kconf.Producer.Timeout = time.Millisecond * time.Duration(bconf.Ack_timeout_ms)
	}

	if bconf.Flush_frequency_ms <= 0 {
		_kconf.Producer.Flush.Frequency = time.Millisecond * 100
	} else {
		_kconf.Producer.Flush.Frequency = time.Millisecond * time.Duration(bconf.Flush_frequency_ms)
	}

	cc := strings.ToLower(bconf.Compression_codec)
	switch {
	case cc == "none":
		_kconf.Producer.Compression = sarama.CompressionNone
	case cc == "gzip":
		_kconf.Producer.Compression = sarama.CompressionGZIP // Compress messages
	case cc == "snappy":
		_kconf.Producer.Compression = sarama.CompressionSnappy // Compress messages
	default:
		_kconf.Producer.Compression = sarama.CompressionNone
	}

	ra := strings.ToLower(bconf.Required_acks)
	switch {
	case ra == "no_response":
		_kconf.Producer.RequiredAcks = sarama.NoResponse
	case ra == "wait_for_local":
		_kconf.Producer.RequiredAcks = sarama.WaitForLocal
	case ra == "wait_for_all":
		_kconf.Producer.RequiredAcks = sarama.WaitForAll
	default:
		_kconf.Producer.RequiredAcks = sarama.NoResponse
	}

	logging.Debugf("kafka conf: %+v", _kconf)

	s := &kafkaBackend{
		name:    name,
		closing: make(chan chan error),
		updates: make(chan newcore.MultiDataPoint),
		conf:    bconf,  // backend config
		kconf:   _kconf, // sarama config
	}
	go s.loop()
	return s
}

func (b *kafkaBackend) connect() error {
	producer, err := sarama.NewAsyncProducer(b.conf.Broker_list, b.kconf)
	if err != nil {
		logging.Errorf("failed to start producer: %v, %v", err, b.conf.Broker_list)
		return fmt.Errorf("failed to start producer: %v, %v", err, b.conf.Broker_list)
	}

	go func() {
		logging.Debug("consuming from producer.Errors()")
		for err := range producer.Errors() {
			logging.Errorf("producer error: %v", err)
		}
		logging.Debug("producer.Errors() closed")
	}()

	logging.Infof("created new producer: %v", b.conf.Broker_list)

	// save producer reference
	b.producer = producer
	return nil
}

func (b *kafkaBackend) loop() {
	var (
		startConsuming    <-chan newcore.MultiDataPoint
		try_connect_first chan bool
		try_connect_tick  <-chan time.Time
	)
	startConsuming = b.updates
	logging.Info("kafkaBackend.loop started")

	for {
		if b.producer == nil && try_connect_first == nil && try_connect_tick == nil {
			startConsuming = nil // disable consuming

			try_connect_first = make(chan bool)
			logging.Debug("trying to connect to kafka first time.")

			// trying to connect to kafka first time
			go func() {
				err := b.connect()
				if b.producer != nil && err == nil {
					logging.Debugf("connect kafka first time OK: %v", b.producer)
					try_connect_first <- true
				} else {
					logging.Criticalf("connect to kafka failed %s", err)
					try_connect_first <- false
				}
			}()
		}
		if startConsuming != nil {
			logging.Trace("kafkaBackend consuming started")
		}

		select {
		case md := <-startConsuming:
			for _, p := range md {
				b.producer.Input() <- &sarama.ProducerMessage{
					Topic: b.conf.Topic_id,
					Key:   sarama.StringEncoder(p.Metric),
					Value: &p,
				}
				logging.Tracef(" -> point: %+v", p)
			}
			logging.Debugf("kafkaBackend consuming finished: count: %d", len(md))
		case connected := <-try_connect_first:
			try_connect_first = nil // disable this branch
			if !connected {
				// failed open it the first time,
				// then we try to open file with time interval, until connected successfully.
				logging.Critical("connect first time failed, try to connect with interval of 1s")
				try_connect_tick = time.Tick(time.Second * 1)
			} else {
				logging.Debug("kafka connected the first time.")
				startConsuming = b.updates
			}
		case <-try_connect_tick:
			// try to connect with interval
			err := b.connect()
			if b.producer != nil && err == nil {
				// finally connected.
				try_connect_tick = nil
				startConsuming = b.updates
			} else {
				logging.Criticalf("kafka backend trying to connect but failed: %s", err)
			}
		case errc := <-b.closing:
			logging.Info("kafaBackend.loop closing")
			startConsuming = nil // stop comsuming
			errc <- nil
			close(b.updates)
			logging.Info("kafaBackend.loop closed")
			return
		}
	}
}

func (b *kafkaBackend) Updates() chan<- newcore.MultiDataPoint {
	return b.updates
}

func (b *kafkaBackend) Close() error {
	errc := make(chan error)
	b.closing <- errc
	if b.producer != nil {
		b.producer.Close()
	}
	return <-errc
}

func (b *kafkaBackend) Name() string {
	return b.name
}
