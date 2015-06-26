package hickwall

import (
	"container/ring"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/franela/goreq"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
	"io/ioutil"
	"net/url"
	"time"
)

var (
	_ = fmt.Sprintf("")
)

// ----------------------------------- RegistryRequest ----------------------------------

type registry_request struct {
	Timestamp  time.Time  `json:"timestamp"`
	SystemInfo SystemInfo `json:"systeminfo"`
}

type hashed_registry_request struct {
	Hash       string `json:"hash"`
	RequestStr string `json:"request_str"`
}

func new_reg_request() (*registry_request, error) {
	sysinfo, err := GetSystemInfo()
	if err != nil {
		return nil, err
	}

	return &registry_request{
		Timestamp:  time.Now(),
		SystemInfo: sysinfo,
	}, nil
}

func new_hashed_reg_request(r *registry_request) (*hashed_registry_request, error) {
	r_str, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	h := md5.New()
	h.Write(r_str)
	hash := hex.EncodeToString(h.Sum(nil))
	return &hashed_registry_request{
		Hash:       hash,
		RequestStr: string(r_str),
	}, nil
}

func new_reg_request_from_hashed(hr *hashed_registry_request) (*registry_request, error) {
	h := md5.New()
	h.Write([]byte(hr.RequestStr))
	hash_expect := hex.EncodeToString(h.Sum(nil))
	if hr.Hash != hash_expect {
		return nil, fmt.Errorf("hash doesn't match: %s != %s", hr.Hash, hash_expect)
	}
	var rr registry_request
	err := json.Unmarshal([]byte(hr.RequestStr), &rr)
	if err != nil {
		return nil, err
	}

	return &rr, nil
}

// ----------------------------------- RegistryResponse ----------------------------------

type registry_response struct {
	Request *registry_request `json:"request",omitempty`

	RequestHash    string    `json:"request_hash"`
	Timestamp      time.Time `json:"timestamp"`
	EtcdMachines   []string  `json:"etcd_machines"`
	EtcdConfigPath string    `json:"etcd_config_path"`
}

type hashed_registry_response struct {
	Hash        string `json:"hash"`
	ResponseStr string `json:"response_str"`
}

func new_hashed_reg_resp(r *registry_response) (*hashed_registry_response, error) {
	dump, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	h := md5.New()
	h.Write(dump)
	hash := hex.EncodeToString(h.Sum(nil))
	return &hashed_registry_response{
		Hash:        hash,
		ResponseStr: string(dump),
	}, nil
}

func new_hashed_reg_response_from_json(dump []byte) (*hashed_registry_response, error) {
	var hr hashed_registry_response
	// fmt.Println("dump ---- : ", string(dump))
	err := json.Unmarshal(dump, &hr)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal HashedRegistryResponse: %v", err)
	}
	return &hr, nil
}

func new_reg_resp_from_hashed(hr hashed_registry_response) (*registry_response, error) {
	h := md5.New()
	h.Write([]byte(hr.ResponseStr))
	hash_expect := hex.EncodeToString(h.Sum(nil))
	if hr.Hash != hash_expect {
		return nil, fmt.Errorf("hash doesn't match: %s != %s", hr.Hash, hash_expect)
	}
	var resp registry_response
	err := json.Unmarshal([]byte(hr.ResponseStr), &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RegistryResponse: %v", err)
	}
	return &resp, nil
}

func new_reg_response_from_json(dump []byte) (*registry_response, error) {
	hr, err := new_hashed_reg_response_from_json(dump)
	if err != nil {
		return nil, err
	}
	return new_reg_resp_from_hashed(*hr)
}

func (r *registry_response) Save() error {
	dump, err := json.Marshal(r)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(config.REGISTRY_FILEPATH, dump, 0644)
	if err != nil {
		return err
	}

	return nil
}

func do_registry(reg_url string) (*registry_response, error) {
	req, err := new_reg_request()
	if err != nil {
		return nil, err
	}
	hReq, err := new_hashed_reg_request(req)
	if err != nil {
		return nil, err
	}

	hResp, err := goreq.Request{
		Method:      "POST",
		Uri:         reg_url,
		Body:        hReq,
		Accept:      "application/json",
		ContentType: "application/json",
		UserAgent:   "hickwall",
		Timeout:     10 * time.Second,
	}.Do()

	if serr, ok := err.(*goreq.Error); ok {
		if serr.Timeout() {
			return nil, fmt.Errorf("registry timed out.")
		}
		return nil, fmt.Errorf("registry failed: %d", hResp.StatusCode)
	}
	defer hResp.Body.Close()

	body, err := ioutil.ReadAll(hResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body %v", err)
	}

	resp, err := new_reg_response_from_json(body)
	if err != nil {
		return nil, err
	}

	if len(resp.EtcdMachines) <= 0 {
		return nil, fmt.Errorf("EtcdMachines is empty")
	}

	for _, m := range resp.EtcdMachines {
		_, err = url.Parse(m)
		if err != nil {
			return nil, fmt.Errorf("invalid etcd machine url: %s", err)
		}
	}

	if resp.EtcdConfigPath == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	if resp.RequestHash != hReq.Hash {
		return nil, fmt.Errorf("request hash and response hash mismatch: %s != %s", hReq.Hash, resp.RequestHash)
	}
	resp.Request = req
	return resp, nil
}

func load_reg_response() (*registry_response, error) {
	data, err := ioutil.ReadFile(config.REGISTRY_FILEPATH)
	if err != nil {
		return nil, err
	}
	var resp registry_response

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func new_core_from_registry(stop chan error) {
	if stop == nil {
		panic("stop chan is nil")
	}

	if len(config.CoreConf.Registry_urls) <= 0 {
		logging.Criticalf("RegistryURLs is empty!!")
		panic("RegistryURLS is empty!!")
		//		return fmt.Errorf("RegistryURLS is empty!!")
	}

	resp, err := load_reg_response()
	if err != nil {
		// we don't have a valid registry info.
		tick := time.Tick(time.Minute * 5)

		// round robin registry machines
		r := ring.New(len(config.CoreConf.Registry_urls))
		for i := 0; i < r.Len(); i++ {
			r.Value = config.CoreConf.Registry_urls[i]
			r = r.Next()
		}

	registry_loop:
		for {
			select {
			case <-tick:
				r = r.Next()
				resp, err = do_registry(r.Value.(string))
				if err == nil {
					// we are registried.
					break registry_loop
				} else {
					logging.Errorf("failed to registry: %v", err)
				}
			}
		}
	}

	// here we got a valid registry info. get config and start to run.
	new_core_from_etcd(resp.EtcdMachines, resp.EtcdConfigPath, stop)
}

//TODO: retrive registry server public key
//TODO: what if SystemInfo changed after registration?
