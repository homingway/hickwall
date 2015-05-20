package newcore

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

	exp := "win.wmi.fs.d.cdfs.freespace.bytes"
	res := NormalizeMetricKey(metric)
	t.Logf("*%s*", res)
	if res != exp {
		t.Error("failed")
	}
}

func Test_normalize_metric_key_3(t *testing.T) {
	// '-'' to '_'
	metric := "  win.wmi.fs.D:.CDFS.free-space.bytes"

	exp := "win.wmi.fs.d.cdfs.freespace.bytes"
	res := NormalizeMetricKey(metric)
	t.Logf("*%s*", res)
	if res != exp {
		t.Error("failed")
	}
}

func Test_normalize_metric_key_4(t *testing.T) {
	// numbers
	metric := " win.wmi.fs.D:.CDFS.free-space.bytes.112121"

	exp := "win.wmi.fs.d.cdfs.freespace.bytes.112121"
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

	exp := "win.wmi._path_to_file"
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

	exp := "win.wmi.c_path_to_file"
	res := NormalizeMetricKey(metric)
	t.Logf("*%s*", res)
	if res != exp {
		t.Error("failed")
	}
	// t.Error("---")
}

func Test_normalize_metric_key_from_tag_1(t *testing.T) {
	// path
	metric := "win.wmi.fs.D.CDFS.free_space.bytes_中文_1"

	exp := "win.wmi.fs.d.cdfs.free_space.bytes__1"
	res := NormalizeMetricKey(metric)
	// res = NormalizeMetricKey(res)
	t.Logf("*%s*", res)
	if res != exp {
		t.Error("failed")
	}
	// t.Error("---")
}

func Test_normalize_metric_key_from_tag_2(t *testing.T) {
	// path
	metric := "win.wmi.fs.D.CDFS.free-space.bytes_/path/to/file"

	exp := "win.wmi.fs.d.cdfs.freespace.bytes__path_to_file"
	res := NormalizeMetricKey(metric)
	// res = NormalizeMetricKey(res)
	t.Logf("*%s*", res)
	if res != exp {
		t.Error("failed")
	}
	// t.Error("---")
}

func Test_normalize_metric_key_from_tag_3(t *testing.T) {
	// path
	metric := "win.wmi.c\\path\\to\\file"

	exp := "win.wmi.c_path_to_file"
	res := NormalizeMetricKey(metric)
	// res = NormalizeMetricKey(res)
	t.Logf("*%s*", res)
	if res != exp {
		t.Error("failed")
	}
	// t.Error("---")
}
