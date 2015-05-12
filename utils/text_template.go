package utils

import (
	// "fmt"
	"regexp"
	"text/template"

	"bytes"
)

var txt_tpl_pat_field, _ = regexp.Compile(`\{\{\.(\w+(([_|\.]?)+\w+)+)\}\}`)

func ExecuteTemplate(tpl string, data interface{}, post_process func(string) string) (string, error) {
	var buf bytes.Buffer
	defer buf.Reset()

	t, err := template.New("").Parse(tpl)
	if err != nil {
		return "", err
	}
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	if post_process != nil {
		return post_process(buf.String()), nil
	}
	return buf.String(), nil
}

func FindAllTemplateKeys(tpl string) (res []string) {
	for _, k := range txt_tpl_pat_field.FindAllStringSubmatch(tpl, -1) {
		res = append(res, k[1])
	}
	return
}
