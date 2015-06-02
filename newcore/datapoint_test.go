package newcore

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func Test_DataPoint_Clean(t *testing.T) {
	d := DataPoint{
		Metric:    "m1-:$^&%$*^&%2",
		Timestamp: time.Now(),
		Value:     1,
	}

	d.Clean()
	t.Log("d: ", d)
	if d.Metric != "m12" {
		t.Error("clean failed")
	}

	// t.Error("--")
}

func Test_DataPoint_MarshalJSON(t *testing.T) {
	d := DataPoint{
		Metric:    "m1-:$^&%$*^&%2",
		Timestamp: time.Now(),
		Value:     1,
	}

	v, err := d.MarshalJSON2String()
	t.Logf("%s, %v\n", v, err)
	if err != nil {
		t.Error("MarshalJSON failed")
	}
	// t.Error("--")

	var tmp_d DataPoint
	json.Unmarshal([]byte(v), &tmp_d)
	t.Logf("%v", tmp_d)
	t.Logf("%v", tmp_d.Timestamp.UnixNano())
	// t.Error("--")
}

func Test_MultiDataPoint_MarshalJSON(t *testing.T) {
	d := DataPoint{
		Metric:    "m1-:$^&%$*^&%2",
		Timestamp: time.Now(),
		Value:     1,
	}
	md := MultiDataPoint{}
	md = append(md, d)

	v, err := json.Marshal(md)
	t.Logf("%s, %v\n", v, err)
	if err != nil {
		t.Error("MarshalJSON failed")
	}
	// t.Error("--")

	var tmp_d MultiDataPoint
	json.Unmarshal(v, &tmp_d)
	t.Logf("%v", tmp_d)
	// t.Error("--")
}

func Test_DataPoint_Json(t *testing.T) {
	var d *DataPoint
	d = &DataPoint{
		Metric:    "hahaha",
		Timestamp: time.Now(),
		Value:     1.1,
	}

	t.Log(d.Json())

	var d2 DataPoint
	err := json.Unmarshal(d.Json(), &d2)
	if err != nil {
		t.Error("failed to unmarshal d.Json()")
	}

	t.Log(d2)

	if d2.Value != d.Value {
		t.Error("...")
	}
}

func Test_DataPoint_value2string(t *testing.T) {
	var d *DataPoint
	d = &DataPoint{
		Metric:    "hahaha",
		Timestamp: time.Now(),
		Value:     1.1,
	}

	t.Log(d.Json())

	assert_equal := func(dp *DataPoint, eq string) {
		s := dp.value2string()
		t.Logf("assert_equal: expect: %s, res: %s, value: %v", eq, s, dp.Value)
		if s != eq {
			t.Error("...")
		}
	}

	d.Value = 1.1
	assert_equal(d, `1.1`)
	d.Value = "1.1"
	assert_equal(d, `"1.1"`)

	d.Value = true
	assert_equal(d, `true`)
	d.Value = false
	assert_equal(d, `false`)

}

func Benchmark_DataPoint_fmt(b *testing.B) {
	// run the Fib function b.N times
	var (
		d     *DataPoint
		value string
	)

	d = &DataPoint{
		Metric:    "hahaha",
		Timestamp: time.Now(),
		Value:     1.1,
	}

	for n := 0; n < b.N; n++ {
		value = fmt.Sprintf("%+v", d)
	}
	b.Log(value)
}

func Benchmark_DataPoint_MarshalJSON(b *testing.B) {
	// run the Fib function b.N times
	var (
		d     *DataPoint
		value string
	)

	d = &DataPoint{
		Metric:    "hahaha",
		Timestamp: time.Now(),
		Value:     1.1,
	}

	for n := 0; n < b.N; n++ {
		value, _ = d.MarshalJSON2String()
	}
	b.Log(value)

}

// Benchmark_DataPoint_fmt   300000              5172 ns/op
// Benchmark_DataPoint_MarshalJSON   100000             12817 ns/op
// Benchmark_DataPoint_Json          300000              4012 ns/op

func Benchmark_DataPoint_Json(b *testing.B) {
	// run the Fib function b.N times
	var (
		d     *DataPoint
		value string
	)

	d = &DataPoint{
		Metric:    "hahaha",
		Timestamp: time.Now(),
		Value:     1.1,
	}

	for n := 0; n < b.N; n++ {
		value = string(d.Json())
	}
	b.Log(value)
}
