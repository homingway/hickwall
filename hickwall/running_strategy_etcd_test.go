package hickwall

import (
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
	"github.com/oliveagle/hickwall/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	_ = pretty.Sprintf("")
)

type _etcd_resp_node struct {
	CreatedIndex  int    `json:"createdIndex"`
	Key           string `json:"key"`
	ModifiedIndex int    `json:"modifiedIndex"`
	Value         string `json:"value"`
}

type _etcd_resp struct {
	Action string          `json:"action"`
	Node   _etcd_resp_node `json:"node"`
}

func _fack_etcd_respose(idx int, value string) (string, error) {
	x := _etcd_resp{"get", _etcd_resp_node{idx, "/message", 1, value}}
	content, err := json.Marshal(x)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// config content, expecting?, error message?

var nil_core_tests = [][]string{
	{
		"",
		"no message",
	},

	{
		"}", // cannot parse this yaml content. will panic. we should handle it.
		"no message",
	},
	{
		`
# ---------- client configurations -------------------------------------------
client:
    # customize hostname, if omit this option, the client will use $(hostname) automatically
    # hostname: "test_server1"

    # minimal interval is 1s
    heartbeat_interval: 1s

    metric_enabled: true
    metric_interval: 2s
`,
		"no message",
	},
}

func Test_LoadConfigStrategyEtcd(t *testing.T) {
	stopCh := make(chan error)

	for idx, tcase := range nil_core_tests {
		request_cnt := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			request_cnt += 1
			fmt.Printf(" case: %d, -- we got request: cnt:%d ------------\n", idx, request_cnt)
			// pretty.Println(r)
			iswait := r.URL.Query().Get("wait")
			if strings.ToLower(iswait) == "true" {
				// watch, long polling
				time.Sleep(time.Second * 5)
			}

			v, _ := _fack_etcd_respose(1, tcase[0])
			fmt.Printf("case: %d, content: %s\n", idx, v)

			fmt.Fprintln(w, v)

		}))

		config.CoreConf.Etcd_url = ts.URL
		config.CoreConf.Etcd_path = "/config/host/DST54869.yml"
		go LoadConfigStrategyEtcd(stopCh)
		tick := time.After(time.Second * 1)
		timeout := time.After(time.Second * 2)

	main_loop:
		for {
			select {
			case <-tick:
				stopCh <- nil
				if the_core != nil {
					t.Error("the_core has been created with invalid configuration!")
					return
				} else {
					t.Log("test case successed")
					ts.Close()
					break main_loop
				}
			case <-timeout:
				t.Errorf("timed out. somethings is blocking")
				return
			}
		}
	}

}
