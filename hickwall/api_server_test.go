package hickwall

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/franela/goreq"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"github.com/oliveagle/hickwall/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"
	"time"
)

func Test_api_info(t *testing.T) {
	logging.SetLevel("debug")
	go serve()

	resp, err := http.Get("http://localhost:3031/sys_info?detail=true")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Log(resp.StatusCode, string(body))
}

func Test_api_registry_revoke(t *testing.T) {
	logging.SetLevel("debug")
	go serve()
	config.CORE_CONF_FILEPATH, _ = filepath.Abs("./test/core_config.yml")
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config_wo_groups.yml")

	ioutil.WriteFile(config.REGISTRY_FILEPATH, []byte(`test`), 0644)

	uri, _ := url.Parse("http://localhost:3031/registry/revoke")

	values := url.Values{}
	values.Add("hostname", newcore.GetHostname())
	values.Encode()
	uri.RawQuery = values.Encode()

	resp, err := goreq.Request{
		Method:  "DELETE",
		Uri:     uri.String(),
		Accept:  "application/json",
		Timeout: 1 * time.Second,
	}.Do()

	if err != nil {
		t.Errorf("failed to do revoke: %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("statuscode != 200, %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("failed to read response")
	}
	t.Logf("data: %s", data)

	b, err := utils.PathExists(config.REGISTRY_FILEPATH)
	if b != false || err != nil {
		t.Errorf("revoke not working.")
	}
}

var private_key = bytes.NewReader([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDCFENGw33yGihy92pDjZQhl0C36rPJj+CvfSC8+q28hxA161QF
NUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6Z4UMR7EOcpfdUE9Hf3m/hs+F
UR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJwoYi+1hqp1fIekaxsyQIDAQAB
AoGBAJR8ZkCUvx5kzv+utdl7T5MnordT1TvoXXJGXK7ZZ+UuvMNUCdN2QPc4sBiA
QWvLw1cSKt5DsKZ8UETpYPy8pPYnnDEz2dDYiaew9+xEpubyeW2oH4Zx71wqBtOK
kqwrXa/pzdpiucRRjk6vE6YY7EBBs/g7uanVpGibOVAEsqH1AkEA7DkjVH28WDUg
f1nqvfn2Kj6CT7nIcE3jGJsZZ7zlZmBmHFDONMLUrXR/Zm3pR5m0tCmBqa5RK95u
412jt1dPIwJBANJT3v8pnkth48bQo/fKel6uEYyboRtA5/uHuHkZ6FQF7OUkGogc
mSJluOdc5t6hI1VsLn0QZEjQZMEOWr+wKSMCQQCC4kXJEsHAve77oP6HtG/IiEn7
kpyUXRNvFsDE0czpJJBvL/aRFUJxuRK91jhjC68sA7NsKMGg5OXb5I5Jj36xAkEA
gIT7aFOYBFwGgQAQkWNKLvySgKbAZRTeLBacpHMuQdl1DfdntvAyqpAZ0lY0RKmW
G6aFKaqQfOXKCyWoUiVknQJAXrlgySFci/2ueKlIE1QqIiLSZ8V8OlpFLRnb1pzI
7U1yQXnTAEFYM560yJlzUpOb1V4cScGd365tiSMvxLOvTA==
-----END RSA PRIVATE KEY-----`))

var public_key = bytes.NewReader([]byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`))

func Test_api_registry_revoke_secure(t *testing.T) {
	logging.SetLevel("debug")
	go serve()
	config.CORE_CONF_FILEPATH, _ = filepath.Abs("./test/core_config.yml")
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config_wo_groups.yml")
	config.CoreConf.SecureAPIWrite = true // need to add signature

	unsigner, _ = utils.LoadPublicKey(public_key) // override interval unsigner
	signer, _ := utils.LoadPrivateKey(private_key)

	hostname := newcore.GetHostname()
	now := fmt.Sprintf("%d", time.Now().Unix())

	uri, _ := url.Parse("http://localhost:3031/registry/revoke")
	values := url.Values{}
	values.Add("hostname", newcore.GetHostname())
	values.Add("time", now)
	values.Encode()
	uri.RawQuery = values.Encode()

	go_req := goreq.Request{
		Method:  "DELETE",
		Uri:     uri.String(),
		Accept:  "application/json",
		Timeout: 1 * time.Second,
	}

	// mock registry file before call api
	ioutil.WriteFile(config.REGISTRY_FILEPATH, []byte(`test`), 0644)
	resp, _ := go_req.Do()
	//	if err == nil {
	//		t.Errorf("should fail but not.")
	//		return
	//	}

	if resp.StatusCode == 200 {
		t.Errorf("should fail but not")
		return
	}

	if b, _ := utils.PathExists(config.REGISTRY_FILEPATH); b != true {
		t.Errorf("should fail but not")
		return
	}
	resp.Body.Close()

	toSign := fmt.Sprintf("%s%s", hostname, now)
	sign, _ := signer.Sign([]byte(toSign))
	sign_str := base64.StdEncoding.EncodeToString(sign)
	go_req.AddHeader("HICKWALL_ADMIN_SIGN", sign_str)

	// mock registry file before call api
	ioutil.WriteFile(config.REGISTRY_FILEPATH, []byte(`test`), 0644)
	resp, err := go_req.Do()
	if err != nil {
		t.Errorf("should work but not: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("statuscode != 200, %d", resp.StatusCode)
		return
	}

	b, err := utils.PathExists(config.REGISTRY_FILEPATH)
	if b != false || err != nil {
		t.Errorf("revoke not working.")
	}
}
