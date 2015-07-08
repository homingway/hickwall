package utils

import (
	"errors"
	"fmt"
	//	"github.com/kr/pretty"
	"net"
)

func IpV4Map() (map[string]string, error) {
	//	fmt.Println("-----------------------------------------------------")
	ipmap := map[string]string{}

	ifaces, err := net.Interfaces()
	//	fmt.Printf("err: %v, ifaces %v\n", err, ifaces)
	//	pretty.Println(ifaces)
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		// go 1.4.2 on windows 2012, iface.Flags always return 0x0

		//		if iface.Flags&net.FlagUp == 0 {
		//			continue // interface down
		//		}

		//		if iface.Flags&net.FlagLoopback != 0 {
		//			continue // loopback interface
		//		}

		addrs, err := iface.Addrs()
		//		fmt.Printf("err: %v, address: %s\n", err, addrs)
		if err != nil {
			return nil, err
		}
		for idx, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			key := fmt.Sprintf("%s(%d)", iface.HardwareAddr.String(), idx)
			ipmap[key] = ip.String()
		}
	}
	if len(ipmap) > 0 {
		return ipmap, nil
	} else {
		return nil, errors.New("are you connected to the network?")
	}
}

func Ipv4List() ([]string, error) {
	ipmap, err := IpV4Map()
	if err != nil {
		return nil, err
	}

	iplist := []string{}
	for _, ip := range ipmap {
		iplist = append(iplist, ip)
	}
	return iplist, nil
}
