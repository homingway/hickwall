// +build windows

package collectors

import (
	"github.com/oliveagle/hickwall/config"
	log "github.com/oliveagle/seelog"
	"os"
	"path/filepath"
	// "runtime"
	"testing"
	"time"
)

func TestCollector_cmd_1(t *testing.T) {
	conf := config.Conf_cmd{
		Cmd:      []string{"ahaha"},
		Interval: "1s",
		Tags: map[string]string{
			"bu": "hotel",
		},
	}

	collector := factory_cmd("test", conf)

	ic := collector.(*IntervalCollector)
	t.Log(ic)
	if ic.Interval != time.Duration(1)*time.Second {
		t.Error("failed to parse interval")
	}
	t.Log(ic.Interval)
}

func TestCollector_cmd_win_1(t *testing.T) {
	runtime_conf = config.GetRuntimeConf()
	runtime_conf.Log_console_level = "debug"
	config.ConfigLogger()
	defer log.Flush()

	dir, _ := os.Getwd()
	path := filepath.Join(dir, `tests\test_cmd.bat`)

	conf := config.Conf_cmd{
		Cmd: []string{
			`c:\windows\system32\cmd.exe`,
			`/c`,
			path,
		},
		Interval: "1s",
		Tags: map[string]string{
			"bu": "hotel",
		},
	}

	collector := factory_cmd("test", conf)

	ic := collector.(*IntervalCollector)
	md, err := ic.F(ic.states)

	if err != nil || len(md) <= 0 {
		t.Error("md len == 0")
	}

	t.Log(dir)

	t.Log(md)
	// t.Error("--")
}
