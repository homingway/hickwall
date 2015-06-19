package config

import (
	"github.com/oliveagle/viper"
	"io"
)

func ReadRuntimeConfig(r io.Reader) (*RuntimeConfig, error) {
	var rconf RuntimeConfig
	var err error

	vp := viper.New()
	vp.SetConfigType("yaml")
	err = vp.ReadConfig(r)
	if err != nil {
		return nil, err
	}
	err = vp.Marshal(&rconf)
	if err != nil {
		return nil, err
	}

	return &rconf, nil
}

type RuntimeConfig struct {
	Client ClientConfig           `json:"client"`
	Groups []CollectorConfigGroup `json:"groups"`
}
