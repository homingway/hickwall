package collectors

import (
	"fmt"
	// "github.com/kr/pretty"
	"github.com/oliveagle/go-collectors/datapoint"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/utils"
	// log "github.com/cihub/seelog"
	"github.com/oliveagle/go-collectors/util"
	"regexp"
	"strings"
	"time"
)

func init() {
	collector_factories["win_wmi"] = factory_win_wmi

	builtin_collectors = append(builtin_collectors, builtin_win_wmi())
}

var (
	win_wmi_pat_format, _ = regexp.Compile("\\/format:\\w+(.xsl)?")
	win_wmi_pat_get, _    = regexp.Compile("\\bget\\b")
)

func wmi_windows_query_cmdline(query string) []map[string]string {

	// fmt.Println("query: ", query)
	the_query := strings.ToLower(query)

	if win_wmi_pat_format.MatchString(the_query) == true {
		the_query = win_wmi_pat_format.ReplaceAllString(the_query, "/format:textvaluelist")
	} else {
		if win_wmi_pat_get.MatchString(the_query) == true {
			the_query = strings.Join([]string{the_query, " /format:textvaluelist"}, "")
		} else {
			the_query = strings.Join([]string{the_query, " get /format:textvaluelist"}, "")
		}
	}
	// fmt.Println(the_query)

	results := []map[string]string{}

	parts := []string{}
	name := ""
	for idx, part := range strings.Split(strings.Trim(the_query, " "), " ") {
		if part != "" {
			if idx == 0 {
				name = part
			} else {
				parts = append(parts, part)
			}
		}
	}
	if name != "" {
		// for _, p := range parts {
		//  fmt.Println(p)
		// }

		line_num := 0
		new_record := false
		record := map[string]string{}

		lines := []string{}
		util.ReadCommand(func(line string) error {
			// lines = append(lines, line)
			if len(lines) < 3 {
				lines = append(lines, line)

				// fmt.Printf("%3d %2d %5v %45s %45s\n", line_num, len(line), new_record, "", line)
				line_num += 1
				return nil
			} else if len(lines) == 3 {
				if len(lines[0]) == len(lines[1]) && len(lines[0]) == 1 {
					new_record = true
					if len(record) > 0 {
						results = append(results, record)

						// record = nil
						record = map[string]string{}
					}
				} else {
					new_record = false
				}
				// fmt.Println(line_num, len(line), new_record, line)

				property := strings.Trim(lines[2], "\r\n")

				// fmt.Println(record)
				// fmt.Println(line_num, len(line), new_record, record, line)
				// fmt.Printf("%3d %2d %5v %45s %45s\n", line_num, len(line), new_record, property, line)

				property_array := strings.Split(property, "=")

				if len(property_array) == 2 && property_array[0] != "" {
					if strings.HasPrefix(property_array[1], `{`) && strings.HasSuffix(property_array[1], `}`) {
						// remove `{` and  `}` from value string
						record[property_array[0]] = property_array[1][1 : len(property_array[1])-1]
					} else {
						record[property_array[0]] = property_array[1]
					}
				} else if len(property_array) == 1 && property_array[0] != "" {
					record[property_array[0]] = ""
				}

				lines = append(lines, line)
				lines = lines[1:]
			}

			line_num += 1
			return nil
		}, name, parts...)

		if len(record) > 0 {
			results = append(results, record)
		}

	}
	return results
}

func builtin_win_wmi() Collector {

	// queries := []config.Conf_win_wmi_query{}

	// queries = append(queries, config.Conf_win_wmi_query{
	// 	Query:  "\\System\\Processes",
	// 	Metric: "win.processes.count"})
	// queries = append(queries, config.Conf_win_wmi_query{
	// 	Query:  "\\Memory\\Available Bytes",
	// 	Metric: "win.memory.available_bytes"})
	// queries = append(queries, config.Conf_win_wmi_query{
	// 	Query:  "\\Processes(_Total)\\Working Set",
	// 	Metric: "win.processes.working_set.total"})
	// queries = append(queries, config.Conf_win_wmi_query{
	// 	Query:  "\\Memory\\Cache Bytes",
	// 	Metric: "win.memory.cache_bytes"})

	// conf := config.Conf_win_wmi{
	// 	Interval: 2,
	// 	Queries:  queries,
	// }
	return nil
	// return factory_win_pdh("builtin_win_pdh", conf)
}

func factory_win_wmi(name string, conf interface{}) Collector {
	var states state_win_wmi
	var cf config.Conf_win_wmi

	if conf != nil {
		cf = conf.(config.Conf_win_wmi)
		// fmt.Println("factory_win_pdh: ", cf)
		// pretty.Println("factory_win_pdh:", cf)

		// states.map_metrics = make(map[string]string)
		states.Interval = time.Duration(cf.Interval) * time.Second
		states.queries = []config.Conf_win_wmi_query{}

		for _, query_obj := range cf.Queries {
			//TODO: validate query

			// merge tags
			query_obj.Tags = AddTags.Copy().Merge(config.Conf.Tags).Merge(cf.Tags).Merge(query_obj.Tags)
			// fmt.Println(query_obj.Tags)

			states.queries = append(states.queries, query_obj)
		}
	}

	return &IntervalCollector{
		F:        c_win_wmi,
		Enable:   nil,
		name:     name,
		states:   states,
		Interval: states.Interval,
	}
}

type state_win_wmi struct {
	Interval time.Duration

	// internal use only
	queries []config.Conf_win_wmi_query
}

func c_win_wmi_parse_metric_key(metric string, data map[string]string) (string, error) {
	if strings.Contains(metric, "{{") {
		return utils.ExecuteTemplate(metric, data, collectorlib.NormalizeMetricKey)
	} else {
		return metric, nil
	}

}

func c_win_wmi_parse_tags(tags map[string]string, data map[string]string) (map[string]string, error) {
	res := map[string]string{}

	for key, tag := range tags {
		if strings.Contains(tag, "{{") {
			tag_value, err := utils.ExecuteTemplate(tag, data, collectorlib.NormalizeTag)
			if err != nil {
				return res, err
			}
			res[key] = tag_value
		} else {
			res[key] = tag
		}
	}
	return res, nil
}

func c_win_wmi(states interface{}) (datapoint.MultiDataPoint, error) {
	var md datapoint.MultiDataPoint
	var st state_win_wmi

	if states != nil {
		st = states.(state_win_wmi)
		// fmt.Println("c_win_pdh states: ", states)
	}

	// fmt.Println(st)
	for _, query := range st.queries {
		// fmt.Println("-----------------------------------------------")
		// fmt.Println(query)
		// fmt.Println(query.Query)
		for _, record := range wmi_windows_query_cmdline(query.Query) {
			for _, item := range query.Metrics {
				// fmt.Println(item.Value_from)
				if value, ok := record[item.Value_from]; ok == true {
					// fmt.Println(item.Value_from, value, item.Metric)

					metric, err := c_win_wmi_parse_metric_key(item.Metric, record)
					if err != nil {
						fmt.Println(err)
						continue
					}

					tags, err := c_win_wmi_parse_tags(item.Tags, record)
					if err != nil {
						fmt.Println(err)
						continue
					}

					tags = AddTags.Copy().Merge(query.Tags).Merge(tags)

					Add(&md, metric, value, tags, "", "", "")
				}
			}
		}
	}

	// fmt.Println("")
	return md, nil
}
