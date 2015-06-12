package config

type KafkaSubscription struct {
	Broker_list    []string
	Topic_id       string
	Flush_interval string
	Max_batch_size int
}
