package config

import (
	"bytes"
	"github.com/oliveagle/viper"
	"testing"
)

func TestParseStrategy1(t *testing.T) {
	var conf_str = []byte(`
config_strategy: "file"
`)
	v := viper.New()
	v.SetConfigType("yaml")
	v.ReadConfig(bytes.NewReader(conf_str))

	var cc CoreConfig
	v.Marshal(&cc)
	t.Log(cc)
	if cc.Config_Strategy.GetString() != "file" {
		t.Error("")
	}

	if cc.Config_Strategy.IsValid() != true {
		t.Error("")
	}
}

func TestParseStrategy2(t *testing.T) {
	var conf_str = []byte(`
config_strategy: "xxxx"
`)
	v := viper.New()
	v.SetConfigType("yaml")
	v.ReadConfig(bytes.NewReader(conf_str))

	var cc CoreConfig
	v.Marshal(&cc)
	t.Log(cc)
	t.Log(cc.Config_Strategy)
	if cc.Config_Strategy.IsValid() != false {
		t.Error("")
	}
}
