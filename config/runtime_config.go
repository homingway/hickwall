package config

import (
	"github.com/oliveagle/viper"
	"io"
)

func ReadRuntimeConfig(r io.Reader) (rc RuntimeConfig, err error) {
	vp := viper.New()
	vp.SetConfigType("yaml")
	err = vp.ReadConfig(r)
	if err != nil {
		return rc, err
	}
	err = vp.Marshal(&rc)
	if err != nil {
		return rc, err
	}

	return rc, nil
}

type RuntimeConfig struct {
	Client ClientConfig           `json:"client"`
	Groups []CollectorConfigGroup `json:"groups"`
}
