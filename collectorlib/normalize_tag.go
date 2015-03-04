package collectorlib

import (
	// "fmt"
	"regexp"
	// "strings"
)

var (
	Re_Trim, _ = regexp.Compile(`(^\s+|\s+$)`)

	ntag_pat_remove, _ = regexp.Compile(`[^\pL\d-_./\\]`)
	ntag_pat_s, _      = regexp.Compile(`\s+`)
)

func NormalizeTag(tag string) string {

	tag = Re_Trim.ReplaceAllString(tag, "")

	tag = ntag_pat_s.ReplaceAllString(tag, "_")
	// fmt.Printf("*%s*\n", tag)

	tag = ntag_pat_remove.ReplaceAllString(tag, "")
	// fmt.Printf("*%s*\n", tag)

	return tag
}

func NormalizeTags(tags map[string]string) map[string]string {
	tmp := map[string]string{}
	for key, value := range tags {
		tmp[key] = NormalizeTag(value)
	}
	return tmp
}
