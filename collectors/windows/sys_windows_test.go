package windows

import (
	c_conf "github.com/oliveagle/hickwall/collectors/config"
	"testing"
	"time"
)

func Test_sys_windows_nil(t *testing.T) {
	cs := MustNewWinSysCollectors("sys", "sys", nil)
	if len(cs) > 0 {
		t.Errorf("nil config should not create collector")
	}
}

func Test_sys_windows_1(t *testing.T) {
	conf := c_conf.Config_win_sys{Pdh_Interval: "1s", Wmi_Interval: "2s"}
	cs := MustNewWinSysCollectors("sys", "sys", &conf)
	if len(cs) != 2 {
		t.Errorf("failed to create sys windows collectors")
	}
	//	t.Log(cs)
	for _, c := range cs {
		if c.ClassName() == "win_pdh_collector" && c.Interval() != time.Second {
			t.Errorf("pdh collector interval should be 1s", c.Interval())
		}
		if c.ClassName() == "win_wmi_collector" && c.Interval() != time.Second*2 {
			t.Errorf("wmi collector interval should be 2s: %v", c.Interval())
		}
	}
}
