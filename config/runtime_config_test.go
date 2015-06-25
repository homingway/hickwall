package config

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func Test_ReadRuntimeConfig_hash(t *testing.T) {
	data := `name: test`
	rconf, err := ReadRuntimeConfig(bytes.NewReader([]byte(data)))
	if rconf == nil || err != nil {
		t.Errorf("cannot read runtime config: %v", err)
	}
	if len(rconf.hash) <= 0 {
		t.Error("hash doesn't exists")
	}
	if len(rconf.rawdata) <= 0 {
		t.Errorf("rawdata len <=0")
	}
}

func Test_DumpRuntimeConfig(t *testing.T) {
	data := `name: test`
	rconf, _ := ReadRuntimeConfig(bytes.NewReader([]byte(data)))
	err := DumpRuntimeConfig(rconf)
	if err != nil {
		t.Errorf("failed dump")
		return
	}
	ndata, err := ioutil.ReadFile(CONF_CACHE_PATH)
	if err != nil {
		t.Errorf("failed to read cache: %v", err)
	}
	if data != string(ndata) {
		t.Errorf("dumpped data is not the same we gave.")
	}
}
