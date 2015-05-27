package config

import (
	"encoding/json"
	"fmt"
	"github.com/oliveagle/hickwall/newcore"
	"io/ioutil"
	"net/http"
	"net/url"
	//	"os"
	//	valid "github.com/asaskevich/govalidator"
	"errors"
	//	"github.com/coreos/go-etcd/etcd"
	"github.com/oliveagle/hickwall/utils"
	"os"
	"regexp"
)

var (
	ETCD_URL_PAT = regexp.MustCompile(`^http[s]?://.*(:\d+)?[/]?`)
)

var (
	ERR_REGISTRY_REVOKED = errors.New("registry has been revoked!")
)

type RegistryInfo struct {
	Etcd_registry_endpoint string `json:"etcd_registry_endpoint"`
	Etcd_config_endpoing   string `json:"etcd_config_endpoint"`
	Etcd_url               string `json:"etcd_url"`
	Cookie                 string `json:"cookie"`
}

func Registry(reg_url string) (*RegistryInfo, error) {
	var (
		regInfo RegistryInfo
		err     error
	)
	v := url.Values{}
	v.Set("hostname", newcore.GetHostname())

	resp, err := http.PostForm(reg_url, v)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if resp.ContentLength > 1024 {
			return nil, fmt.Errorf("response is too big")
		}
		d, _ := ioutil.ReadAll(resp.Body)

		// get resp.Body and save it to registry file
		err = json.Unmarshal(d, &regInfo)
		if err != nil {
			return nil, err
		}

		if !ETCD_URL_PAT.MatchString(regInfo.Etcd_url) {
			return nil, fmt.Errorf("etcd url is not a valid url: %s", regInfo.Etcd_url)
		}

		// write registry info into file
		err = ioutil.WriteFile(REGISTRY_FILEPATH, d, 0644)
		if err != nil {
			return nil, fmt.Errorf("write registry file failed: %v", err)
		}
	} else {
		return nil, fmt.Errorf("error response,  code: %d, status: %v", resp.StatusCode, resp.Status)
	}
	return &regInfo, nil
}

//TODO: check registry file , read back registry info

func IsRegistryFileExists() bool {
	exists, err := utils.PathExists(REGISTRY_FILEPATH)
	if exists == true && err == nil {
		return true
	}
	return false
}

func ReadRegistryFile() (*RegistryInfo, error) {
	var (
		regInfo RegistryInfo
		err     error
	)

	stat, err := os.Stat(REGISTRY_FILEPATH)
	if err != nil {
		return nil, err
	}

	if stat.Size() > 1024 {
		return nil, fmt.Errorf("registry file is too big")
	}

	data, err := ioutil.ReadFile(REGISTRY_FILEPATH)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &regInfo)
	if err != nil {
		return nil, err
	}

	return &regInfo, nil
}

// 1. Check Do we have registry info stored on disk
// 2. if we don't, do registry process repeatedly until success. interval default is 5m
// 3. if we do have registry info.  use etcd to get runtime config.
// 4. if we found our registry into on etcd has been revoked, remove local registry info and
// 	  go to step-2, start registry process again
//func LoadConfigFromRegistry() error {
//
//}

//func WatchEtcdRegistryInfo(url, path string) error {
//	//TODO: remove viper and use etcd directly.
//	client := etcd.NewClient([]string{url})
//
//	return nil
//}
