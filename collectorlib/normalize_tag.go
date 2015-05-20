package collectorlib

// import (
// 	// "fmt"
// 	"regexp"
// 	// "strings"
// )

// var (
// 	Re_Trim, _ = regexp.Compile(`(^\s+|\s+$)`)

// 	ntag_pat_remove, _ = regexp.Compile(`[^\pL\d-_./\\]`)
// 	ntag_pat_s, _      = regexp.Compile(`\s+`)
// )

// func NormalizeTag(tag string) string {
// 	var tmp string

// 	tmp = Re_Trim.ReplaceAllString(tag, "")

// 	tmp = ntag_pat_s.ReplaceAllString(tmp, "_")
// 	// fmt.Printf("*%s*\n", tag)

// 	tmp = ntag_pat_remove.ReplaceAllString(tmp, "")
// 	// fmt.Printf("*%s*\n", tag)

// 	return tmp
// }

// func NormalizeTags(tags map[string]string) map[string]string {
// 	tmp := map[string]string{}
// 	for key, value := range tags {
// 		// tmp[NormalizeTag(key)] = NormalizeTag(value)
// 		nkey := NormalizeTag(key)
// 		nvalue := NormalizeTag(value)
// 		tmp[nkey] = nvalue

// 	}
// 	return tmp
// }
