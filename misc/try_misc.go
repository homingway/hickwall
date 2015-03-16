package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func main() {
	var buf bytes.Buffer

	tpl := "{{.Key}}"

	t, err := template.New("").Parse(tpl)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := map[string]interface{}{
		"Key": 111,
	}

	err = t.Execute(&buf, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(buf.String())
}
