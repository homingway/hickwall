package utils

import (
	// "fmt"
	"github.com/oliveagle/stringio"
	"text/template"
)

func ExecuteTemplate(tpl string, data map[string]string, post_process func(string) string) (string, error) {
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
