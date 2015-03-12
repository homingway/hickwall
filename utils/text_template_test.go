package utils

import (
	"testing"
)

func Test_text_template_FindAllTemplateKeys_1(t *testing.T) {
	tpl := "hotel.{{.Tags.Bu.Any_hahah.XX}}.{{.Host_name}}.{{.Key}}.{{.Tags}}"
	keys := FindAllTemplateKeys(tpl)
	exp := map[string]bool{
		"Tags.Bu.Any_hahah.XX": true,
		"Host_name":            true,
		"Key":                  true,
		"Tags":                 true,
	}
	if len(keys) != 4 {
		t.Error("sholud be 4")
	}

	for _, k := range keys {
		_, ok := exp[k]
		if ok != true {
			t.Error("key not find: ", k)
		}
	}
}
