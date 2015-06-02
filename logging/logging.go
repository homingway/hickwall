package logging

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type LEVEL int

const (
	TRACE    LEVEL = 0
	DEBUG          = 10
	INFO           = 20
	WARNING        = 30
	ERROR          = 40
	CRITICAL       = 50
)

var (
	output *os.File
	mu     sync.Mutex
	level  LEVEL
)

func init() {
	level = DEBUG
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

func SetLevel(lvl LEVEL) error {
	switch lvl {
	case TRACE:
		level = TRACE
	case DEBUG:
		level = DEBUG
	case INFO:
		level = INFO
	case WARNING:
		level = WARNING
	case ERROR:
		level = ERROR
	case CRITICAL:
		level = CRITICAL
	default:
		return fmt.Errorf("invalid level: ", lvl)
	}
	return nil
}

//func Println(v ...interface{}) {
//	log.SetPrefix("")
//	log.Println(v...)
//}
//func Printf(format string, v ...interface{}) {
//	log.SetPrefix("")
//	log.Printf(format, v...)
//}

func Trace(v ...interface{}) {
	if level >= TRACE {
		log.SetPrefix("[TRACE] ")
		log.Println(v...)
	}
}

func Tracef(format string, v ...interface{}) {
	if level >= TRACE {
		log.SetPrefix("[TRACE] ")
		log.Printf(format, v...)
	}
}

func Debug(v ...interface{}) {
	if level >= DEBUG {
		log.SetPrefix("[DEBUG] ")
		log.Println(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if level >= DEBUG {
		log.SetPrefix("[DEBUG] ")
		log.Printf(format, v...)
	}
}

func Info(v ...interface{}) {
	if level >= INFO {
		log.SetPrefix("[INFO] ")
		log.Println(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if level >= INFO {
		log.SetPrefix("[INFO] ")
		log.Printf(format, v...)
	}
}

func Error(v ...interface{}) {
	if level >= ERROR {
		log.SetPrefix("[ERROR] ")
		log.Println(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if level >= ERROR {
		log.SetPrefix("[ERROR] ")
		log.Printf(format, v...)
	}
}

func Critical(v ...interface{}) {
	if level >= CRITICAL {
		log.SetPrefix("[CRITICAL] ")
		log.Println(v...)
	}
}

func Criticalf(format string, v ...interface{}) {
	if level >= CRITICAL {
		log.SetPrefix("[CRITICAL] ")
		log.Printf(format, v...)
	}
}
