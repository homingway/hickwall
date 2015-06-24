package utils

import (
	"errors"
	"fmt"
	"net"
)

func IpV4Map() (map[string]string, error) {
	ipmap := map[string]string{}

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
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
	iplist := []string{}

	ifaces, err := net.Interfaces()
	if err != nil {
		return []string{}, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return []string{}, err
		}
		for _, addr := range addrs {
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
			// fmt.Println(ip.String())
			iplist = append(iplist, ip.String())
		}
	}
	if len(iplist) > 0 {
		return iplist, nil
	} else {
		return []string{}, errors.New("are you connected to the network?")
	}
}
