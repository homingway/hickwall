package pdh

import (
	"fmt"
	// "log"
	"time"
)

var (
	_ = fmt.Sprint("")
)

type PdhCollectResult struct {
	Query     string
	Timestamp time.Time
	Value     float64
	Err       error
}

type PdhCollector struct {
	handle   uintptr
	counters map[string]uintptr
}

func NewPdhCollector() PdhCollector {
	// var handle uintptr
	// PdhOpenQuery(0, 0, &handle)
	_, handle := PdhOpenQuery_1(0, 0)

	return PdhCollector{
		handle:   handle,
		counters: make(map[string]uintptr),
	}
}

func (p PdhCollector) GetHandle() uintptr {
	return p.handle
}

func (p PdhCollector) Close() {

	keys := make([]string, 0, len(p.counters))
	for key, _ := range p.counters {
		keys = append(keys, key)
	}

	for _, key := range keys {
		// p.counters[key] = nil
		delete(p.counters, key)
	}

	// p.handle = nil

	PdhCloseQuery(p.handle)
}

func (p PdhCollector) AddEnglishCounter(query string) {
	// var handle uintptr
	// PdhAddEnglishCounter(p.handle, query, 0, &handle)
	_, handle := PdhAddEnglishCounter_1(p.handle, query, 0)
	p.counters[query] = handle
}

func valid_pdh_cstatus(cs uint32) bool {
	if cs == uint32(PDH_CSTATUS_VALID_DATA) || cs == uint32(PDH_CSTATUS_NEW_DATA) {
		return true
	}
	return false
}

func (p PdhCollector) CollectData() []*PdhCollectResult {
	PdhCollectQueryData(p.handle)
	data := []*PdhCollectResult{}

	var perf PDH_FMT_COUNTERVALUE_DOUBLE
	for key, chandle := range p.counters {
		cstatus := PdhValidatePath(key)
		if valid_pdh_cstatus(cstatus) == true {
			PdhGetFormattedCounterValueDouble(chandle, 0, &perf)
			if valid_pdh_cstatus(perf.CStatus) == true {
				pd := PdhCollectResult{
					Query:     key,
					Timestamp: time.Now(),
					Value:     perf.DoubleValue,
					Err:       nil,
				}
				data = append(data, &pd)
			} else {
				pd := PdhCollectResult{
					Query:     key,
					Timestamp: time.Now(),
					Value:     perf.DoubleValue,
					Err:       fmt.Errorf("invalid data: CSTATUS: 0x%X\n", perf.CStatus),
				}
				data = append(data, &pd)
				// log.Printf("invalid data: CSTATUS: 0x%X\n", perf.CStatus)
			}
		} else {
			pd := PdhCollectResult{
				Query:     key,
				Timestamp: time.Now(),
				Value:     perf.DoubleValue,
				Err:       fmt.Errorf("invalid path: CSTATUS: 0x%X Path: %s\n", cstatus, key),
			}
			data = append(data, &pd)
			// log.Printf("invalid path: CSTATUS: 0x%X Path: %s\n", cstatus, key)
		}
	}
	return data
}
