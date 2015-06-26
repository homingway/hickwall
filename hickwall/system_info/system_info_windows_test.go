package hickwall

import (
	"github.com/kr/pretty"
	"github.com/oliveagle/hickwall/newcore"
	"testing"
)

func Test_GetSystemInfo(t *testing.T) {
	newcore.SetHostname("hahah")

	res, err := GetSystemInfo()
	if err != nil {
		t.Error("...")
	}
	pretty.Println(res)
	if res.NumberOfProcessors <= 0 {
		t.Error("...")
	}

	if res.Name != "hahah" {
		t.Error("newcore SetHostname doesn't work here.")
	}
}
