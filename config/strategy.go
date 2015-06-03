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

func (s *Strategy) IsValid() bool {
	k := strings.ToLower(string(*s))
	switch k {
	case "file", "etcd", "registry":
		return true
	default:
		return false
	}
}

func (s *Strategy) GetString() string {
	return strings.ToLower(string(*s))
}
