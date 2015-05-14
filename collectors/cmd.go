package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/collectorlib/metadata"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/utils"
	log "github.com/oliveagle/seelog"
	// "regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	defer utils.Recover_and_log()
	collector_factories["cmd"] = factory_cmd
}

func factory_cmd(name string, conf interface{}) <-chan Collector {
	// func factory_cmd(conf interface{}) <-chan Collector {
	defer utils.Recover_and_log()

	var out = make(chan Collector)
	go func() {
		var (
			config_list []config.Conf_cmd
			// cf               config.Conf_cmd
			state            state_cmd
			default_interval = time.Duration(1) * time.Minute
			runtime_conf     = config.GetRuntimeConf()
		)

		if conf != nil {
			config_list = conf.([]config.Conf_cmd)

			for idx, cf := range config_list {
				// cf = conf.(config.Conf_cmd)

				interval, err := collectorlib.ParseInterval(cf.Interval)
				if err != nil {
					log.Errorf("cannot parse interval of collector_cmd: %s - %v", cf.Interval, err)
					interval = default_interval
				}
				state.Interval = interval
				state.Cmd = cf.Cmd
				state.Tags = AddTags.Copy().Merge(runtime_conf.Client.Tags).Merge(cf.Tags)

				out <- &IntervalCollector{
					F:            c_cmd,
					EnableFunc:   nil,
					name:         fmt.Sprintf("cmd_%s_%d", name, idx),
					states:       state,
					Interval:     state.Interval,
					factory_name: "cmd",
				}
			}
		}
		close(out)
	}()
	return out
}

type state_cmd struct {
	Interval time.Duration
	Cmd      []string
	Tags     map[string]string
}

func c_cmd(states interface{}) (collectorlib.MultiDataPoint, error) {
	defer utils.Recover_and_log()

	log.Debugf("collector:c_cmd start")
	var md collectorlib.MultiDataPoint
	st := states.(state_cmd)

	name := st.Cmd[0]
	args := []string{}
	if len(st.Cmd) > 1 {
		args = st.Cmd[1:]
	}

	fmt.Println("")
	err := collectorlib.ReadCommand(func(s string) error {
		// log.Debugf("c_cmd: line: %s", s)
		slices := strings.Split(s, "|")
		if len(slices) == 3 {
			// first supported format.  metric|timestamp|value
			metric := collectorlib.NormalizeMetricKey(slices[0])

			// timestamp
			sec, err := strconv.ParseInt(slices[1], 10, 64)
			if err != nil {
				log.Errorf("cannot parse epoch timestamp: %v", slices[1])
				return nil
			}
			timestamp := time.Unix(sec, 0)

			// value
			value, err := strconv.ParseFloat(slices[2], 64)
			if err != nil {
				log.Errorf("cannot parse float result: %s  err: %v", slices[2], err)
				return nil
			}

			AddTS(&md, metric, timestamp, value, st.Tags, "", "", "")

		} else if len(slices) == 5 {
			// longer format:  type|unit|metric|timestamp|value
			ratetype, err := metadata.ParseRateType(slices[0])
			if err != nil {
				log.Errorf("cannot parse epoch ratetype: %v", slices[0])
				// should continue
			}

			unit := slices[1]

			metric := collectorlib.NormalizeMetricKey(slices[2])

			// timestamp
			sec, err := strconv.ParseInt(slices[3], 10, 64)
			if err != nil {
				log.Errorf("cannot parse epoch timestamp: %v", slices[3])
				return nil
			}
			timestamp := time.Unix(sec, 0)

			// value
			value, err := strconv.ParseFloat(slices[4], 64)
			if err != nil {
				log.Errorf("cannot parse float result: %s  err: %v", slices[4], err)
				return nil
			}

			AddTS(&md, metric, timestamp, value, st.Tags, ratetype, unit, "")
		} else {
			log.Errorf("unsupported output format in stdout: %s", s)
		}
		return nil
	}, name, args...)

	if err != nil {
		log.Errorf("collector:c_cmd error: %v", err)
	}

	log.Debugf("collector:c_cmd return")
	return md, nil
}
