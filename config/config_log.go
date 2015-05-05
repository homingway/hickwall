package config

import (
	"fmt"
	log "github.com/oliveagle/seelog"
	"github.com/oliveagle/stringio"
	"os"
	// "path"
	"path/filepath"
	// "runtime"
	"strings"
	"text/template"
)

const (
	// LOG_FORMAT = "%Date(2006-01-02T15:04:05.00 MST) [%Level] %RelFile:%Line(%FuncShort) %Msg%n"
	LOG_FORMAT = "%Date(2006-01-02T15:04:05.00 MST) [%Level] %File:%Line(%FuncShort) %Msg%n"
)

var (
	Logger        log.LoggerInterface
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

func Mkdir_p_logdir(logfile string) {
	dir, _ := filepath.Split(logfile)
	if dir != "" {
		// fmt.Println("log dir: ", dir)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("Error: cannot create log dir: %s, err: %s", dir, err)
		}
	}
}

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
	// CoreConf.setDefaultByKey("Log_console_level", "info")
	// CoreConf.setDefaultByKey("Log_level", "debug")
	// CoreConf.setDefaultByKey("Log_file_filepath", LOG_FILEPATH)
	// CoreConf.setDefaultByKey("Log_file_maxsize", 300)
	// CoreConf.setDefaultByKey("Log_file_maxrolls", 5)

	core_viper.SetDefault("Log_file_maxrolls", 5)
	core_viper.SetDefault("Log_file_maxsize", 300)
	core_viper.SetDefault("Log_console_level", "info")
	core_viper.SetDefault("Log_level", "debug")
	core_viper.SetDefault("Log_file_filepath", LOG_FILEPATH)

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
			// if ALLOWED_COLOR_LOG && args.colored_console {
			if args.colored_console {
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

func ConfigLogger() error {
	var (
		maxsize   = 1
		maxrolls  = 1
		log_level = "debug"
	)

	setLoggerDefaults()

	formats_args := &gen_formats_args{
		fmt_console: LOG_FORMAT,
		fmt_file:    LOG_FORMAT,
	}

	if CoreConf.Log_file_maxsize > maxsize {
		maxsize = CoreConf.Log_file_maxsize
	}
	if CoreConf.Log_file_maxrolls > maxrolls {
		maxrolls = CoreConf.Log_file_maxrolls
	}

	_log_level := strings.ToLower(CoreConf.Log_level)

	switch _log_level {
	case "trace", "debug", "info", "error", "critical":
		log_level = _log_level
	}

	outputs_args := &gen_outputs_args{
		colored_console: false,
		console_level:   "info",
		file_level:      log_level,
		file_path:       LOG_FILEPATH,
		maxsize:         maxsize * 1024 * 1024,
		maxrolls:        maxrolls,
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
	// log.Debug(config_str)

	// log.Debug("log_filepath: ", LOG_FILEPATH)
	return nil
}
