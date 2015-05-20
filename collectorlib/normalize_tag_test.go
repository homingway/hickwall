package collectorlib

// import (
// 	"testing"
// )

// func Test_normalize_tag_1(t *testing.T) {
// 	tag := "win.wmi.fs.d.cdfs.free_space.bytes"

// 	exp := "win.wmi.fs.d.cdfs.free_space.bytes"
// 	res := NormalizeTag(tag)
// 	t.Logf("*%s*", res)
// 	if res != exp {
// 		t.Error("failed")
// 	}
// }

// func Test_normalize_tag_2(t *testing.T) {
// 	// 中文ok，spaces will be transfer to "_"
// 	tag := "  win.wmi.fs.D:.CDFS.free space.bytes 中文 1 \r\n"

// 	exp := "win.wmi.fs.D.CDFS.free_space.bytes_中文_1"
// 	res := NormalizeTag(tag)
// 	t.Logf("*%s*", res)
// 	if res != exp {
// 		t.Error("failed")
// 	}
// 	// t.Error("1")
// }

// func Test_normalize_tag_3(t *testing.T) {
// 	tag := "  win.wmi.fs.D:.CDFS.free-space.bytes /path/to/file\r\n"

// 	exp := "win.wmi.fs.D.CDFS.free-space.bytes_/path/to/file"
// 	res := NormalizeTag(tag)
// 	t.Logf("*%s*", res)
// 	if res != exp {
// 		t.Error("failed")
// 	}
// 	// t.Error("1")
// }

// func Test_normalize_tag_4(t *testing.T) {
// 	// path
// 	tag := " win.wmi. /path/to/file"

// 	exp := "win.wmi._/path/to/file"
// 	res := NormalizeTag(tag)
// 	t.Logf("*%s*", res)
// 	if res != exp {
// 		t.Error("failed")
// 	}
// 	// t.Error("---")
// }

// func Test_normalize_tag_5(t *testing.T) {
// 	// path
// 	tag := " win.wmi. c:\\path\\to\\file"

// 	exp := "win.wmi._c\\path\\to\\file"
// 	res := NormalizeTag(tag)
// 	t.Logf("*%s*", res)
// 	if res != exp {
// 		t.Error("failed")
// 	}
// 	// t.Error("---")
// }

// func Test_normalize_tag_6(t *testing.T) {
// 	// path
// 	tag := ""

// 	exp := ""
// 	res := NormalizeTag(tag)
// 	t.Logf("*%s*", res)
// 	if res != exp {
// 		t.Error("failed")
// 	}
// 	// t.Error("---")
// }

// func Test_normalize_tags(t *testing.T) {
// 	// path
// 	tags := map[string]string{
// 		"t1": " win.wmi. c:\\path\\to\\file",
// 		"t2": " win.wmi. c:\\path\\to\\file",
// 	}

// 	res := NormalizeTags(tags)
// 	t.Logf("*%s*", res)
// 	if len(res) != 2 {
// 		t.Error("failed")
// 	}
// 	// t.Error("---")
// }

// func Test_normalize_tags_1(t *testing.T) {
// 	// path
// 	tags := map[string]string{
// 		" win.wmi. c:\\path\\to\\file": " win.wmi. c:\\path\\to\\file",
// 		"t2": " win.wmi. c:\\path\\to\\file",
// 	}

// 	res := NormalizeTags(tags)
// 	t.Logf("*%s*", res)
// 	if len(res) != 2 {
// 		t.Error("failed")
// 	}
// 	// t.Error("---")
// }
