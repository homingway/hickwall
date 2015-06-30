// write MultiDataPoint into ES.

package backends

import (
	"fmt"
	// "github.com/mozillazg/request"
	"github.com/franela/goreq"
	"github.com/oliveagle/hickwall/backends/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	// "net/http"
	"time"
)

var (
	_ = time.Now()
	_ = fmt.Sprintf("")
)

type elasticsearchBackend struct {
	name    string
	closing chan chan error             // for Close
	updates chan newcore.MultiDataPoint // for receive updates

	// elasticsearch backend specific attributes
	conf *config.Transport_elasticsearch

	// output *request.Request
}

func NewElasticsearchBackend(name string, conf *config.Transport_elasticsearch) (newcore.Publication, error) {
	s := &elasticsearchBackend{
		name:    name,
		closing: make(chan chan error),
		updates: make(chan newcore.MultiDataPoint),
		conf:    conf,
	}

	go s.loop()
	return s, nil
}

// func (b *elasticsearchBackend) ElasticsearchClient() error {
// 	c := new(http.Client)
// 	req := request.NewRequest(c)
// 	b.output = req
// 	return nil
// }

func (b *elasticsearchBackend) loop() {
	var (
		startConsuming <-chan newcore.MultiDataPoint
		// try_create_client_once chan bool
		// try_create_client_tick <-chan time.Time
	)
	startConsuming = b.updates
	logging.Info("elasticsearch backend loop started")

	for {
		// if b.output == nil && try_create_client_once == nil && try_create_client_tick == nil {
		// 	startConsuming = nil // disable consuming
		// 	try_create_client_once = make(chan bool)
		// 	// try to create elasticsearch client the first time async.
		// 	go func() {
		// 		err := b.ElasticsearchClient()
		// 		if err == nil {
		// 			try_create_client_once <- true
		// 		} else {
		// 			try_create_client_once <- false
		// 		}
		// 	}()
		// }

		select {
		case md := <-startConsuming:
			// if b.output != nil {
			logging.Tracef("elasticsearch backend consuming: 0x%X", &md)
			url := b.conf.URL + "/" + b.conf.Index + "/" + b.conf.Type
			for _, p := range md {
				data := map[string]interface{}{
					"metric":     p.Metric.Clean(),
					"@timestamp": p.Timestamp.Format(time.RFC3339Nano),
					"value":      p.Value,
					"tags":       p.Tags,
				}
				// resp, err := b.output.Post(url)
				resp, err := goreq.Request{
					Method: "POST",
					Uri:    url,
					Body:   data,
				}.Do()
				if err != nil {
					logging.Critical("post data to elasticsearch fail: %s", err)
				} else {
					resp.Body.Close() // Don't forget close the response body
				}
			}
			// }

		// case opened := <-try_create_client_once:
		// 	try_create_client_once = nil // disable this branch
		// 	if !opened {
		// 		// failed open it the first time,
		// 		// then we try to open file with time interval, until opened successfully.
		// 		logging.Debug("open the first time failed, try to open with interval of 1s")
		// 		try_create_client_tick = time.Tick(time.Second * 1)
		// 	} else {
		// 		startConsuming = b.updates
		// 	}
		// case <-try_create_client_tick:
		// 	// try to open with interval
		// 	err := b.ElasticsearchClient()
		// 	if b.output != nil && err == nil {
		// 		// finally opened.
		// 		try_create_client_tick = nil
		// 		startConsuming = b.updates
		// 	} else {
		// 		logging.Critical("elasticsearch backend trying to open file but failed: %s", err)
		// 	}
		case errc := <-b.closing:
			// fmt.Println("errc <- b.closing")
			logging.Debug("elasticsearch backend .loop closing")
			startConsuming = nil // stop comsuming
			errc <- nil
			close(b.updates)
			logging.Debug("elasticsearch backend .loop stopped")
			return
		}
	}
}

func (b *elasticsearchBackend) Updates() chan<- newcore.MultiDataPoint {
	return b.updates
}

func (b *elasticsearchBackend) Close() error {
	errc := make(chan error)
	b.closing <- errc
	return <-errc
}

func (b *elasticsearchBackend) Name() string {
	return b.name
}
