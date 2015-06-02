package utils

import (
	// "fmt"
	"bytes"
	"regexp"
	"sync"
	"text/template"
)

var (
	txt_tpl_pat_field, _ = regexp.Compile(`\{\{\.(\w+(([_|\.]?)+\w+)+)\}\}`)
)

var tplBufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

//FIXME: leaking param: data
func ExecuteTemplate(tpl string, data interface{}, post_process func(string) string) (string, error) {
	buffer := tplBufPool.Get().(*bytes.Buffer)
	defer tplBufPool.Put(buffer)
	defer buffer.Reset()

	t, err := template.New("").Parse(tpl[:])
	if err != nil {
		return "", err
	}

	err = t.Execute(buffer, data)
	if err != nil {
		return "", err
	}
	if post_process != nil {
		return post_process(buffer.String()), nil
	}
	return buffer.String(), nil
}

func FindAllTemplateKeys(tpl string) (res []string) {
	for _, k := range txt_tpl_pat_field.FindAllStringSubmatch(tpl[:], -1) {
		res = append(res, k[1])
	}
	return
}
