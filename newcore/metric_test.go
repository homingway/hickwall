package newcore

import (
	"bytes"
	"github.com/oliveagle/viper"
	"testing"
)

func TestMetric(t *testing.T) {
	var metric Metric
	// metric := Metric("win.wmi.fs.d.cdfs.free_space.bytes\r\n")

	metric = "win.wmi.fs.d.cdfs.free_space.bytes\r\n"
	exp := "win.wmi.fs.d.cdfs.free_space.bytes"
	res := metric.Clean()

	t.Logf("*%s*", res)
	if res != exp {
		t.Error("failed")
	}
}

func TestMetricCleanWithTags(t *testing.T) {
	var metric Metric

	metric = "win.wmi.fs.d.cdfs.free_space.bytes\r\n"
	tpl := "hotel.{{.Tags_bu}}.{{.Tags_host}}.{{.Key}}.{{.Tags}}"
	tags := TagSet{
		"bu":   "hotel",
		"ax":   "ax001",
		"host": "host1",
	}

	exp := "hotel.hotel.host1.win.wmi.fs.d.cdfs.free_space.bytes.ax_ax001"

	res, err := metric.CleanWithTags(tpl, &tags)
	t.Log(res, err)
	if res != exp {
		t.Error("--")
	}
}

func TestViperParseMetric(t *testing.T) {
	type conf struct {
		Metric Metric `json:"metric"`
	}

	var c conf

	c_str := []byte(`{"metric": "win.wmi.fs.d.cdfs.free_space.bytes ahahah"}`)

	vp := viper.New()
	vp.SetConfigType("json")
	vp.ReadConfig(bytes.NewBuffer(c_str))
	vp.Marshal(&c)
	res := c.Metric.Clean()
	t.Log(res)
	if res != "win.wmi.fs.d.cdfs.free_space.bytesahahah" {
		t.Error("failed")
	}
}
