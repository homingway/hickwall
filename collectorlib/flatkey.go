package collectorlib

// import (
// 	"bytes"
// 	"fmt"
// 	"github.com/oliveagle/hickwall/utils"
// 	"sort"
// 	"strings"
// )

// type FlatTags map[string]string

// func (t FlatTags) String() string {
// 	var keys []string
// 	for k := range t {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)
// 	b := &bytes.Buffer{}
// 	for i, k := range keys {
// 		if i > 0 {
// 			fmt.Fprint(b, ",")
// 		}
// 		fmt.Fprintf(b, "%s_%s", k, t[k])
// 	}
// 	return NormalizeMetricKey(b.String())
// }

// func FlatMetricKeyAndTags(tpl, key string, tags map[string]string) (string, error) {
// 	if tpl == "" {
// 		return "", fmt.Errorf("template is empty")
// 	}

// 	keys := utils.FindAllTemplateKeys(tpl)
// 	if len(keys) == 0 {
// 		return "", fmt.Errorf("template don't have any substitutions: %s", tpl)
// 	}

// 	_tags := map[string]string{}
// 	for _, k := range keys {
// 		if strings.Count(k, ".") > 0 {
// 			return "", fmt.Errorf("We don't allowed multple leveled template. {{.Lv1.Lv2}}")
// 		}
// 		if strings.HasPrefix(k, "Tags_") == true {
// 			_key := strings.TrimLeft(k, "Tags_")
// 			_value, ok := tags[_key]
// 			if ok == true {
// 				_tags[k] = _value
// 			} else {
// 				return "", fmt.Errorf("tag is not found: %s", _key)
// 			}
// 			delete(tags, _key)
// 		}
// 	}
// 	// fmt.Println(_tags)
// 	// fmt.Println(tags)
// 	data := map[string]interface{}{
// 		"Tags": FlatTags(tags),
// 		"Key":  NormalizeMetricKey(key),
// 	}

// 	for k, v := range _tags {
// 		data[k] = v
// 	}

// 	res, err := utils.ExecuteTemplate(tpl, data, nil)
// 	if err != nil {
// 		return "", err
// 	}
// 	res = NormalizeMetricKey(res)
// 	return res, nil
// }
