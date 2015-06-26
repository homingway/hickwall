package config

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/viper"
	"io"
	"io/ioutil"
)

// return md5 hash of cached runtime config.
func GetCachedRuntimeConfigHash() (string, error) {
	data, err := ioutil.ReadFile(CONF_CACHE_PATH)
	if err != nil {
		return "", err
	}
	h := md5.New()
	h.Write(data)
	hash := hex.EncodeToString(h.Sum(nil))
	return hash, nil
}

// Dump RuntimeConfig rawdata content into a file. override if file already exists
func DumpRuntimeConfig(rconf *RuntimeConfig) error {
	logging.Debug("DumpRuntimeConfig")
	if len(rconf.rawdata) <= 0 {
		return fmt.Errorf("runtime config rawdata length <= 0")
	}
	h := md5.New()
	h.Write(rconf.rawdata)
	hash := hex.EncodeToString(h.Sum(nil))
	if hash != rconf.hash {
		return fmt.Errorf("rawdata has been modified!")
	}
	err := ioutil.WriteFile(CONF_CACHE_PATH, rconf.rawdata, 0644)
	if err != nil {
		return fmt.Errorf("failed to dump RuntimeConfig: %v", err)
	}
	logging.Debug("DumpRuntimeConfig Finished")
	return nil
}

// Read RuntimeConfig from io.Reader. will append hash of read data.
func ReadRuntimeConfig(r io.Reader) (*RuntimeConfig, error) {
	var rconf RuntimeConfig
	var err error

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// we have to use viper to load config. coz yaml cannot unmarshal
	// complicated struct. inside of viper. it use mapstructure to do it.
	nr := bytes.NewReader(data)
	vp := viper.New()
	vp.SetConfigType("yaml")
	err = vp.ReadConfig(nr)
	if err != nil {
		return nil, err
	}
	err = vp.Marshal(&rconf)
	if err != nil {
		return nil, err
	}

	h := md5.New()
	h.Write(data)
	rconf.hash = hex.EncodeToString(h.Sum(nil))
	rconf.rawdata = data

	return &rconf, nil
}

type RuntimeConfig struct {
	hash    string // hash of yaml content
	rawdata []byte // raw yaml content

	Client ClientConfig           `json:"client"`
	Groups []CollectorConfigGroup `json:"groups"`
}

func (r *RuntimeConfig) GetHash() string {
	return r.hash
}

func (r *RuntimeConfig) GetRawdata() []byte {
	return r.rawdata[:]
}
