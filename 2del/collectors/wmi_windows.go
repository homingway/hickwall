// +build windows

package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/utils"

	"github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"

	log "github.com/oliveagle/seelog"
	"regexp"
	"strings"
	"time"
)

func init() {
	defer utils.Recover_and_log()

	collector_factories["win_wmi"] = factory_win_wmi

	// for collector := range builtin_win_wmi() {
	// 	builtin_collectors = append(builtin_collectors, collector)
	// }
}

//TODO: we don't allow  multpile leveled template {{.A.B}}
var (
	win_wmi_pat_format, _ = regexp.Compile("\\/format:\\w+(.xsl)?")
	win_wmi_pat_get, _    = regexp.Compile("\\bget\\b")
	win_wmi_pat_field, _  = regexp.Compile(`\{\{\.\w+((_?)+\w+)+\}\}`)
)

func WmiQueryWithFields(query string, fields []string) ([]map[string]string, error) {
	defer utils.Recover_and_log()

	resultMap := []map[string]string{}

	// init COM, oh yeah
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		log.Error("oleutil.CreateObject Failed: ", err)
		return resultMap, err
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Error("QueryInterface Failed: ", err)
		return resultMap, err
	}
	defer wmi.Release()

	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer")
	if err != nil {
		log.Error("Connect to Server Failed", err)
		return resultMap, err
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	resultRaw, err := oleutil.CallMethod(service, "ExecQuery", query)
	if err != nil {
		log.Error("ExecQuery Failed: ", err)
		return resultMap, err
	}

	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, err := oleutil.GetProperty(result, "Count")
	if err != nil {
		log.Error("Get result count Failed: ", err)
		return resultMap, err
	}
	count := int(countVar.Val)

	for i := 0; i < count; i++ {
		itemMap := make(map[string]string)

		itemRaw, err := oleutil.CallMethod(result, "ItemIndex", i)
		if err != nil {
			log.Error("ItemIndex failed: ", err)
			return resultMap, err
		}

		item := itemRaw.ToIDispatch()
		defer item.Release()

		for _, field := range fields {
			asString, err := oleutil.GetProperty(item, field)

			if err == nil {
				itemMap[field] = fmt.Sprintf("%v", asString.Value())
			} else {
				fmt.Println(err)
			}
		}

		resultMap = append(resultMap, itemMap)
	}

	return resultMap, nil
}

func builtin_win_wmi() <-chan Collector {
	defer utils.Recover_and_log()

	// large_interval_queries -----------------------------------------------------------

	large_interval_queries := []config.Conf_win_wmi_query{}

	large_interval_queries = append(large_interval_queries, config.Conf_win_wmi_query{
		Query: "select Name, NumberOfCores, NumberOfLogicalProcessors from Win32_Processor",
		Metrics: []config.Conf_win_wmi_query_metric{
			config.Conf_win_wmi_query_metric{
				Value_from: "Name",
				Metric:     "win.wmi.cpu.name",
			},
			config.Conf_win_wmi_query_metric{
				Value_from: "NumberOfCores",
				Metric:     "win.wmi.cpu.numberofcores",
			},
			config.Conf_win_wmi_query_metric{
				Value_from: "NumberOfLogicalProcessors",
				Metric:     "win.wmi.cpu.numberoflogicalprocessors",
			},
		}})

	large_interval_queries = append(large_interval_queries, config.Conf_win_wmi_query{
		Query: "select * from Win32_ComputerSystem",
		Metrics: []config.Conf_win_wmi_query_metric{
			config.Conf_win_wmi_query_metric{
				Value_from: "TotalPhysicalMemory",
				Metric:     "win.wmi.mem.totalphysicalmemory",
			},
			config.Conf_win_wmi_query_metric{
				Value_from: "Domain",
				Metric:     "win.wmi.net.domain",
			},
		}})

	large_interval_queries = append(large_interval_queries, config.Conf_win_wmi_query{
		Query: "select Name, FileSystem, FreeSpace, Size from Win32_LogicalDisk where MediaType=11 or mediatype=12",
		Metrics: []config.Conf_win_wmi_query_metric{
			config.Conf_win_wmi_query_metric{
				Value_from: "Size",
				Metric:     "win.wmi.fs.size.bytes",
				Tags: map[string]string{
					"mount": "{{.Name}}",
					// "fs_type": "{{.FileSystem}}",
				},
			},
		}})

	large_interval_queries = append(large_interval_queries, config.Conf_win_wmi_query{
		Query: "select * from Win32_OperatingSystem",
		Metrics: []config.Conf_win_wmi_query_metric{
			config.Conf_win_wmi_query_metric{
				Value_from: "Caption",
				Metric:     "win.wmi.os.caption",
			},
			config.Conf_win_wmi_query_metric{
				Value_from: "CSDVersion",
				Metric:     "win.wmi.os.csdversion",
			},
		}})

	// TODO: iisInstalled, when W3svc is not installed, should give a value.
	large_interval_queries = append(large_interval_queries, config.Conf_win_wmi_query{
		Query: "select * from Win32_Service where Name='W3svc'",
		Metrics: []config.Conf_win_wmi_query_metric{
			config.Conf_win_wmi_query_metric{
				Value_from: "State",
				Metric:     "win.wmi.service.iis.state",
				Default:    "IIS Not Installed",
			},
		}})

	// TODO: rsaInstalled

	// small_interval_queries -----------------------------------------------------------

	small_interval_queries := []config.Conf_win_wmi_query{}
	small_interval_queries = append(small_interval_queries, config.Conf_win_wmi_query{
		Query: "select Name, FileSystem, FreeSpace, Size from Win32_LogicalDisk where MediaType=11 or mediatype=12",
		Metrics: []config.Conf_win_wmi_query_metric{
			config.Conf_win_wmi_query_metric{
				Value_from: "FreeSpace",
				Metric:     "win.wmi.fs.freespace.bytes",
				Tags: map[string]string{
					"mount": "{{.Name}}",
					// "fs_type": "{{.FileSystem}}",
				},
			},
		}})

	// config_list ------------------------------------------------------------------------

	config_list := []config.Conf_win_wmi{
		config.Conf_win_wmi{
			Interval: "60m",
			Queries:  large_interval_queries,
		},
		config.Conf_win_wmi{
			Interval: "1s",
			Queries:  small_interval_queries,
		},
	}

	return factory_win_wmi("builtin_win_wmi", config_list)
	// return factory_win_wmi(config_list)
}

