package hickwall

import (
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
	"github.com/oliveagle/hickwall/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
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
	resp := &RegistryResponse{
		RequestHash: "hadhfasdfadsfasdf",
		Request: &RegistryRequest{
			Timestamp:  time.Now(),
			SystemInfo: sysinfo,
		},
		Timestamp:   time.Now(),
		Etcd_url:    "hahah",
		Config_path: "path",
	}

	err := resp.Save()
	if err != nil {
		t.Error("...")
		return
	}
	new_resp, err := LoadRegistryResponse()
	if err != nil {
		t.Error("load failed")
		return
	}
	// pretty.Println(new_resp)
	// t.Log(new_resp)

	if new_resp == resp {
		t.Error("the should not be the same one")
	}

	if new_resp.Etcd_url != "hahah" {
		t.Error("data is not correct")
	}

}

func Test_LoadRegistryResponse_Failed_Open(t *testing.T) {
	os.Remove(config.REGISTRY_FILEPATH)
	_, err := LoadRegistryResponse()
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

		var hashed HashedRegistryRequest
		err = json.Unmarshal(body, &hashed)
		if err != nil {
			t.Error("err", err)
			http.Error(w, "failed to unmarshal", 500)
			return
		}

		resp := RegistryResponse{
			RequestHash: hashed.Hash,
			Timestamp:   time.Now(),
			Etcd_url:    "http://localhost",
			Config_path: "/config/xxxx",
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
