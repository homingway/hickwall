package logging

import (
	"fmt"
	"github.com/oliveagle/hickwall/logging/level"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	output *os.File
	mu     sync.Mutex
	_level level.LEVEL
	logger *log.Logger
)

// var ppFree = sync.Pool{
// 	New: func() interface{} { return new(pp) },
// }

func InitFileLogger(log_filepath string) {
	_level = level.DEBUG
	writer, err := create_output(log_filepath[:])
	if err != nil {
		panic(err)
	}
	logger = log.New(writer, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func InitStdoutLogger() {
	_level = level.DEBUG
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func initNullLogger() {
	_level = level.DEBUG
	logger = log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func Close() error {
	if output != nil {
		return output.Close()
	}
	return nil
}

func create_output(log_filepath string) (writer io.Writer, err error) {
	mu.Lock()
	defer mu.Unlock()

	if output != nil {
		output.Close()
		output = nil
	}

	output, err = os.OpenFile(log_filepath[:], os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	writer = io.MultiWriter(os.Stdout, output)
	return
}

func SetLevel(lvl string) error {
	switch strings.ToLower(lvl[:]) {
	case "trace":
		_level = level.TRACE
		Debug("set logging level to TRACE")
	case "debug":
		_level = level.DEBUG
		Debug("set logging level to DEBUG")
	case "info":
		_level = level.INFO
	case "warn":
		_level = level.WARNING
	case "error":
		_level = level.ERROR
	case "critical":
		_level = level.CRITICAL
	default:
		return fmt.Errorf("invalid level: ", lvl)
	}
	return nil
}

func Trace(v ...interface{}) {
	if _level <= level.TRACE {
		logger.Output(2, fmt.Sprintf("[TRACE] %s", v...))
	}
}

func Tracef(format string, v ...interface{}) {
	if _level <= level.TRACE {
		logger.Output(2, fmt.Sprintf("[TRACE] %s", fmt.Sprintf(format, v...)))
	}
}

func Debug(v ...interface{}) {
	if _level <= level.DEBUG {
		logger.Output(2, fmt.Sprintf("[DEBUG] %s", v...)) // 1728 ns/op
	}
}

func Debugf(format string, v ...interface{}) {
	if _level <= level.DEBUG {
		logger.Output(2, fmt.Sprintf("[DEBUG] %s", fmt.Sprintf(format, v...))) //2046 ns/op
	}
}

func Info(v ...interface{}) {
	if _level <= level.INFO {
		logger.Output(2, fmt.Sprintf("[INFO] %s", v...))
	}
}

func Infof(format string, v ...interface{}) {
	if _level <= level.INFO {
		logger.Output(2, fmt.Sprintf("[INFO] %s", fmt.Sprintf(format, v...)))
	}
}

func Error(v ...interface{}) {
	if _level <= level.ERROR {
		logger.Output(2, fmt.Sprintf("[ERROR] %s", v...))
	}
}

func Errorf(format string, v ...interface{}) {
	if _level <= level.ERROR {
		logger.Output(2, fmt.Sprintf("[ERROR] %s", fmt.Sprintf(format, v...)))
	}
}

func Critical(v ...interface{}) {
	if _level <= level.CRITICAL {
		logger.Output(2, fmt.Sprintf("[CRITICAL] %s", v...))
	}
}

func Criticalf(format string, v ...interface{}) {
	if _level <= level.CRITICAL {
		logger.Output(2, fmt.Sprintf("[CRITICAL] %s", fmt.Sprintf(format, v...)))
	}
}
