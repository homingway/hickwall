package config

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/oliveagle/stringio"
	"os"
	"path/filepath"
	"text/template"
)

func Mkdir_p_logdir(logfile string) {
	dir, _ := filepath.Split(logfile)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Println("Error: cannot create log dir: %s, err: %s", dir, err)
	}
}

var (
	ordered_level = []string{
		"trace",
		"debug",
		"info",
		"warn",
		"error",
		"critical",
	}
	log_color_tpl_map = map[string]string{
		"trace":    "{{.Format}}",
		"debug":    "%EscM(1;30){{.Format}}%EscM(0)",
		"info":     "%EscM(1;34){{.Format}}%EscM(0)",
		"warn":     "%EscM(33){{.Format}}%EscM(0)",
		"error":    "%EscM(31){{.Format}}%EscM(0)",
		"critical": "%EscM(35){{.Format}}%EscM(0)",
	}
)

func evalTpl(tpl string, data interface{}) (str string, err error) {
	sio := stringio.NewStringIO()
	defer sio.Close()

	t := template.Must(template.New("tmp").Parse(tpl))
	err = t.Execute(sio, data)
	if err != nil {
		fmt.Printf("Error to eval logging tpl err: %v", err)
		return "", err
	}
	return sio.GetValueString(), err
}

func setLoggerDefaults() {
	Conf.setDefaultByKey("Log_colored_console", false)
	Conf.setDefaultByKey("Log_console_level", "debug")
	Conf.setDefaultByKey("Log_console_format", "%Date/%Time [%LEV] %Msg%n")

	Conf.setDefaultByKey("Log_file_level", "debug")
	Conf.setDefaultByKey("Log_file_filepath", "/var/log/hickwall/hickwall.log")
	Conf.setDefaultByKey("Log_file_format", "%Date(2006 Jan 02/3:04:05.000000000 PM MST) [%Level] %File %FullPath %RelFile %Msg%n")
	Conf.setDefaultByKey("Log_file_maxsize", 300)
	Conf.setDefaultByKey("Log_file_maxrolls", 5)
}

type gen_outputs_args struct {
	colored_console bool
	console_level   string
	file_level      string
	file_path       string
	maxsize         int
	maxrolls        int
}

func gen_outputs(args *gen_outputs_args) (str string, err error) {
	idx_console_level := len(ordered_level) + 1
	idx_file_level := len(ordered_level) + 1
	outputs := "<outputs>"
	for idx, level := range ordered_level {
		var (
			tmp_filter_console = ""
			tmp_filter_file    = ""
		)

		// t.Log(idx, level)
		if args.console_level == level {
			idx_console_level = idx
		}
		if args.file_level == level {
			idx_file_level = idx
		}

		tpl_log_out_console := `<console formatid="console-{{.Level}}"/>`
		// console logger can enable or disable colored output
		if idx >= idx_console_level {
			if ALLOWED_COLOR_LOG && args.colored_console {
				tmp_filter_console, err = evalTpl(tpl_log_out_console, struct {
					Level string
				}{
					Level: level,
				})
			} else {
				tmp_filter_console, err = evalTpl(tpl_log_out_console, struct {
					Level string
				}{
					Level: "",
				})
			}

			if err != nil {
				fmt.Println("gen_filters: Error to eval template", err)
				return "", err
			}
		}

		tpl_log_out_file := `<rollingfile formatid="file-" type="size" filename="{{.Filepath}}" maxsize="{{.Maxsize}}" maxrolls="{{.Maxrolls}}"/>`
		// file logger don't allow color code
		if idx >= idx_file_level {
			tmp_filter_file, err = evalTpl(tpl_log_out_file, struct {
				// Level    string
				Filepath string
				Maxsize  int
				Maxrolls int
			}{
				// Level:    level,
				Filepath: args.file_path,
				Maxsize:  args.maxsize,
				Maxrolls: args.maxrolls,
			})
			if err != nil {
				fmt.Println("gen_filters: Error to eval template", err)
				return "", err
			}
		}

		tpl_log_filter := `<filter levels="{{.Level}}">{{.ConsoleOut}}{{.FileOut}}</filter>`
		if tmp_filter_console != "" || tmp_filter_file != "" {
			filter, err := evalTpl(tpl_log_filter, struct {
				Level, ConsoleOut, FileOut string
			}{
				Level:      level,
				ConsoleOut: tmp_filter_console,
				FileOut:    tmp_filter_file,
			})
			if err != nil {
				fmt.Println("gen_filters: Error to eval template", err)
				return "", err
			}
			outputs += filter
		}
	}
	outputs += "</outputs>"
	return outputs, nil
}

