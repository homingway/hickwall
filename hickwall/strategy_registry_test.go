package hickwall

import (
	"encoding/json"
	"fmt"
	"github.com/franela/goreq"
	"github.com/kr/pretty"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	_ = pretty.Sprintf("")
	_ = fmt.Sprint("")
)

func Test_RegistryRequest_Hashed(t *testing.T) {
	req, _ := new_reg_request()
	hash, err := new_hashed_reg_request(req)
	if err != nil {
		t.Error("failed")
	}
	t.Log(err, hash)

	new_req, err := new_reg_request_from_hashed(hash)
	if err != nil {
		t.Error("failed")
	}
	if new_req.Timestamp != req.Timestamp {
		t.Error("unmarshal failed.")
	}
}

func Test_RegistryResponse_Save_and_Load(t *testing.T) {
	sysinfo, _ := GetSystemInfo()
	resp := &registry_response{
		RequestHash: "hadhfasdfadsfasdf",
		Request: &registry_request{
			Timestamp:  time.Now(),
			SystemInfo: sysinfo,
		},
		Timestamp:      time.Now(),
		EtcdMachines:   []string{"hahah"},
		EtcdConfigPath: "path",
	}

	err := resp.Save()
	if err != nil {
		t.Error("...")
		return
	}
	new_resp, err := load_reg_response()
	if err != nil {
		t.Error("load failed")
		return
	}
	// pretty.Println(new_resp)
	// t.Log(new_resp)

	if new_resp == resp {
		t.Error("the should not be the same one")
	}

	if new_resp.EtcdMachines[0] != "hahah" {
		t.Error("data is not correct")
	}

}

func Test_LoadRegistryResponse_Failed_Open(t *testing.T) {
	os.Remove(config.REGISTRY_FILEPATH)
	_, err := load_reg_response()
	if err == nil {
		t.Error("should raise error if nothing to load")
	}

}

func Test_Do_Registry(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// pretty.Println(r)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("err", err)
			http.Error(w, "failed to read body", 500)
			return
		}
		defer r.Body.Close()

		var hashed hashed_registry_request
		err = json.Unmarshal(body, &hashed)
		if err != nil {
			t.Error("err", err)
			http.Error(w, "failed to unmarshal", 500)
			return
		}

		resp := registry_response{
			RequestHash:    hashed.Hash,
			Timestamp:      time.Now(),
			EtcdMachines:   []string{"http://localhost"},
			EtcdConfigPath: "/config/xxxx",
		}
		hResp, err := new_hashed_reg_resp(&resp)
		if err != nil {
			t.Error("err", err)
			http.Error(w, "new_hashed_reg_resp failed", 500)
			return
		}

		dump, err := json.Marshal(hResp)
		if err != nil {
			t.Error("err", err)
			http.Error(w, "Marshal failed", 500)
			return
		}
		fmt.Println("dump: ", string(dump))

		fmt.Fprintln(w, string(dump))
	}))

	resp, err := do_registry(ts.URL)
	if err != nil {
		t.Error("do_registry failed")
	}
	t.Log(err, resp)
}

func Test_strategy_registry_enable_api_server(t *testing.T) {
	logging.SetLevel("info")
	config.CORE_CONF_FILEPATH, _ = filepath.Abs("./test/core_config_registry_enable_api.yml")
	//	config.CONF_FILEPATH, _ = filepath.Abs("./test/config.yml")
	err := config.LoadCoreConfig()
	if err != nil || config.CoreConf.Enable_http_api != true {
		t.Error("CoreConf.EnableHTTPAPI != true, err: %v", err)
		return
	}

	Stop() // stop if already exists while test all cases
	Start()

	//	resp, err := http.Get("http://localhost:3031/sys_info")
	hResp, err := goreq.Request{
		Method:      "Get",
		Uri:         "http://localhost:3031/sys_info",
		Accept:      "application/json",
		ContentType: "application/json",
		UserAgent:   "hickwall",
		Timeout:     100 * time.Millisecond,
	}.Do()
	if err != nil {
		t.Errorf("api server doesn't work. %v", err)
		return
	}
	defer hResp.Body.Close()
	t.Log(hResp)

}

func Test_strategy_registry_disable_api_server(t *testing.T) {
	config.CORE_CONF_FILEPATH, _ = filepath.Abs("./test/core_config_registry_disable_api.yml")
	err := config.LoadCoreConfig()
	if err != nil || config.CoreConf.Enable_http_api != false {
		t.Error("CoreConf.EnableHTTPAPI != false, err: %v, %v", err, config.CoreConf.Enable_http_api)
		return
	}

	Stop() // stop if already exists while test all cases
	Start()

	//	resp, err := http.Get("http://localhost:3031/sys_info")
	hResp, err := goreq.Request{
		Method:      "Get",
		Uri:         "http://localhost:3034/sys_info", // use disabled api port 3034
		Accept:      "application/json",
		ContentType: "application/json",
		UserAgent:   "hickwall",
		Timeout:     100 * time.Millisecond,
	}.Do()
	if err == nil {
		t.Error("api server still working.")
		defer hResp.Body.Close()
	}
}
