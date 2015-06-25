package config

import (
	"strings"
)

type Strategy string

const (
	FILE     Strategy = "file"
	ETCD              = "etcd"
	REGISTRY          = "registry"
)

func ValidStrategy(s string) string {
	k := strings.ToLower(s)
	switch k {
	case "file", "etcd", "registry":
		return k
	default:
		return ""
	}
}
