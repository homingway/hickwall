// write MultiDataPoint into a file.
// no rotation currently

package backends

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"github.com/oliveagle/hickwall/backends/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"github.com/oliveagle/hickwall/utils"
	"time"
)

var (
	_ = time.Now()
	_ = fmt.Sprintf("")
)

type influxdbBackend struct {
	name    string
	closing chan chan error             // for Close
	updates chan newcore.MultiDataPoint // for receive updates

	// influxdb backend specific attributes
	conf    *config.Transport_influxdb
	output  InfluxdbClient
	version string
}

func NewInfluxdbBackend(name string, conf *config.Transport_influxdb) (*influxdbBackend, error) {
	s := &influxdbBackend{
		name:    name,
		closing: make(chan chan error),
		updates: make(chan newcore.MultiDataPoint),
		conf:    conf,
		version: influxdbParseVersionFromString(conf.Version),
	}

	go s.loop()
	return s, nil
}

func (b *influxdbBackend) newInfluxdbClientFromConf() error {
	iclient, err := NewInfluxdbClient(map[string]interface{}{
		"Host":         b.conf.Host,
		"URL":          b.conf.URL,
		"Username":     b.conf.Username,
		"Password":     b.conf.Password,
		"UserAgent":    "",
		"Database":     b.conf.Database,
		"FlatTemplate": b.conf.FlatTemplate,
	}, b.version)
	if err != nil && iclient == nil {
		logging.Error("failed to create influxdb client: ", err)
		return fmt.Errorf("failed to create influxdb client: ", err)
	}
	b.output = iclient
	return nil
}

func (b *influxdbBackend) loop() {
	var (
		startConsuming         <-chan newcore.MultiDataPoint
		try_create_client_once chan bool
		try_create_client_tick <-chan time.Time
	)
	startConsuming = b.updates
	logging.Debug("influxdb backend loop started ")

	for {
		if b.output == nil && try_create_client_once == nil && try_create_client_tick == nil {
			startConsuming = nil // disable consuming
			try_create_client_once = make(chan bool)
			// try to create influxdb client the first time async.
			go func() {
				err := b.newInfluxdbClientFromConf()
				if err == nil {
					try_create_client_once <- true
				} else {
					try_create_client_once <- false
				}
			}()
		}

		//TODO: Flush_interval and Max_batch_size
		select {
		case md := <-startConsuming:
			if b.output != nil {
				points := []client.Point{}
				for _, p := range md {
					// logging.Debug(p.Metric.Clean())
					// logging.Debug(utils.Convert(p.Value))
					points = append(points, client.Point{
						Measurement: p.Metric.Clean(),
						Time:        p.Timestamp,
						Fields: map[string]interface{}{
							"value": utils.Convert(p.Value),
						},
						Tags: p.Tags, //TODO: Tags
					})
				}
				write := client.BatchPoints{
					Database:        b.conf.Database,
					RetentionPolicy: b.conf.RetentionPolicy,
					Points:          points,
				}
				// logging.Debugf("write: count: %d", len(md))

				//FIXME: connection timeout?
				resp, err := b.output.Write(write)
				if err != nil {
					logging.Errorf("failed to write into influxdb: %v, %+v", err, resp)
				}
			}
		case opened := <-try_create_client_once:
			try_create_client_once = nil // disable this branch
			if !opened {
				// failed open it the first time,
				// then we try to open file with time interval, until opened successfully.
				logging.Debug("open the first time failed, try to open with interval of 1s")
				try_create_client_tick = time.Tick(time.Second * 1)
			} else {
				startConsuming = b.updates
			}
		case <-try_create_client_tick:
			// try to open with interval
			err := b.newInfluxdbClientFromConf()
			if b.output != nil && err == nil {
				// finally opened.
				try_create_client_tick = nil
				startConsuming = b.updates
			} else {
				logging.Critical("influxdb backend trying to open file but failed: %s", err)
			}
		case errc := <-b.closing:
			// fmt.Println("errc <- b.closing")
			logging.Debug("influxdb backend .loop closing")
			startConsuming = nil // stop comsuming
			errc <- nil
			close(b.updates)
			logging.Debug("influxdb backend .loop stopped")
			return
		}
	}
}

func (b *influxdbBackend) Updates() chan<- newcore.MultiDataPoint {
	return b.updates
}

func (b *influxdbBackend) Close() error {
	errc := make(chan error)
	b.closing <- errc
	return <-errc
}

func (b *influxdbBackend) Name() string {
	return b.name
}
