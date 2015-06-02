package logging

import (
	"fmt"
	"github.com/oliveagle/hickwall/logging/level"
	"log"
	"os"
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

func SetLevel(lvl level.LEVEL) error {
	switch lvl {
	case level.TRACE:
		_level = level.TRACE
	case level.DEBUG:
		_level = level.DEBUG
	case level.INFO:
		_level = level.INFO
	case level.WARNING:
		_level = level.WARNING
	case level.ERROR:
		_level = level.ERROR
	case level.CRITICAL:
		_level = level.CRITICAL
	default:
		return fmt.Errorf("invalid level: ", lvl)
	}
	return nil
}

func Trace(v ...interface{}) {
	if _level >= level.TRACE {
		log.SetPrefix("[TRACE] ")
		log.Println(v...)
	}
}

func Tracef(format string, v ...interface{}) {
	if _level >= level.TRACE {
		log.SetPrefix("[TRACE] ")
		log.Printf(format, v...)
	}
}

func Debug(v ...interface{}) {
	if _level >= level.DEBUG {
		log.SetPrefix("[DEBUG] ")
		log.Println(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if _level >= level.DEBUG {
		log.SetPrefix("[DEBUG] ")
		log.Printf(format, v...)
	}
}

func Info(v ...interface{}) {
	if _level >= level.INFO {
		log.SetPrefix("[INFO] ")
		log.Println(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if _level >= level.INFO {
		log.SetPrefix("[INFO] ")
		log.Printf(format, v...)
	}
}

func Error(v ...interface{}) {
	if _level >= level.ERROR {
		log.SetPrefix("[ERROR] ")
		log.Println(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if _level >= level.ERROR {
		log.SetPrefix("[ERROR] ")
		log.Printf(format, v...)
	}
}

func Critical(v ...interface{}) {
	if _level >= level.CRITICAL {
		log.SetPrefix("[CRITICAL] ")
		log.Println(v...)
	}
}

func Criticalf(format string, v ...interface{}) {
	if _level >= level.CRITICAL {
		log.SetPrefix("[CRITICAL] ")
		log.Printf(format, v...)
	}
}
