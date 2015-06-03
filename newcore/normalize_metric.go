package newcore

import (
	// "bytes"
	// "fmt"
	"regexp"
	"strings"
	// "unicode"
	// "unicode/utf8"
)

var (
	nmetric_pat_s, _      = regexp.Compile(`\s+|-+|/+|\\+|_+|,|;`)
	nmetric_pat_remove, _ = regexp.Compile(`[^\w-_.]`)
)

func NormalizeMetricKey(metric string) string {
	var tmp string

	tmp = Re_Trim.ReplaceAllString(metric[:], "")

	tmp = nmetric_pat_s.ReplaceAllString(tmp[:], "_")

	tmp = nmetric_pat_remove.ReplaceAllString(tmp[:], "")

	tmp = nmetric_pat_s.ReplaceAllString(tmp[:], "_")

	tmp = strings.Trim(tmp[:], "_")

	tmp = strings.ToLower(tmp[:])

	// this line will coz heap growing slowly.
	// tmp = strings.ToLower(strings.Trim(tmp, "_"))

	// fmt.Printf("*%s*\n", tmp)
	return tmp
}

// func NormalizeMetricKey(metric string) string {
// 	var tmp = string(metric)

// 	res, _ := Replace2_1(tmp, "/", "_")
// 	res, _ = Replace2_1(res, "\\", "_")

// 	// fmt.Println(res)

// 	res, _ = Replace_1(res, "")
// 	// res, _ = Replace(res, "")
// 	res = strings.ToLower(res)
// 	return res
// }

// func Replace2_1(s, target, replacement string) (string, error) {
// 	var tmp = string(s)

// 	var b bytes.Buffer
// 	var c string

// 	// defer b.Reset()

// 	replaced := false
// 	for len(tmp) > 0 {
// 		r, size := utf8.DecodeRuneInString(tmp)
// 		if string(r) == target && !replaced {
// 			b.WriteString(replacement)
// 			replaced = true
// 		} else {
// 			b.WriteString(string(r))
// 			replaced = false
// 		}
// 		tmp = tmp[size:]
// 	}
// 	c = b.String()

// 	if len(c) == 0 {
// 		return "", fmt.Errorf("clean result is empty")
// 	}
// 	b.Reset()

// 	return c, nil
// }

// func Replace2(s, target, replacement string) (string, error) {
// 	var c string
// 	replaced := false
// 	for len(s) > 0 {
// 		r, size := utf8.DecodeRuneInString(s)
// 		if string(r) == target && !replaced {
// 			c += replacement
// 			replaced = true
// 		} else {
// 			c += string(r)
// 			replaced = false
// 		}
// 		s = s[size:]
// 	}

// 	if len(c) == 0 {
// 		return "", fmt.Errorf("clean result is empty")
// 	}
// 	return c, nil
// }

// func Replace_1(s, replacement string) (string, error) {
// 	var tmp = s
// 	var b bytes.Buffer
// 	var c string

// 	// defer b.Reset()

// 	replaced := false
// 	for len(tmp) > 0 {
// 		r, size := utf8.DecodeRuneInString(tmp)
// 		switch {
// 		case 'a' <= r && r <= 'z':
// 			b.WriteString(string(r))
// 			replaced = false
// 		case 'A' <= r && r <= 'Z':
// 			b.WriteString(string(r))
// 			replaced = false
// 		case unicode.IsDigit(r) || r == '_' || r == '.':
// 			b.WriteString(string(r))
// 			replaced = false
// 		default:
// 			if !replaced {
// 				b.WriteString(replacement)
// 				replaced = true
// 			}
// 		}
// 		tmp = tmp[size:]
// 	}

// 	c = b.String()
// 	if len(c) == 0 {
// 		return "", fmt.Errorf("clean result is empty")
// 	}
// 	b.Reset()

// 	return c, nil
// }

// func Replace(s, replacement string) (string, error) {
// 	var c string
// 	replaced := false
// 	for len(s) > 0 {
// 		r, size := utf8.DecodeRuneInString(s)
// 		switch {
// 		case 'a' <= r && r <= 'z':
// 			c += string(r)
// 			replaced = false
// 		case 'A' <= r && r <= 'Z':
// 			c += string(r)
// 			replaced = false
// 		case unicode.IsDigit(r) || r == '_' || r == '.':
// 			c += string(r)
// 			replaced = false
// 		default:
// 			if !replaced {
// 				c += replacement
// 				replaced = true
// 			}
// 		}
// 		s = s[size:]
// 	}
// 	if len(c) == 0 {
// 		return "", fmt.Errorf("clean result is empty")
// 	}
// 	return c, nil
// }