// conf interface{}: []config.Conf_win_wmi
func factory_win_wmi(name string, conf interface{}) <-chan Collector {
	defer utils.Recover_and_log()

	var out = make(chan Collector)
	go func() {
		var (
			states      state_win_wmi
			config_list []config.Conf_win_wmi
			// cf               config.Conf_win_wmi
			default_interval = time.Duration(60) * time.Minute
			runtime_conf     = config.GetRuntimeConf()
		)

		if conf != nil {
			config_list = conf.([]config.Conf_win_wmi)
			for idx, cf := range config_list {

				// cf = conf.(config.Conf_win_wmi)

				interval, err := collectorlib.ParseInterval(cf.Interval)
				if err != nil {
					log.Errorf("cannot parse interval of collector_wmi: %s - %v", cf.Interval, err)
					interval = default_interval
				}
				states.Interval = interval

				states.queries = []config.Conf_win_wmi_query{}

				for _, query_obj := range cf.Queries {
					//TODO: validate query

					// merge tags
					query_obj.Tags = AddTags.Copy().Merge(runtime_conf.Client.Tags).Merge(cf.Tags).Merge(query_obj.Tags)

					states.queries = append(states.queries, query_obj)
				}

				out <- &IntervalCollector{
					F:            c_win_wmi,
					EnableFunc:   nil,
					name:         fmt.Sprintf("win_wmi_%s_%d", name, idx),
					states:       states,
					Interval:     states.Interval,
					factory_name: "win_wmi",
				}
			}

		}

		close(out)
	}()
	return out
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

func get_fields_of_query(query config.Conf_win_wmi_query) []string {
	fields := map[string]bool{}
	for _, item := range query.Metrics {
		if len(item.Value_from) > 0 {
			fields[item.Value_from] = true
		}

		for _, f := range win_wmi_pat_field.FindAllString(item.Metric, -1) {
			key := f[3 : len(f)-2]
			if len(key) > 0 {
				fields[key] = true
			}

		}

		for _, value := range item.Tags {
			// fmt.Println("item.Tags.value: ", value)
			for _, f := range win_wmi_pat_field.FindAllString(value, -1) {
				key := f[3 : len(f)-2]
				if len(key) > 0 {
					fields[key] = true
				}
			}
		}
	}

	results := []string{}
	for key, _ := range fields {
		results = append(results, key)
	}

	// fmt.Println("results: ", results)

	return results
}

func c_win_wmi(states interface{}) (collectorlib.MultiDataPoint, error) {
	defer utils.Recover_and_log()

	var md collectorlib.MultiDataPoint
	var st state_win_wmi

	if states != nil {
		st = states.(state_win_wmi)
	}

	for _, query := range st.queries {

		fields := get_fields_of_query(query)

		results, err := WmiQueryWithFields(query.Query, fields)
		if err != nil {
			continue
		}
		if len(results) > 0 {
			for _, record := range results {
				for _, item := range query.Metrics {

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

					if value, ok := record[item.Value_from]; ok == true {

						Add(&md, metric, value, tags, "", "", "")

					} else if item.Default != "" {

						Add(&md, metric, item.Default, tags, "", "", "")

					}
				}
			}
		} else {
			for _, item := range query.Metrics {
				if item.Default != "" {
					// no templating support if no data got
					if strings.Contains(item.Metric, "{{") {
						continue
					}
					for _, value := range item.Tags {
						if strings.Contains(value, "{{") {
							continue
						}
					}

					tags := AddTags.Copy().Merge(query.Tags).Merge(item.Tags)

					Add(&md, item.Metric, item.Default, tags, "", "", "")
				}
			}
		}

	}

	return md, nil
}