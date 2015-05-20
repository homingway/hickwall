package meta

import (
	"fmt"
	"strings"
)

// DataType is the type of rate for a metric: gauge, counter, or rate.
type DataType string

const (
	// Unknown is a not-yet documented data type.
	Unknown DataType = ""
	// Gauges are numbers that are expected to fluctuate over time.
	Gauge = "gauge"
	//Counters are numbers that increase over time, and never decrease (with the exception that the counter resets to zero).
	Counter = "counter"
	// A delta is the change from the previous data point.
	Delta = "delta"
)

// Unit is the unit for a metric.
type Unit string

const (
	// None is a not-yet documented unit.
	None           Unit = ""
	A                   = "A" // Amps
	Bool                = "bool"
	BitsPerSecond       = "bits per second"
	Bytes               = "bytes"
	BytesPerSecond      = "bytes per second"
	C                   = "C" // Celsius
	CHz                 = "CentiHertz"
	Context             = "contexts"
	ContextSwitch       = "context switches"
	Count               = ""
	Document            = "documents"
	Entropy             = "entropy"
	Event               = ""
	Eviction            = "evictions"
	Fault               = "faults"
	Flush               = "flushes"
	Files               = "files"
	Get                 = "gets"
	GetExists           = "get exists"
	Interupt            = "interupts"
	KBytes              = "kbytes"
	Load                = "load"
	MHz                 = "MHz" // MegaHertz
	Megabit             = "Mbit"
	Merge               = "merges"
	MilliSecond         = "milliseconds"
	Ok                  = "ok" // "OK" or not status, 0 = ok, 1 = not ok
	Operation           = "Operations"
	Page                = "pages"
	Pct                 = "percent" // Range of 0-100.
	PerSecond           = "per second"
	Process             = "processes"
	Query               = "queries"
	Refresh             = "refreshes"
	Replica             = "replicas"
	RPM                 = "RPM" // Rotations per minute.
	Second              = "seconds"
	Segment             = "segments"
	Shard               = "shards"
	Socket              = "sockets"
	Suggest             = "suggests"
	StatusCode          = "status code"
	Syscall             = "system calls"
	V                   = "V" // Volts
	V10                 = "tenth-Volts"
	Watt                = "Watts"
)

func ParseDataType(r string) (DataType, error) {
	r1 := strings.ToLower(r)
	switch r1 {
	case Gauge:
		return Gauge, nil
	case Counter:
		return Counter, nil
	case Delta:
		return Delta, nil
	default:
		return Unknown, fmt.Errorf("unknown data type: %s", r)
	}
}
