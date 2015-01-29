package config

import (
	"testing"
)

func TestEvalTpl(t *testing.T) {
	str, err := evalTpl(`<format id="{{.Logger}}-{{.Level}}" format="{{.Format}}"/>`, struct {
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
		t.Errorf("eval: %s != %s", str, expected)
	}
}

func TestGenOutputs(t *testing.T) {
	expected := `<outputs><filter levels="debug"><console formatid="console-debug"/></filter><filter levels="info"><console formatid="console-info"/></filter><filter levels="warn"><console formatid="console-warn"/></filter><filter levels="error"><console formatid="console-error"/><rollingfile formatid="file-error" type="size" filename="/var/log/hickwall/hickwall.log" maxsize="300" maxrolls="5"/></filter><filter levels="critical"><console formatid="console-critical"/><rollingfile formatid="file-critical" type="size" filename="/var/log/hickwall/hickwall.log" maxsize="300" maxrolls="5"/></filter></outputs>`
	outputs, _ := gen_outputs("debug", "error", "/var/log/hickwall/hickwall.log", 300, 5)
	t.Log("expected: ", expected)

	if outputs != expected {
		t.Error(outputs)
	}
	// t.Error("-------------------")
}

func TestGenFormatsTplNoColor(t *testing.T) {
	// expected :=
	formats, err := gen_formats("console_format", "file_format")
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
