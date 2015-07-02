package windows

import (
	"bytes"
	c_conf "github.com/oliveagle/hickwall/collectors/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"github.com/oliveagle/viper"
	"strings"
)

func MustNewWinSysCollectors(name, prefix string, conf *c_conf.Config_win_sys) []newcore.Collector {
	if conf == nil {
		return nil
	}
	var collectors []newcore.Collector

	collectors = append(collectors, _must_win_sys_pdh(name, prefix, conf.Pdh_Interval))
	collectors = append(collectors, _must_win_sys_wmi(name, prefix, conf.Wmi_Interval))
	return collectors
}

func _must_win_sys_pdh(name, prefix, interval string) newcore.Collector {
	pdh_config_str_tpl := `
interval: {-----}
queries:
    -
        # CPU load
        query: "\\System\\Processor Queue Length"
        metric: "win.pdh.processor_queue_length"
    -
        query: "\\Processor(_Total)\\% Processor Time"
        metric: "win.pdh.pct_processor_time"
    -
        query: "\\Memory\\Available KBytes"
        metric: "win.pdh.memory.available_kbytes"
    -
        query: "\\Memory\\% Committed Bytes In Use"
        metric: "win.pdh.memory.pct_committed_bytes_in_use"
    -
        query: "\\PhysicalDisk(_Total)\\Avg. Disk sec/Read"
        metric: "win.pdh.physical_disk.avg_disk_sec_read"
    -
        query: "\\PhysicalDisk(_Total)\\Avg. Disk sec/Write"
        metric: "win.pdh.physical_disk.avg_disk_sec_write"
    -
        query: "\\TCPv4\\Connections Established"
        metric: "win.pdh.tcpv4.connections_established"
    -
        query: "\\System\\System Calls/sec"
        metric: "win.pdh.system.system_calls_sec"
    -
        query: "\\PhysicalDisk(_Total)\\Avg. Disk Bytes/Transfer"
        metric: "win.pdh.physical_disk.avg_disk_bytes_transfer"
    -
        query: "\\PhysicalDisk(_Total)\\Avg. Disk Queue Length"
        metric: "win.pdh.physical_disk.avg_disk_queue_length"
    -
        query: "\\PhysicalDisk(_Total)\\% Disk Time"
        metric: "win.pdh.physical_disk.pct_disk_time"
    -
        query: "\\LogicalDisk(C:)\\% Free Space"
        metric: "win.pdh.logicaldisk.pct_free_space_c"
    -
        query: "\\LogicalDisk(C:)\\Free Megabytes"
        metric: "win.pdh.logicaldisk.free_mbytes_c"
    -
        query: "\\LogicalDisk(D:)\\% Free Space"
        metric: "win.pdh.logicaldisk.pct_free_space_d"
    -
        query: "\\LogicalDisk(D:)\\Free Megabytes"
        metric: "win.pdh.logicaldisk.free_mbytes_d"
    -
        query: "\\TCPv4\\Connections Reset"
        metric: "win.pdh.tcpv4.connections_reset"
    -
        query: "\\TCPv4\\Connection Failures"
        metric: "win.pdh.tcpv4.connections_failures"

`
	pdh_config_str := strings.Replace(pdh_config_str_tpl, "{-----}", interval, 1)

	logging.Tracef("pdh_config_str: %s", pdh_config_str)

	pdh_viper := viper.New()
	pdh_viper.SetConfigType("yaml")
	pdh_viper.ReadConfig(bytes.NewBuffer([]byte(pdh_config_str)))
	var tmp_pdh_conf c_conf.Config_win_pdh_collector
	pdh_viper.Marshal(&tmp_pdh_conf)

	return MustNewWinPdhCollector(name, prefix, tmp_pdh_conf)
}

func _must_win_sys_wmi(name, prefix, interval string) newcore.Collector {
	wmi_config_str_tpl := `
interval: {-----}
queries:
    -
        query: "select Name, NumberOfCores from Win32_Processor"
        metrics:
            -
                value_from: "Name"
                metric: "win.wmi.cpu.name"
            -
                value_from: "NumberOfCores"
                metric: "win.wmi.cpu.numberofcores"
    -
        query: "select Name, Size from Win32_LogicalDisk where MediaType=11 or mediatype=12"
        metrics:
            -
                value_from: "Size"
                metric: "win.wmi.logicaldisk.total_size_bytes"
                tags: {
                    "mount": "{{.Name}}",
                }
`
	wmi_config_str := strings.Replace(wmi_config_str_tpl, "{-----}", interval, 1)

	logging.Tracef("pdh_config_str: %s", wmi_config_str)

	wmi_viper := viper.New()
	wmi_viper.SetConfigType("yaml")
	wmi_viper.ReadConfig(bytes.NewBuffer([]byte(wmi_config_str)))

	var tmp_wmi_conf c_conf.Config_win_wmi
	wmi_viper.Marshal(&tmp_wmi_conf)

	return MustNewWinWmiCollector(name, prefix, tmp_wmi_conf)
}
