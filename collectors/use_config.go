package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
)

func UseConfigCreateSubscription(rconf config.RuntimeConfig) ([]newcore.Subscription, error) {
	var subs []newcore.Subscription

	kafka_sub_names := make(map[string]bool)
	for _, conf := range rconf.Client.Subscribe_kafka {
		if conf != nil {
			// fmt.Printf("kafka_sub_names: %v\n", kafka_sub_names)
			_, ok := kafka_sub_names[conf.Name]
			if ok == true {
				logging.Errorf("duplicated kafka subscribe name are not allowed: %s", conf.Name)
				return nil, fmt.Errorf("duplicated kafka subscribe name are not allowed: %s", conf.Name)
			}
			kafka_sub_names[conf.Name] = true

			sub, err := NewKafkaSubscription(*conf)
			if err != nil {
				logging.Errorf("failed to create kafka subscription: %v", err)
				return nil, err
			}
			subs = append(subs, sub)
		}
	}
	return subs, nil
}
