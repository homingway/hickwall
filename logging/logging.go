package logging

import (
	"fmt"
	"github.com/oliveagle/hickwall/logging/level"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	output *os.File
	mu     sync.Mutex
	_level level.LEVEL
)

func init() {
	_level = level.DEBUG
	create_output()
	log.SetFlags(log.Ldate | log.Ltime)
}

func create_output() {
	mu.Lock()
	defer mu.Unlock()
	var err error

	output, err = os.OpenFile(LOG_FILEPATH, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	log.SetOutput(output)
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
		log.SetPrefix("[TRACE] ")
		log.Println(v...)
		log.SetPrefix("[-] ")
	}
}

func Tracef(format string, v ...interface{}) {
	if _level <= level.TRACE {
		log.SetPrefix("[TRACE] ")
		log.Printf(format, v...)
		log.SetPrefix("[-] ")
	}
}

func Debug(v ...interface{}) {
	if _level <= level.DEBUG {
		log.SetPrefix("[DEBUG] ")
		log.Println(v...)
		log.SetPrefix("[-] ")
	}
}

func Debugf(format string, v ...interface{}) {
	if _level <= level.DEBUG {
		log.SetPrefix("[DEBUG] ")
		log.Printf(format, v...)
		log.SetPrefix("[-] ")
	}
}

func Info(v ...interface{}) {
	if _level <= level.INFO {
		log.SetPrefix("[INFO] ")
		log.Println(v...)
		log.SetPrefix("[-] ")
	}
}

func Infof(format string, v ...interface{}) {
	if _level <= level.INFO {
		log.SetPrefix("[INFO] ")
		log.Printf(format, v...)
		log.SetPrefix("[-] ")
	}
}

func Error(v ...interface{}) {
	if _level <= level.ERROR {
		log.SetPrefix("[ERROR] ")
		log.Println(v...)
		log.SetPrefix("[-] ")
	}
}

func Errorf(format string, v ...interface{}) {
	if _level <= level.ERROR {
		log.SetPrefix("[ERROR] ")
		log.Printf(format, v...)
		log.SetPrefix("[-] ")
	}
}

func Critical(v ...interface{}) {
	if _level <= level.CRITICAL {
		log.SetPrefix("[CRITICAL] ")
		log.Println(v...)
		log.SetPrefix("[-] ")
	}
}

func Criticalf(format string, v ...interface{}) {
	if _level <= level.CRITICAL {
		log.SetPrefix("[CRITICAL] ")
		log.Printf(format, v...)
		log.SetPrefix("[-] ")
	}
}
