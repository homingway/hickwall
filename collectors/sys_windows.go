package collectors

import (
	"bytes"
	// "fmt"
	// "github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/utils"
	log "github.com/oliveagle/seelog"
	"github.com/oliveagle/viper"
	// "runtime"
	// "time"
)

func init() {
	defer utils.Recover_and_log()

	collector_factories["win_sys"] = factory_win_sys
}

func factory_win_sys(name string, conf interface{}) <-chan Collector {
	defer utils.Recover_and_log()

	log.Debug("factory_win_sys")

	_win_sys_pdh("100ms")
	_win_sys_wmi("60s")

	var out = make(chan Collector)
	go func() {
		// if conf != nil {
		// 	config_list := conf.(config.Conf_ping)

		// }
		close(out)
	}()
	return out
}

func _win_sys_pdh(interval string) {
	pdh_config_str_tpl := `
collector_win_pdh:
    -
        interval: {{.Interval}}
        queries:
            -
                query: "\\System\\Processes"
                metric: "sys.processes.count"
            -
                query: "\\Memory\\Available Bytes"
                metric: "sys.memory.avaiable.bytes"
            -
                query: "\\Processor(_Total)\\% Processor Time"
                metric: "sys.cpu.processor_time_pct.total"
            -
                query: "\\Process(hickwall)\\Working Set - Private"
                metric: "hickwall.client.mem.private_working_set.bytes"

`
	pdh_config_str, _ := utils.ExecuteTemplate(pdh_config_str_tpl, map[string]string{"Interval": interval}, nil)
	pdh_viper := viper.New()
	pdh_viper.SetConfigType("yaml")
	pdh_viper.ReadBufConfig(bytes.NewBuffer([]byte(pdh_config_str)))

	var tmp_conf = config.RuntimeConfig{}
	pdh_viper.Marshal(&tmp_conf)
	AddCollector("win_pdh", "win_sys_pdh", tmp_conf.Collector_win_pdh)
}

func _win_sys_wmi(interval string) {
	wmi_config_str_tpl := `
collector_win_wmi:
    -
        interval: {{.Interval}}
        queries:
            -
                query: "select Name, NumberOfCores from Win32_Processor"
                metrics:
                    -
                        value_from: "Name"
                        metric: "sys.cpu.name"
                    -
                        value_from: "NumberOfCores"
                        metric: "sys.cpu.numberofcores"
`
	wmi_config_str, _ := utils.ExecuteTemplate(wmi_config_str_tpl, map[string]string{"Interval": interval}, nil)
	wmi_viper := viper.New()
	wmi_viper.SetConfigType("yaml")
	wmi_viper.ReadBufConfig(bytes.NewBuffer([]byte(wmi_config_str)))

	var tmp_conf = config.RuntimeConfig{}
	wmi_viper.Marshal(&tmp_conf)

	AddCollector("win_wmi", "win_sys_wmi", tmp_conf.Collector_win_wmi)
}
