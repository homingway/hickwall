package config

type Transport_kafka struct {
	Broker_list        []string `json:"broker_list"`        // ["localhost:xxx", "remote:xxx"]
	Topic_id           string   `json:"topic_id"`           //
	Compression_codec  string   `json:"compression_codec"`  // none, gzip or snappy
	Ack_timeout_ms     int      `json:"ack_timeout_ms"`     // milliseconds
	Required_acks      string   `json:"required_acks"`      // no_response, wait_for_local, wait_for_all
	Flush_frequency_ms int      `json:"flush_frequency_ms"` // milliseconds
	Write_timeout      string   `json:"write_timeout"`      // string, 100ms, 1s, default 1s
	Dail_timeout       string   `json:"dail_timeout"`       // string, 100ms, 1s, default 5s
	Keepalive          string   `json:"keepalive"`          // string, 100ms, 1s, 0 to disable it. default 30s
}

// type kafka struct {
// 	// Net is the namespace for network-level properties used by the Broker, and shared by the Client/Producer/Consumer.
// 	Net struct {
// 		MaxOpenRequests int // How many outstanding requests a connection is allowed to have before sending on it blocks (default 5).

// 		// All three of the below configurations are similar to the `socket.timeout.ms` setting in JVM kafka.
// 		DialTimeout  time.Duration // How long to wait for the initial connection to succeed before timing out and returning an error (default 30s).
// 		ReadTimeout  time.Duration // How long to wait for a response before timing out and returning an error (default 30s).
// 		WriteTimeout time.Duration // How long to wait for a transmit to succeed before timing out and returning an error (default 30s).

// 		// KeepAlive specifies the keep-alive period for an active network connection.
// 		// If zero, keep-alives are disabled. (default is 0: disabled).
// 		KeepAlive time.Duration
// 	}

// 	// Metadata is the namespace for metadata management properties used by the Client, and shared by the Producer/Consumer.
// 	Metadata struct {
// 		Retry struct {
// 			Max     int           // The total number of times to retry a metadata request when the cluster is in the middle of a leader election (default 3).
// 			Backoff time.Duration // How long to wait for leader election to occur before retrying (default 250ms). Similar to the JVM's `retry.backoff.ms`.
// 		}
// 		// How frequently to refresh the cluster metadata in the background. Defaults to 10 minutes.
// 		// Set to 0 to disable. Similar to `topic.metadata.refresh.interval.ms` in the JVM version.
// 		RefreshFrequency time.Duration
// 	}

// 	// Producer is the namespace for configuration related to producing messages, used by the Producer.
// 	Producer struct {
// 		// The maximum permitted size of a message (defaults to 1000000). Should be set equal to or smaller than the broker's `message.max.bytes`.
// 		MaxMessageBytes int
// 		// The level of acknowledgement reliability needed from the broker (defaults to WaitForLocal).
// 		// Equivalent to the `request.required.acks` setting of the JVM producer.
// 		RequiredAcks RequiredAcks
// 		// The maximum duration the broker will wait the receipt of the number of RequiredAcks (defaults to 10 seconds).
// 		// This is only relevant when RequiredAcks is set to WaitForAll or a number > 1. Only supports millisecond resolution,
// 		// nanoseconds will be truncated. Equivalent to the JVM producer's `request.timeout.ms` setting.
// 		Timeout time.Duration
// 		// The type of compression to use on messages (defaults to no compression). Similar to `compression.codec` setting of the JVM producer.
// 		Compression CompressionCodec
// 		// Generates partitioners for choosing the partition to send messages to (defaults to hashing the message key).
// 		// Similar to the `partitioner.class` setting for the JVM producer.
// 		Partitioner PartitionerConstructor

// 		// Return specifies what channels will be populated. If they are set to true, you must read from
// 		// the respective channels to prevent deadlock.
// 		Return struct {
// 			// If enabled, successfully delivered messages will be returned on the Successes channel (default disabled).
// 			Successes bool

