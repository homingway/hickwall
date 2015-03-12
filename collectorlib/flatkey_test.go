package collectorlib

import (
	"testing"
)

func Test_FlatMetricKeyAndTags_NotAllowMultiLevelTpl(t *testing.T) {
	tpl := "hotel.{{.Tags.bu}}.{{.Tags_host}}.{{.Key}}.{{.Tags}}"
	key := "win.wmi.fs.d.cdfs.free_space.bytes\r\n"
	tags := map[string]string{
		"bu":   "hotel",
		"ax":   "ax001",
		"host": "host1",
	}

	res, err := FlatMetricKeyAndTags(tpl, key, tags)
	t.Log(res, err)
	if err == nil {
		t.Error("Multple Leveled Template should raise an error")
	}

	// t.Error("--")
}

func Test_FlatMetricKeyAndTags_TagNotFoundError(t *testing.T) {
	tpl := "hotel.{{.Tags_AAA}}.{{.Tags_host}}.{{.Key}}.{{.Tags}}"
	key := "win.wmi.fs.d.cdfs.free_space.bytes\r\n"
	tags := map[string]string{
		"bu":   "hotel",
		"ax":   "ax001",
		"host": "host1",
	}

	res, err := FlatMetricKeyAndTags(tpl, key, tags)
	t.Log(res, err)
	if err == nil {
		t.Error("tag AAA is not in tags")
	}
}

func Test_FlatMetricKeyAndTags_1(t *testing.T) {
	tpl := "hotel.{{.Tags_bu}}.{{.Tags_host}}.{{.Key}}.{{.Tags}}"
	key := "win.wmi.fs.d.cdfs.free_space.bytes\r\n"
	tags := map[string]string{
		"bu":   "hotel",
		"ax":   "ax001",
		"host": "host1",
	}

	exp := "hotel.hotel.host1.win.wmi.fs.d.cdfs.free_space.bytes.ax_ax001"

	res, err := FlatMetricKeyAndTags(tpl, key, tags)
	t.Log(res, err)
	if res != exp {
		t.Error("--")
	}
	// if err == nil {
	// 	t.Error("tag AAA is not in tags")
	// }
	// t.Error("---")
}
