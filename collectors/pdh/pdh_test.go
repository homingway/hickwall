package pdh

import (
	"github.com/kr/pretty"
	// "fmt"
	"testing"
	// "time"
	// "unsafe"
	"math/rand"
)

func Test_simple(t *testing.T) {
	var handle uintptr
	var counterHandle uintptr
	ret := PdhOpenQuery(0, 0, &handle)
	// ret = PdhAddEnglishCounter(handle, "\\Processor(_Total)\\% Idle Time", 0, &counterHandle)
	// ret = PdhAddEnglishCounter(handle, "\\Memory\\% Committed Bytes In Use", 0, &counterHandle)
	// ret = PdhAddEnglishCounter(handle, "\\Memory\\Page Faults/sec", 0, &counterHandle)
	// ret = PdhAddEnglishCounter(handle, "\\System\\Processor Queue Length", 0, &counterHandle)
	// ret = PdhAddEnglishCounter(handle, "\\LogicalDisk(C:)\\% Free Space", 0, &counterHandle)
	ret = PdhAddEnglishCounter(handle, "\\System\\Processes", 0, &counterHandle)

	ret = PdhCollectQueryData(handle)

	// ret = PdhCollectQueryData(handle)
	// fmt.Printf("Collect return code is %x\n", ret) // return code will be ERROR_SUCCESS
	var perf PDH_FMT_COUNTERVALUE_DOUBLE
	ret = PdhGetFormattedCounterValueDouble(counterHandle, 0, &perf)
	t.Logf("ret: %v, perf: %v", ret, perf) // return code will be ERROR_SUCCESS
	// pretty.Println(perf)

	if perf.DoubleValue <= 0 {
		t.Error("perf.DoubleValue <= 0")
	}

	// t.Error("hhh")
}

func Test_simple_3(t *testing.T) {
	var handle uintptr
	// var counterHandle uintptr
	ret := PdhOpenQuery(0, 0, &handle)

	cHandles := make([]uintptr, 3)

	ret = PdhAddEnglishCounter(handle, "\\System\\Processes", 0, &cHandles[0])
	ret = PdhAddEnglishCounter(handle, "\\LogicalDisk(C:)\\% Free Space", 0, &cHandles[1])
	ret = PdhAddEnglishCounter(handle, "\\Memory\\Available MBytes", 0, &cHandles[2])

	ret = PdhCollectQueryData(handle)

	// ret = PdhCollectQueryData(handle)
	// fmt.Printf("Collect return code is %x\n", ret) // return code will be ERROR_SUCCESS
	var perf PDH_FMT_COUNTERVALUE_DOUBLE
	for i := 0; i < 3; i++ {
		ret = PdhGetFormattedCounterValueDouble(cHandles[i], 0, &perf)
		t.Logf("ret: %v, perf: %v", ret, perf) // return code will be ERROR_SUCCESS
		// pretty.Println(perf)

		if perf.DoubleValue <= 0 {
			t.Error("perf.DoubleValue <= 0")
		}
	}

	// t.Error("hhh")
}

func Test_wrapper(t *testing.T) {
	collector := NewPdhCollector()
	defer collector.Close()

	// collector.AddEnglishCounter("\\System\\Processes")
	// collector.AddEnglishCounter("\\LogicalDisk(C:)\\% Free Space")
	collector.AddEnglishCounter("\\Memory\\Available Bytes")
	collector.AddEnglishCounter("\\Processes(_Total)\\Working Set")
	collector.AddEnglishCounter("\\Memory\\Cache Bytes")

	// t.Log(collector.CollectData())
	pretty.Println(collector.CollectData())

	data := collector.CollectData()

	// pretty.Println(data[0].Value / 1024 / 1024)
	if len(data) != 3 {
		pretty.Println(data)
		t.Error("Wrapper error")
	}

	// pretty.Println((data[0].Value + data[1].Value + data[2].Value) / 1024 / 1024)

	// t.Error("--")
}

func Benchmark_200(b *testing.B) {
	var handle uintptr
	// var counterHandle uintptr
	PdhOpenQuery(0, 0, &handle)

	cHandles := make([]uintptr, 200)
	for i := 0; i < 200; i++ {
		PdhAddEnglishCounter(handle, "\\System\\Processes", 0, &cHandles[i])
	}
	PdhCollectQueryData(handle)

	var perf PDH_FMT_COUNTERVALUE_DOUBLE
	for i := 0; i < b.N; i++ {
		PdhCollectQueryData(handle)

		PdhGetFormattedCounterValueDouble(cHandles[rand.Intn(200)], 0, &perf)
		// t.Logf("ret: %v, perf: %v", ret, perf) // return code will be ERROR_SUCCESS
		// pretty.Println(perf)

		if perf.DoubleValue <= 0 {
			b.Error("perf.DoubleValue <= 0")
		}
	}

	// b.Error("hhh")
}
