package config

type Config_KafkaSubscription struct {
	Name           string
	Broker_list    []string
	Topic          string
	Flush_interval string
	Max_batch_size int
}
