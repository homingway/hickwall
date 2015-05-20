package utils

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReadLine(t *testing.T) {
	ioutil.WriteFile("test.txt", []byte("data"), 0644)
	ReadLine("test.txt", func(line string) error {
		t.Log(line)
		if line != "data" {
			t.Error("failed")
		}
		return nil
	})
	os.Remove("test.txt")
}
