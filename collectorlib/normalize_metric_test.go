package collectorlib

import (
	"testing"
)

func Test_normalize_metric_key_1(t *testing.T) {
	metric := "win.wmi.fs.d.cdfs.free_space.bytes\r\n"

	exp := "win.wmi.fs.d.cdfs.free_space.bytes"
	res := NormalizeMetricKey(metric)
	t.Logf("*%s*", res)
	if res != exp {
		t.Error("failed")
	}

}

func Test_normalize_metric_key_2(t *testing.T) {
	// ascii, and space
	metric := "  win.wmi.fs.D:.CDFS.free space.bytes 中文"

	exp := "win.wmi.fs.d.cdfs.free_space.bytes"
	res := NormalizeMetricKey(metric)
	t.Logf("*%s*", res)
	if res != exp {
		t.Error("failed")
	}
}

func Test_normalize_metric_key_3(t *testing.T) {
	// '-'' to '_'
	metric := "  win.wmi.fs.D:.CDFS.free-space.bytes"

	exp := "win.wmi.fs.d.cdfs.free_space.bytes"
	res := NormalizeMetricKey(metric)
	t.Logf("*%s*", res)
	if res != exp {
		t.Error("failed")
	}
}

func Test_normalize_metric_key_4(t *testing.T) {
	// numbers
	metric := " win.wmi.fs.D:.CDFS.free-space.bytes.112121"

	exp := "win.wmi.fs.d.cdfs.free_space.bytes.112121"
	res := NormalizeMetricKey(metric)
	t.Logf("*%s*", res)
	if res != exp {
		t.Error("failed")
	}
	// t.Error("---")
}

func Test_normalize_metric_key_5(t *testing.T) {
	// path
	metric := " win.wmi. /path/to/file"

	exp := "win.wmi.__path_to_file"
	res := NormalizeMetricKey(metric)
	t.Logf("*%s*", res)
	if res != exp {
		t.Error("failed")
	}
	// t.Error("---")
}

func Test_normalize_metric_key_6(t *testing.T) {
	// path
	metric := " win.wmi. c:\\path\\to\\file"

	exp := "win.wmi._c_path_to_file"
	res := NormalizeMetricKey(metric)
	t.Logf("*%s*", res)
	if res != exp {
		t.Error("failed")
	}
	// t.Error("---")
}