type gen_formats_args struct {
	fmt_console string
	fmt_file    string
}

func gen_formats(args *gen_formats_args) (str string, err error) {
	tpl_log_fmt_color := `<format id="console-{{.Level}}" format="{{.ConsoleFormat}}"/><format id="file-{{.Level}}" format="{{.FileFormat}}"/>`
	type struct_log_fmt_color struct {
		Level, ConsoleFormat, FileFormat string
	}
	formats := ""

	// generate no colored formats for console and file
	no_color_format, err := evalTpl(tpl_log_fmt_color, struct_log_fmt_color{
		Level:         "",
		ConsoleFormat: args.fmt_console,
		FileFormat:    args.fmt_file,
	})
	if err != nil {
		fmt.Println("Error: cannot eval color format template: ", err)
		return "", err
	}
	formats += no_color_format

	// colored formats for console and file
	for level, tpl_color := range log_color_tpl_map {
		// fmt.Println(key, value)
		ConsoleFormat, err := evalTpl(tpl_color, struct {
			Format string
		}{
			Format: args.fmt_console,
		})
		if err != nil {
			fmt.Println("Error: cannot eval color format template: ", err)
			return "", err
		}
		FileFormat, err := evalTpl(tpl_color, struct {
			Format string
		}{
			Format: args.fmt_file,
		})
		if err != nil {
			fmt.Println("Error: cannot eval color format template: ", err)
			return "", err
		}

		format, err := evalTpl(tpl_log_fmt_color, struct_log_fmt_color{
			Level: level, ConsoleFormat: ConsoleFormat, FileFormat: FileFormat,
		})
		if err != nil {
			return "", err
		}
		formats += format
	}
	return "<formats>" + formats + "</formats>", nil
}

func gen_config(formats_args *gen_formats_args, outputs_args *gen_outputs_args) (config_str string, err error) {
	tpl_log_conf := `<seelog type="asynctimer" asyncinterval="100000" minlevel="trace" maxlevel="critical">{{.Outputs}}{{.Formats}}</seelog>`
	outputs, err := gen_outputs(outputs_args)
	if err != nil {
		fmt.Println("Error: failed to generate outputs ", outputs)
		return "", err
	}
	formats, err := gen_formats(formats_args)
	if err != nil {
		fmt.Println("Error: failed to generate formats ", formats)
		return "", err
	}
	config_str, err = evalTpl(tpl_log_conf, struct {
		Outputs, Formats string
	}{
		Outputs: outputs,
		Formats: formats,
	})
	if err != nil {
		fmt.Println("Error: failed to generate config ", err)
		return "", err
	}
	return config_str, nil
}

var Logger log.LoggerInterface

func ConfigLogger() error {
	setLoggerDefaults()

	Mkdir_p_logdir(Conf.Log_file_filepath)

	formats_args := &gen_formats_args{
		fmt_console: Conf.Log_console_format,
		fmt_file:    Conf.Log_file_format,
	}

	outputs_args := &gen_outputs_args{
		colored_console: Conf.Log_colored_console,
		console_level:   Conf.Log_console_level,
		file_level:      Conf.Log_file_level,
		file_path:       Conf.Log_file_filepath,
		// maxsize:         Conf.Log_file_maxsize * 1024 * 1024,
		maxsize:  Conf.Log_file_maxsize * 1024,
		maxrolls: Conf.Log_file_maxrolls,
	}

	config_str, err := gen_config(formats_args, outputs_args)
	if err != nil {
		fmt.Println("Error: failed to generate config ", err)
		return err
	}
	// fmt.Println(config_str)
	Logger, err = log.LoggerFromConfigAsString(config_str)
	if err != nil {
		fmt.Println("Error: cannot load log config from string: ", err)
		return err
	}
	log.ReplaceLogger(Logger)
	log.Info(config_str)
	return nil
}
