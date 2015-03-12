package utils

import (
	// "fmt"
	"github.com/oliveagle/stringio"
	"regexp"
	"text/template"
)

var txt_tpl_pat_field, _ = regexp.Compile(`\{\{\.(\w+(([_|\.]?)+\w+)+)\}\}`)

func ExecuteTemplate(tpl string, data interface{}, post_process func(string) string) (string, error) {
	buf := stringio.NewStringIO()
	defer buf.Close()

	t, err := template.New("").Parse(tpl)
	if err != nil {
		return "", err
	}
	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	if post_process != nil {
		return post_process(buf.GetValueString()), nil
	}
	return buf.GetValueString(), nil
}

func FindAllTemplateKeys(tpl string) (res []string) {
	for _, k := range txt_tpl_pat_field.FindAllStringSubmatch(tpl, -1) {
		res = append(res, k[1])
	}
	return
}
