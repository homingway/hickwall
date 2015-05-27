package backends

import (
	"bytes"
	"github.com/oliveagle/hickwall/config"
	"testing"
)

var configs = map[string][]byte{
	"file": []byte(`
client:
    transport_file:
        path: "/var/lib/hickwall/fileoutput.txt"`),
}

func Test_UseConfig_CreateBackends(t *testing.T) {
	for key, data := range configs {
		rconf, err := config.ReadRuntimeConfig(bytes.NewBuffer(data))
		if err != nil {
			t.Errorf("%s, err %v", key, err)
			return
		}

		bks, err := UseConfigCreateBackends(rconf)
		if err != nil {
			t.Errorf("%s, err %v", key, err)
		}
		if len(bks) <= 0 {
			t.Errorf("%s, nothing created", key)
		}
		t.Logf("%s - %v", key, bks)
	}
}

var fail_configs = map[string][]byte{
	"file": []byte(`
client:
    transport_file:
        path: ""`),
}

func Test_UseConfig_CreateBackends_Fails(t *testing.T) {
	for key, data := range fail_configs {
		rconf, err := config.ReadRuntimeConfig(bytes.NewBuffer(data))
		if err != nil {
			t.Errorf("%s, err %v", key, err)
			return
		}

		bks, err := UseConfigCreateBackends(rconf)
		if err == nil {
			t.Errorf("%s, should fail but not", key)
		}
		if len(bks) > 0 {
			t.Errorf("%s, should fail but not", key)
		}
		t.Logf("%s - %v", key, bks)
	}
}
