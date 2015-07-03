package utils

import (
	"fmt"
	"testing"
)

func TestWal(t *testing.T) {

	wal, err := NewWal("D:\\gocodez\\src\\github.com\\oliveagle\\shared\\logs\\test\\test", 1, 5, true)
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
