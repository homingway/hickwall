package pdh

import (
	"testing"
)

func TestPdh(t *testing.T) {
	pc := NewPdhCollector()
	pc.AddEnglishCounter("\\System\\Processes")
	data := pc.CollectData()
	defer pc.Close()

	for _, d := range data {
		t.Log(d.Value)
		if d.Err != nil || d.Value < 10 {
			t.Error("processes count less than 10 ?", d.Err)
		}
	}
}

func TestPdhAddInvalidCounterPath(t *testing.T) {
	pc := NewPdhCollector()
	pc.AddEnglishCounter("\\Systemhjahahah\\Processes")
	data := pc.CollectData()
	defer pc.Close()

	for _, d := range data {
		t.Log(d.Value)
		if d.Err == nil || d.Value > 0 {
			t.Error("should emit error")
		}
	}
}
