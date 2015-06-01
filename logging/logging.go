package logging

import (
	"log"
	"os"
	"sync"
)

var (
	output *os.File
	mu     sync.Mutex
)

func init() {
	create_output()
}

func create_output() {
	mu.Lock()
	defer mu.Unlock()
	var err error

	filename := "d:\\hickwall.log"
	output, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	log.SetOutput(output)
}

func Println(v ...interface{}) {
	log.SetPrefix("")
	log.Println(v...)
}
func Printf(format string, v ...interface{}) {
	log.SetPrefix("")
	log.Printf(format, v...)
}

func Info(v ...interface{}) {
	log.SetPrefix("[INFO]")
	log.Println(v...)
}

func Infof(format string, v ...interface{}) {
	log.SetPrefix("[INFO]")
	log.Printf(format, v...)
}

func Error(v ...interface{}) {
	log.SetPrefix("[ERROR]")
	log.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	log.SetPrefix("[ERROR]")
	log.Printf(format, v...)
}
