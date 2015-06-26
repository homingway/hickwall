package hickwall

import (
	"fmt"
	"github.com/oliveagle/hickwall/lib/wmi"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"github.com/oliveagle/hickwall/utils"
	"strconv"
)

func GetSystemInfo() (SystemInfo, error) {
	var info = SystemInfo{}
	cs_info, err := wmi.QueryWmi("SELECT Name, Domain, NumberOfLogicalProcessors, NumberOfProcessors, TotalPhysicalMemory FROM Win32_ComputerSystem")
	logging.Tracef("err: %v, cs_info: %v", err, cs_info)
	if err != nil {
		return info, err
	}
	if len(cs_info) != 1 {
		return info, fmt.Errorf("invalid query result: %v", cs_info)
	}
	cs_info_m := cs_info[0]

	info.Name = newcore.GetHostname()

	if string_value, ok := cs_info_m["Domain"]; ok == true {
		info.Domain = string_value
	}

	if string_value, ok := cs_info_m["NumberOfLogicalProcessors"]; ok == true {
		int_value, err := strconv.Atoi(string_value)
		if err != nil {
			return info, err
		}
		info.NumberOfLogicalProcessors = int_value
	}

	if string_value, ok := cs_info_m["NumberOfProcessors"]; ok == true {
		int_value, err := strconv.Atoi(string_value)
		if err != nil {
			return info, err
		}
		info.NumberOfProcessors = int_value
	}

	if string_value, ok := cs_info_m["TotalPhysicalMemory"]; ok == true {
		int_value, err := strconv.Atoi(string_value)
		if err != nil {
			return info, err
		}
		info.TotalPhsycialMemoryKb = int_value / 1024
	}

	os_info, err := wmi.QueryWmi("Select Caption, CSDVersion, OSArchitecture, Version From Win32_OperatingSystem")
	logging.Tracef("err: %v, os_info: %v", err, os_info)
	if err != nil {
		return info, err
	}
	if len(os_info) != 1 {
		return info, fmt.Errorf("invalid query result: %v", os_info)
	}
	os_info_m := os_info[0]

	if string_value, ok := os_info_m["Caption"]; ok == true {
		info.OS = string_value
	}

	csdversion := ""
	if string_value, ok := os_info_m["CSDVersion"]; ok == true {
		csdversion = string_value
	}

	if string_value, ok := os_info_m["Version"]; ok == true {
		version := string_value
		info.OSVersion = fmt.Sprintf("%s - %s", csdversion, version)
	}

	if string_value, ok := os_info_m["OSArchitecture"]; ok == true {
		if string_value == "64-bit" {
			info.Architecture = 64
		} else {
			info.Architecture = 32
		}

	}

	ipv4list, err := utils.Ipv4List()
	if err != nil {
		return info, err
	}
	info.IPv4 = ipv4list

	return info, nil
}
