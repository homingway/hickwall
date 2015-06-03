package logging

import (
	"io/ioutil"
	"log"
	"testing"
)

func BenchmarkBuiltinLogger_Debug(b *testing.B) {
	logger := log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
	for n := 0; n < b.N; n++ {
		logger.Println("this is debug")
	}
}

func BenchmarkNullLogger_Debug(b *testing.B) {
	initNullLogger()
	for n := 0; n < b.N; n++ {
		Debug("this is debug")
	}
}

func BenchmarkBuiltinLogger_Printf(b *testing.B) {
	logger := log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
	for n := 0; n < b.N; n++ {
		logger.Printf("this is debug: %d", n)
	}
}

func BenchmarkNullLogger_Debugf(b *testing.B) {
	initNullLogger()
	for n := 0; n < b.N; n++ {
		Debug("this is debug: %d", n)
	}
}
