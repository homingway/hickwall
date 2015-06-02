package config

import (
	"testing"
)

func TestEvalTpl(t *testing.T) {
	str, _ := evalTpl(`<format id="{{.Logger}}-{{.Level}}" format="{{.Format}}"/>`, struct {
		Format string
		Level  string
		Logger string
	}{
		Logger: "console",
		Level:  "debug",
		Format: "hahah",
	})
	// t.Log("eval_tpl\n\n", str, err)
	expected := `<format id="console-debug" format="hahah"/>`
	if str != expected {
		t.Errorf("eval: '%v' != '%v'", []byte(str), []byte(expected))
	}
}

func TestGenOutputs(t *testing.T) {
	expected := `<outputs><filter levels="debug"><console formatid="console-"/></filter><filter levels="info"><console formatid="console-"/></filter><filter levels="warn"><console formatid="console-"/></filter><filter levels="error"><console formatid="console-"/><rollingfile formatid="file-" type="size" filename="/var/log/hickwall/hickwall.log" maxsize="300" maxrolls="5"/></filter><filter levels="critical"><console formatid="console-"/><rollingfile formatid="file-" type="size" filename="/var/log/hickwall/hickwall.log" maxsize="300" maxrolls="5"/></filter></outputs>`
	args := &gen_outputs_args{
		colored_console: false,
		console_level:   "debug",
		file_level:      "error",
		file_path:       "/var/log/hickwall/hickwall.log",
		maxsize:         300,
		maxrolls:        5,
	}
	outputs, _ := gen_outputs(args)
	t.Log("expected: ", expected)

	if outputs != expected {
		t.Error(outputs)
	}
	// t.Error("-------------------")
}

func TestGenFormatsTplNoColor(t *testing.T) {
	// expected :=

	args := &gen_formats_args{
		fmt_console: "console_format",
		fmt_file:    "file_format",
	}
	formats, err := gen_formats(args)
	if err != nil {
		t.Error(err)
	}
	// t.Error(formats)
	if formats == "" {
		t.Errorf("Generated Leveled Formats is Empty")
	}
	head := "<formats>"
	if formats[:len(head)] != head {
		t.Errorf("Generated Leveled Formats Error")
	}
	// t.Error("-------------------")
}
