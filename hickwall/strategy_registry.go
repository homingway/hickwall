package hickwall

import (
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

type RegistryRequest struct {
	Timestamp  time.Time  `json:"timestamp"`
	SystemInfo SystemInfo `json:"systeminfo"`
}

type HashedRegistryRequest struct {
	Hash       string `json:"hash"`
	RequestStr string `json:"request_str"`
}

func new_reg_request() (*RegistryRequest, error) {
	sysinfo, err := GetSystemInfo()
	if err != nil {
		return nil, err
	}

	return &RegistryRequest{
		Timestamp:  time.Now(),
		SystemInfo: sysinfo,
	}, nil
}

func new_hashed_reg_request(r *RegistryRequest) (*HashedRegistryRequest, error) {
	r_str, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	h := md5.New()
	h.Write(r_str)
	hash := hex.EncodeToString(h.Sum(nil))
	return &HashedRegistryRequest{
		Hash:       hash,
		RequestStr: string(r_str),
	}, nil
}

func new_reg_request_from_hashed(hr *HashedRegistryRequest) (*RegistryRequest, error) {
	h := md5.New()
	h.Write([]byte(hr.RequestStr))
	hash_expect := hex.EncodeToString(h.Sum(nil))
	if hr.Hash != hash_expect {
		return nil, fmt.Errorf("hash doesn't match: %s != %s", hr.Hash, hash_expect)
	}
	var rr RegistryRequest
	err := json.Unmarshal([]byte(hr.RequestStr), &rr)
	if err != nil {
		return nil, err
	}

	return &rr, nil
}

// ----------------------------------- RegistryResponse ----------------------------------

type RegistryResponse struct {
	Request *RegistryRequest `json:"request",omitempty`

	RequestHash string    `json:"request_hash"`
	Timestamp   time.Time `json:"timestamp"`
	Etcd_url    string    `json:"etcd_url"`
	Config_path string    `json:"config_path"`
}

type HashedRegistryResponse struct {
	Hash        string `json:"hash"`
	ResponseStr string `json:"response_str"`
}

func new_hashed_reg_resp(r *RegistryResponse) (*HashedRegistryResponse, error) {
	dump, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	h := md5.New()
	h.Write(dump)
	hash := hex.EncodeToString(h.Sum(nil))
	return &HashedRegistryResponse{
		Hash:        hash,
		ResponseStr: string(dump),
	}, nil
}

func new_hashed_reg_response_from_json(dump []byte) (*HashedRegistryResponse, error) {
	var hr HashedRegistryResponse
	// fmt.Println("dump ---- : ", string(dump))
	err := json.Unmarshal(dump, &hr)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal HashedRegistryResponse: %v", err)
	}
	return &hr, nil
}

func new_reg_resp_from_hashed(hr HashedRegistryResponse) (*RegistryResponse, error) {
	h := md5.New()
	h.Write([]byte(hr.ResponseStr))
	hash_expect := hex.EncodeToString(h.Sum(nil))
	if hr.Hash != hash_expect {
		return nil, fmt.Errorf("hash doesn't match: %s != %s", hr.Hash, hash_expect)
	}
	var resp RegistryResponse
	err := json.Unmarshal([]byte(hr.ResponseStr), &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RegistryResponse: %v", err)
	}
	return &resp, nil
}

func NewRegistryResponseFromJson(dump []byte) (*RegistryResponse, error) {
	hr, err := new_hashed_reg_response_from_json(dump)
	if err != nil {
		return nil, err
	}
	return new_reg_resp_from_hashed(*hr)
}

func (r *RegistryResponse) Save() error {
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

func do_registry(reg_url string) (*RegistryResponse, error) {
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

	resp, err := NewRegistryResponseFromJson(body)
	if err != nil {
		return nil, err
	}

	_, err = url.Parse(resp.Etcd_url)
	if err != nil {
		return nil, fmt.Errorf("invalid etcd_url: %s", err)
	}
	if resp.Config_path == "" {
		return nil, fmt.Errorf("config path is empty")
	}
	if resp.RequestHash != hReq.Hash {
		return nil, fmt.Errorf("request hash and response hash mismatch: %s != %s", hReq.Hash, resp.RequestHash)
	}
	resp.Request = req
	return resp, nil
}

func LoadRegistryResponse() (*RegistryResponse, error) {
	data, err := ioutil.ReadFile(config.REGISTRY_FILEPATH)
	if err != nil {
		return nil, err
	}
	var resp RegistryResponse

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func RegistryAndRun(stop chan error) {
	if stop == nil {
		panic("stop chan is nil")
	}

	resp, err := LoadRegistryResponse()
	if err != nil {
		// we don't have a valid registry info.
		tick := time.Tick(time.Minute * 5)

	registry_loop:
		for {
			select {
			case <-tick:
				resp, err = do_registry(config.CoreConf.RegistryURL)
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
	LoadConfigStrategyEtcd(resp.Etcd_url, resp.Config_path, stop)
}

//TODO: retrive registry server public key
