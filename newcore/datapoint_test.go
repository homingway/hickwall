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
	// if d.Metric != "m12" { // ReplaceVersion
	if d.Metric != "m1_2" { // Regex Version
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

	v, err := d.MarshalJSON()
	t.Logf("%s, %v\n", v, err)
	if err != nil {
		t.Error("MarshalJSON failed")
	}
	// t.Error("--")

	var tmp_d DataPoint
	json.Unmarshal(v, &tmp_d)
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
	md = append(md, &d)

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

func Test_NewDP_set_hostname(t *testing.T) {
	SetHostname("test")
	d := NewDP("prefix", "metric", 1, nil, "", "", "")
	if d.Tags["host"] != "test" {
		t.Log(d.Tags)
		t.Error("SetHostname Doesn't work!")
	}
}

func Test_NewDPFromJson(t *testing.T) {
	content := []byte(`{"metric":"metric.c2","timestamp":"2015-06-09T18:44:49+08:00","value":1}`)
	d, err := NewDPFromJson(content)
	if err != nil {
		t.Log("err: %v", err)
		t.Error("should not raise error")
	}
	if d.Metric != "metric.c2" {
		t.Logf("%+v", d)
		t.Error("metric name is differnt.")
	}
	if d.Value.(float64) != 1.0 {
		t.Logf("%+v, value: %v", d, d.Value)
		t.Error("")
	}
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
		value []byte
	)

	d = &DataPoint{
		Metric:    "hahaha",
		Timestamp: time.Now(),
		Value:     1.1,
	}

	for n := 0; n < b.N; n++ {
		value, _ = d.MarshalJSON()
	}
	b.Log(value)

}
