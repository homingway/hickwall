package utils

import (
	"testing"
)

func Test_Ipv4List(t *testing.T) {
	iplist, err := Ipv4List()
	if err != nil {
		t.Error("")
	}
	if len(iplist) <= 0 {
		t.Error("")
	}

	t.Log("iplist: ", iplist)
}

func Test_Ipv4Map(t *testing.T) {
	ipmap, err := IpV4Map()
	if err != nil {
		t.Error("")
	}
	if len(ipmap) <= 0 {
		t.Error("")
	}

	t.Log("ipmap: ", ipmap)
}
