package utils

import (
	"reflect"
)

func Convert(x interface{}) interface{} {
	switch t := x.(type) {
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(t).Int()) // a has type int64
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(t).Uint()) // a has type uint64
	case float64:
		return float64(reflect.ValueOf(t).Float()) // a has type float64
	default:
		return reflect.ValueOf(t).String()
	}
}
