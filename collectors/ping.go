package collectors

import (
	"fmt"
	"github.com/GaryBoone/GoStats/stats"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/utils"
	log "github.com/oliveagle/seelog"
	"github.com/tatsushid/go-fastping"
	"math"
	"net"
	"strings"
	"time"
)

func init() {
	defer utils.Recover_and_log()

	collector_factories["ping"] = factory_ping
}

// func factory_ping(conf interface{}) <-chan Collector {
func factory_ping(name string, conf interface{}) <-chan Collector {
	defer utils.Recover_and_log()

	log.Debug("factory_ping")

	var out = make(chan Collector)
	go func() {
		var (
			config_list []config.Conf_ping
			// cf               config.Conf_ping
			default_interval = time.Duration(1) * time.Second
			default_timeout  = time.Duration(5) * time.Second
		)

		if conf != nil {
			config_list = conf.([]config.Conf_ping)

			for collector_idx, cf := range config_list {
				// cf = conf.(config.Conf_ping)

				interval, err := collectorlib.ParseInterval(cf.Interval)
				if err != nil {
					log.Errorf("cannot parse interval of collector_ping: %s - %v", cf.Interval, err)
					interval = default_interval
				}
				timeout, err := collectorlib.ParseInterval(cf.Timeout)
				if err != nil {
					log.Errorf("cannot parse timeout of collector_ping: %s - %v", cf.Timeout, err)
					timeout = default_timeout
				}

				for target_idx, target := range cf.Targets {
					var (
						states state_c_ping
					)

					states.Interval = interval
					states.Target = target
					states.Conf = cf
					states.Timeout = timeout

					out <- &IntervalCollector{
						F:            C_ping,
						Enable:       nil,
						name:         fmt.Sprintf("ping_%s_%d_%d", name, collector_idx, target_idx),
						states:       states,
						Interval:     states.Interval,
						factory_name: "ping",
					}
				}
			}
		}
		close(out)
	}()
	return out
}

type state_c_ping struct {
	Interval time.Duration
	Target   string
	Conf     config.Conf_ping
	Timeout  time.Duration
}

// hickwall process metrics, only runtime stats
func C_ping(states interface{}) (collectorlib.MultiDataPoint, error) {
	defer utils.Recover_and_log()

	var (
		md           collectorlib.MultiDataPoint
		runtime_conf = config.GetRuntimeConf()
		p            = fastping.NewPinger()
		d            stats.Stats

		state = states.(state_c_ping)

		rtt_chan = make(chan float64)
	)

	tags := AddTags.Copy().Merge(runtime_conf.Client.Tags)
	tags["target"] = state.Target

	if state.Conf.Packets <= 0 {
		return md, fmt.Errorf("collector_ping: packets should be greater than zero")
	}

	ip, err := net.ResolveIPAddr("ip4:icmp", state.Target)
	if err != nil {
		log.Errorf("collector_ping: DNS resolve error: %v", err)
		return md, fmt.Errorf("collector_ping: DNS resolve error: %v", err)
	}

	p.MaxRTT = state.Timeout
	p.AddIPAddr(ip)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		rtt_chan <- float64(rtt.Nanoseconds() / 1000 / 1000)
	}

	go func() {
		for i := 0; i < state.Conf.Packets; i++ {
			err = p.Run()
			if err != nil {
				fmt.Println("run err", err)
			}
		}
		close(rtt_chan)
	}()

	for rtt := range rtt_chan {
		d.Update(rtt)
	}

	for _, sl := range state.Conf.Collect {
		switch strings.ToLower(sl) {
		case "time_min":
			Add(&md, fmt.Sprintf("%s.%s", state.Conf.Metric_key, "time_min"), d.Min(), tags, "", "", "")
		case "time_max":
			Add(&md, fmt.Sprintf("%s.%s", state.Conf.Metric_key, "time_max"), d.Max(), tags, "", "", "")
		case "time_avg":
			Add(&md, fmt.Sprintf("%s.%s", state.Conf.Metric_key, "time_avg"), d.Mean(), tags, "", "", "")
		case "time_mdev":
			std := d.SampleStandardDeviation()
			if math.IsNaN(std) {
				std = 0
			}
			Add(&md, fmt.Sprintf("%s.%s", state.Conf.Metric_key, "time_mdev"), std, tags, "", "", "")
		case "ip":
			Add(&md, fmt.Sprintf("%s.%s", state.Conf.Metric_key, "ip"), ip.IP.String(), tags, "", "", "")
		case "lost_pct":
			lost_pct := float64((state.Conf.Packets-d.Count())/state.Conf.Packets) * 100
			Add(&md, fmt.Sprintf("%s.%s", state.Conf.Metric_key, "lost_pct"), lost_pct, tags, "", "", "")
		}
	}
	return md, nil
}
