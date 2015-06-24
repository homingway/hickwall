package wmi

import (
	"testing"
)

func Test_QueryWmi(t *testing.T) {
	query := "SELECT TotalPhysicalMemory, Name, Domain, NumberOfLogicalProcessors, NumberOfProcessors, SystemType FROM Win32_ComputerSystem"
	res, err := QueryWmi(query)
	if err != nil {
		t.Error("...")
	}
	t.Log(res)
	if len(res) <= 0 {
		t.Error("...")
	}
	if name, ok := res[0]["Name"]; ok != true || name == "" {
		t.Error("...")
	}
}

func Test_parseFieldsFromQuery(t *testing.T) {
	query := "SELECT TotalPhysicalMemory, Name, Domain, NumberOfLogicalProcessors, NumberOfProcessors, SystemType FROM Win32_ComputerSystem"
	fields, err := parseFieldsFromQuery(query)
	if err != nil {
		t.Error("...")
	}
	if len(fields) != 6 {
		t.Error("...")
	}
	t.Log(fields)

	query = "SELECT * FROM Win32_ComputerSystem"
	fields, err = parseFieldsFromQuery(query)
	if err != nil {
		t.Error("...")
	}
	if len(fields) != 1 {
		t.Error("...")
	}
	if fields[0] != "*" {
		t.Error("...")
	}
}
