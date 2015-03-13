package collectorlib

import (
	"testing"
	"time"
)

func Test_DataPoint_MarshalJSON(t *testing.T) {
	d := DataPoint{
		Metric:    "m1",
		Timestamp: time.Now(),
		Value:     1,
	}

	v, err := d.MarshalJSON()
	t.Logf("%s, %v\n", v, err)
	t.Error("--")
}
