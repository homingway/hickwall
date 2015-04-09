package backends

import (
	// "fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	// log "github.com/oliveagle/seelog"
	"time"
)

type FileWriter struct {
	tick <-chan time.Time
	mdCh chan collectorlib.MultiDataPoint
	conf config.Transport_file
}

// func NewFileWriter(conf config.Transport_file) *FileWriter {

// }
