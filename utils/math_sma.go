// simplest mean for quick reference only.

package utils

import (
	"fmt"
	"math"
)

type SMA struct {
	array []float64
	N     int
}

//FIXME: .\math_sma.go:15: leaking param: m
func (m *SMA) Calc(x float64) (float64, error) {
	if m.N <= 0 {
		return math.NaN(), fmt.Errorf("cannot calculate mean of zero")
	}

	if m.array == nil {
		m.array = []float64{}
	}

	m.array = append(m.array, x)

	if len(m.array) < m.N {
		return math.NaN(), fmt.Errorf("not enough data")
	}

	if len(m.array) > m.N {
		m.array = m.array[len(m.array)-m.N : len(m.array)]
	}

	sum := float64(0.0)
	for _, v := range m.array {
		sum += v
	}

	return sum / float64(m.N), nil
}
