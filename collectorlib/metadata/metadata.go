package metadata

import (
	// "log"
	// "reflect"
	// "time"
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"strings"
	"sync"
)

// RateType is the type of rate for a metric: gauge, counter, or rate.
type RateType string

const (
	// Unknown is a not-yet documented rate type.
	Unknown RateType = ""
	// Gauge rate type.
	Gauge = "gauge"
	// Counter rate type.
	Counter = "counter"
	// Rate rate type.
	Rate = "rate"
)

// // Unit is the unit for a metric.
// type Unit string

// const (
// 	// None is a not-yet documented unit.
// 	None           Unit = ""
// 	A                   = "A" // Amps
// 	Bool                = "bool"
// 	BitsPerSecond       = "bits per second"
// 	Bytes               = "bytes"
// 	BytesPerSecond      = "bytes per second"
// 	C                   = "C" // Celsius
// 	CHz                 = "CentiHertz"
// 	Context             = "contexts"
// 	ContextSwitch       = "context switches"
// 	Count               = ""
// 	Document            = "documents"
// 	Entropy             = "entropy"
// 	Event               = ""
// 	Eviction            = "evictions"
// 	Fault               = "faults"
// 	Flush               = "flushes"
// 	Files               = "files"
// 	Get                 = "gets"
// 	GetExists           = "get exists"
// 	Interupt            = "interupts"
// 	KBytes              = "kbytes"
// 	Load                = "load"
// 	MHz                 = "MHz" // MegaHertz
// 	Megabit             = "Mbit"
// 	Merge               = "merges"
// 	MilliSecond         = "milliseconds"
// 	Ok                  = "ok" // "OK" or not status, 0 = ok, 1 = not ok
// 	Operation           = "Operations"
// 	Page                = "pages"
// 	Pct                 = "percent" // Range of 0-100.
// 	PerSecond           = "per second"
// 	Process             = "processes"
// 	Query               = "queries"
// 	Refresh             = "refreshes"
// 	Replica             = "replicas"
// 	RPM                 = "RPM" // Rotations per minute.
// 	Second              = "seconds"
// 	Segment             = "segments"
// 	Shard               = "shards"
// 	Socket              = "sockets"
// 	Suggest             = "suggests"
// 	StatusCode          = "status code"
// 	Syscall             = "system calls"
// 	V                   = "V" // Volts
// 	V10                 = "tenth-Volts"
// 	Watt                = "Watts"
// )

// Metakey uniquely identifies a metadata entry.
type Metakey struct {
	Metric string
	Tags   string
	Name   string
}

var (
	metadata  = make(map[Metakey]interface{})
	metalock  sync.Mutex
	metahost  string
	metafuncs []func()
	metadebug bool
)

// AddMeta adds a metadata entry to memory, which is queued for later sending.
func AddMeta(metric string, tags collectorlib.TagSet, name string, value interface{}, setHost bool) {
	// fmt.Printf("AddMeta: metric: %s, tags: %v, name: %s, value: %v, setHost: %v\n", metric, tags, name, value, setHost)
	// if tags == nil {
	// 	tags = make(TagSet)
	// }
	// if _, present := tags["host"]; setHost && !present {
	// 	//TODO: tags["host"]
	// 	// tags["host"] = util.Hostname
	// }
	// if err := tags.Clean(); err != nil {
	// 	//TODO: slog.Error
	// 	// slog.Error(err)
	// 	return
	// }
	// ts := tags.Tags()
	// metalock.Lock()
	// defer metalock.Unlock()
	// prev, present := metadata[Metakey{metric, ts, name}]
	// if present && !reflect.DeepEqual(prev, value) {
	// 	//TODO: slog.Infof
	// 	// slog.Infof("metadata changed for %s/%s/%s: %v to %v", metric, ts, name, prev, value)

	// 	//TODO: go sendMetadata
	// 	// go sendMetadata([]Metasend{{
	// 	// 	Metric: metric,
	// 	// 	Tags:   tags,
	// 	// 	Name:   name,
	// 	// 	Value:  value,
	// 	// }})
	// } else if metadebug {
	// 	//TODO slog.Infof
	// 	// slog.Infof("AddMeta for %s/%s/%s: %v", metric, ts, name, value)
	// }
	// metadata[Metakey{metric, ts, name}] = value
}

func ParseRateType(r string) (RateType, error) {
	rl := strings.ToLower(r)
	switch rl {
	case Gauge:
		return Gauge, nil
	case Counter:
		return Counter, nil
	case Rate:
		return Rate, nil
	default:
		return Unknown, fmt.Errorf("unknown rate type: %s", r)
	}
}
