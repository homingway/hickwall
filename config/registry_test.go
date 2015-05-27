package config

import (
	"encoding/json"
	"fmt"
	"github.com/oliveagle/hickwall/newcore"
	//	"github.com/oliveagle/hickwall/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	_ = fmt.Sprint("")
)

func init() {
	REGISTRY_FILEPATH = "./test_registry_file"
}

func TestRegistryCalledWithHostname(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.FormValue("hostname"))
		hostname := r.FormValue("hostname")
		if hostname != newcore.GetHostname() {
			t.Error("hostname mismatch")
		}

	}))
	defer ts.Close()

	Registry(ts.URL)
}

func TestRegistryReturnError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer ts.Close()

	reg, err := Registry(ts.URL)
	if err == nil {
		t.Error("should return error if no response body")
	}
	if reg != nil {
		t.Error("registry should not return registryinfo when something goes wrong.")
	}
}

func TestRegistrySuccess(t *testing.T) {
	os.Remove(REGISTRY_FILEPATH)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.FormValue("hostname"))
		hostname := r.FormValue("hostname")
		t.Log(hostname)

		reg_info := map[string]string{
			"etcd_registry_endpoint": "registry",
			"etcd_config_endpoint":   "config",
			"etcd_url":               "http://url",
			"cookie":                 "cookie",
		}
		d, _ := json.Marshal(reg_info)
		fmt.Println(string(d))
		fmt.Fprint(w, string(d))
	}))
	defer ts.Close()

	reg, err := Registry(ts.URL)
	if err != nil {
		t.Error(err)
	}

	if reg.Etcd_config_endpoing != "config" {
		t.Error("")
	}

	if IsRegistryFileExists() == false {
		t.Error("registry file dosen't exists")
	}
}

func TestReadRegistryFile(t *testing.T) {
	os.Remove(REGISTRY_FILEPATH)
	reg_info := map[string]string{
		"etcd_registry_endpoint": "registry",
		"etcd_config_endpoint":   "config",
		"etcd_url":               "http://url",
		"cookie":                 "cookie",
	}
	d, _ := json.Marshal(reg_info)
	ioutil.WriteFile(REGISTRY_FILEPATH, d, 0644)

	regInfo, err := ReadRegistryFile()
	if err != nil {
		t.Error("failed to read registry file")
	}

	if regInfo.Etcd_url != "http://url" {
		t.Error("reg info is not the one we saved!")
	}
}

func Test_ETCD_URL_PAT(t *testing.T) {
	failling_urls := []string{
		"some",
	}

	for _, u := range failling_urls {
		if ETCD_URL_PAT.MatchString(u) {
			t.Error("failling_urls should not be matched!", u)
		}
	}

	passing_urls := []string{
		"http://some",
		"http://some:1234",
		"https://some:1234/",
	}
	for _, u := range passing_urls {
		if !ETCD_URL_PAT.MatchString(u) {
			t.Error("passing_urls didn't matche!", u)
		}
	}
}

func TestRegistryCheckEtcdURLFailed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.FormValue("hostname"))
		hostname := r.FormValue("hostname")
		t.Log(hostname)

		reg_info := map[string]string{
			"etcd_registry_endpoint": "registry",
			"etcd_config_endpoint":   "config",
			"etcd_url":               "url",
			"cookie":                 "cookie",
		}
		d, _ := json.Marshal(reg_info)
		fmt.Println(string(d))
		fmt.Fprint(w, string(d))
	}))
	defer ts.Close()

	_, err := Registry(ts.URL)
	if err == nil {
		t.Error("etcd_url should be checked, and return error")
	}
	t.Log("err: ", err)
}
