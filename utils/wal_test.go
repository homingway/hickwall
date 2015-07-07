package utils

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestWal(t *testing.T) {

	f, _ := ioutil.TempFile("", "test")
	f.Close()

	wal, err := NewWal(f.Name(), 1, 5, true)
	if err != nil {
		t.Error("err:%v", err)
	}

	wal.WriteLine("123123123")
	wal.WriteLine("4567845678")
	d, err := wal.ReadLine()

	fmt.Println(d)
	if err != nil {
		t.Error("err:%v", err)
	}

	err = wal.Commit()
	if err != nil {
		t.Error("err:%v", err)
	}

}