// 			// If enabled, messages that failed to deliver will be returned on the Errors channel, including error (default enabled).
// 			Errors bool
// 		}

// 		// The following config options control how often messages are batched up and sent to the broker. By default,
// 		// messages are sent as fast as possible, and all messages received while the current batch is in-flight are placed
// 		// into the subsequent batch.
// 		Flush struct {
// 			Bytes     int           // The best-effort number of bytes needed to trigger a flush. Use the global sarama.MaxRequestSize to set a hard upper limit.
// 			Messages  int           // The best-effort number of messages needed to trigger a flush. Use `MaxMessages` to set a hard upper limit.
// 			Frequency time.Duration // The best-effort frequency of flushes. Equivalent to `queue.buffering.max.ms` setting of JVM producer.
// 			// The maximum number of messages the producer will send in a single broker request.
// 			// Defaults to 0 for unlimited. Similar to `queue.buffering.max.messages` in the JVM producer.
// 			MaxMessages int
// 		}

// 		Retry struct {
// 			// The total number of times to retry sending a message (default 3).
// 			// Similar to the `message.send.max.retries` setting of the JVM producer.
// 			Max int
// 			// How long to wait for the cluster to settle between retries (default 100ms).
// 			// Similar to the `retry.backoff.ms` setting of the JVM producer.
// 			Backoff time.Duration
// 		}
// 	}

// 	// Consumer is the namespace for configuration related to consuming messages, used by the Consumer.
// 	Consumer struct {
// 		Retry struct {
// 			// How long to wait after a failing to read from a partition before trying again (default 2s).
// 			Backoff time.Duration
// 		}

// 		// Fetch is the namespace for controlling how many bytes are retrieved by any given request.
// 		Fetch struct {
// 			// The minimum number of message bytes to fetch in a request - the broker will wait until at least this many are available.
// 			// The default is 1, as 0 causes the consumer to spin when no messages are available. Equivalent to the JVM's `fetch.min.bytes`.
// 			Min int32
// 			// The default number of message bytes to fetch from the broker in each request (default 32768). This should be larger than the
// 			// majority of your messages, or else the consumer will spend a lot of time negotiating sizes and not actually consuming. Similar
// 			// to the JVM's `fetch.message.max.bytes`.
// 			Default int32
// 			// The maximum number of message bytes to fetch from the broker in a single request. Messages larger than this will return
// 			// ErrMessageTooLarge and will not be consumable, so you must be sure this is at least as large as your largest message.
// 			// Defaults to 0 (no limit). Similar to the JVM's `fetch.message.max.bytes`. The global `sarama.MaxResponseSize` still applies.
// 			Max int32
// 		}
// 		// The maximum amount of time the broker will wait for Consumer.Fetch.Min bytes to become available before it
// 		// returns fewer than that anyways. The default is 250ms, since 0 causes the consumer to spin when no events are available.
// 		// 100-500ms is a reasonable range for most cases. Kafka only supports precision up to milliseconds; nanoseconds will be truncated.
// 		// Equivalent to the JVM's `fetch.wait.max.ms`.
// 		MaxWaitTime time.Duration

// 		// Return specifies what channels will be populated. If they are set to true, you must read from
// 		// them to prevent deadlock.
// 		Return struct {
// 			// If enabled, any errors that occured while consuming are returned on the Errors channel (default disabled).
// 			Errors bool
// 		}
// 	}

// 	// A user-provided string sent with every request to the brokers for logging, debugging, and auditing purposes.
// 	// Defaults to "sarama", but you should probably set it to something specific to your application.
// 	ClientID string
// 	// The number of events to buffer in internal and external channels. This permits the producer and consumer to
// 	// continue processing some messages in the background while user code is working, greatly improving throughput.
// 	// Defaults to 256.
// 	ChannelBufferSize int
// }
