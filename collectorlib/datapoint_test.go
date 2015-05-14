package collectorlib

import (
	"encoding/json"
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
