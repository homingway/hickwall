package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectors/config"
	"github.com/oliveagle/hickwall/newcore"
	"github.com/oliveagle/hickwall/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	_ = fmt.Sprintf("")
)

type cmd_collector struct {
	name     string // collector name
	interval time.Duration
	enabled  bool
	prefix   string

	// cmd_collector specific attributes
	config   config.Config_command
	timeout  time.Duration
	cmd_name string
	cmd_args []string
	tags     newcore.TagSet
}

// newCollector returns a Collector for uri.
func NewCmdCollector(name, prefix string, conf config.Config_command) newcore.Collector {

	var (
		cmd_name string
		cmd_args []string
	)

	if len(conf.Cmd) > 0 {
		cmd_name = conf.Cmd[0]
	}
	if len(conf.Cmd) > 1 {
		cmd_args = conf.Cmd[1:]
	}

	f := &cmd_collector{
		name:     name,
		enabled:  true,
		prefix:   prefix,
		interval: conf.Interval.MustDuration(time.Second),
		config:   conf,
		timeout:  conf.Interval.MustDuration(time.Second * 10),
		cmd_name: cmd_name,
		cmd_args: cmd_args,
		tags:     conf.Tags.Copy(),
	}
	return f
}

func (f *cmd_collector) Name() string {
	return f.name
}

func (f *cmd_collector) Close() error {
	return nil
}

func (f *cmd_collector) ClassName() string {
	return "cmd_collector"
}

func (f *cmd_collector) IsEnabled() bool {
	return f.enabled
}

func (f *cmd_collector) Interval() time.Duration {
	return f.interval
}

func (f *cmd_collector) CollectOnce() *newcore.CollectResult {
	var items newcore.MultiDataPoint

	err := utils.ReadCommandTimeout(f.timeout, func(line string) error {
		// fmt.Println("read command timeout: ", line)
		slices := strings.Split(line, "|")
		if len(slices) == 3 {
			// first supported format.  metric|timestamp|value
			metric := slices[0]

			// timestamp
			sec, err := strconv.ParseInt(slices[1], 10, 64)
			if err != nil {
				log.Println("ERROR: cannot parse epoch timestamp: %v", slices[1])
				return nil
			}
			timestamp := time.Unix(sec, 0)

			// value
			value, err := strconv.ParseFloat(slices[2], 64)
			if err != nil {
				log.Println("ERROR: cannot parse float result: %s  err: %v", slices[2], err)
				return nil
			}

			AddTS(&items, f.prefix, metric, timestamp, value, f.tags, "", "", "")

			//TODO: add DataType Support for Cmd collector.
			// } else if len(slices) == 5 {
			// 	// longer format:  type|unit|metric|timestamp|value
			// 	ratetype, err := metadata.ParseRateType(slices[0])
			// 	if err != nil {
			// 		log.Errorf("cannot parse epoch ratetype: %v", slices[0])
			// 		// should continue
			// 	}

			// 	unit := slices[1]

			// 	metric := collectorlib.NormalizeMetricKey(slices[2])

			// 	// timestamp
			// 	sec, err := strconv.ParseInt(slices[3], 10, 64)
			// 	if err != nil {
			// 		log.Errorf("cannot parse epoch timestamp: %v", slices[3])
			// 		return nil
			// 	}
			// 	timestamp := time.Unix(sec, 0)

			// 	// value
			// 	value, err := strconv.ParseFloat(slices[4], 64)
			// 	if err != nil {
			// 		log.Errorf("cannot parse float result: %s  err: %v", slices[4], err)
			// 		return nil
			// 	}

			// 	newcore.AddTS(&md, metric, timestamp, value, st.Tags, ratetype, unit, "")

		}
		return nil
	}, f.cmd_name, f.cmd_args...)

	if err != nil {
		return &newcore.CollectResult{
			Collected: &items,
			Next:      time.Now().Add(f.interval),
			Err:       err,
		}
	} else {
		return &newcore.CollectResult{
			Collected: &items,
			Next:      time.Now().Add(f.interval),
			Err:       nil,
		}
	}

}
