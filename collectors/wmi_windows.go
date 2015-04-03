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

	collector_factories["win_wmi"] = factory_win_wmi

	builtin_collectors = append(builtin_collectors, builtin_win_wmi())

	// log.Debug("Initialized builtin collector: win_wmi")
}

//TODO: we don't allow  multpile leveled template {{.A.B}}
var (
	win_wmi_pat_format, _ = regexp.Compile("\\/format:\\w+(.xsl)?")
	win_wmi_pat_get, _    = regexp.Compile("\\bget\\b")
	win_wmi_pat_field, _  = regexp.Compile(`\{\{\.\w+((_?)+\w+)+\}\}`)
)

func WmiQueryWithFields(query string, fields []string) []map[string]string {
	// init COM, oh yeah
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, _ := oleutil.CreateObject("WbemScripting.SWbemLocator")
	defer unknown.Release()

	wmi, _ := unknown.QueryInterface(ole.IID_IDispatch)
	defer wmi.Release()

	serviceRaw, _ := oleutil.CallMethod(wmi, "ConnectServer")
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	resultRaw, _ := oleutil.CallMethod(service, "ExecQuery", query)

	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, _ := oleutil.GetProperty(result, "Count")
	count := int(countVar.Val)

	resultMap := []map[string]string{}

	for i := 0; i < count; i++ {
		itemMap := make(map[string]string)

		itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i)

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

	return resultMap
}

func builtin_win_wmi() Collector {

	queries := []config.Conf_win_wmi_query{}

	queries = append(queries, config.Conf_win_wmi_query{
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

	queries = append(queries, config.Conf_win_wmi_query{
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

	queries = append(queries, config.Conf_win_wmi_query{
		Query: "select Name, FileSystem, FreeSpace, Size from Win32_LogicalDisk where MediaType=11 or mediatype=12",
		Metrics: []config.Conf_win_wmi_query_metric{
			config.Conf_win_wmi_query_metric{
				Value_from: "Size",
				Metric:     "win.wmi.fs.size.{{.Name}}.bytes",
				Tags: map[string]string{
					"mount":   "{{.Name}}",
					"fs_type": "{{.FileSystem}}",
				},
			},
		}})

	queries = append(queries, config.Conf_win_wmi_query{
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
	queries = append(queries, config.Conf_win_wmi_query{
		Query: "select * from Win32_Service where Name='W3svc'",
		Metrics: []config.Conf_win_wmi_query_metric{
			config.Conf_win_wmi_query_metric{
				Value_from: "State",
				Metric:     "win.wmi.service.iis.state",
				Default:    "IIS Not Installed",
			},
		}})

	// TODO: rsaInstalled

	conf := config.Conf_win_wmi{
		Interval: "60m",
		Queries:  queries,
	}

	return factory_win_wmi("builtin_win_wmi", conf)
}

func factory_win_wmi(name string, conf interface{}) Collector {
	var states state_win_wmi
	var cf config.Conf_win_wmi
	var default_interval = time.Duration(60) * time.Minute

	if conf != nil {
		cf = conf.(config.Conf_win_wmi)

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
			query_obj.Tags = AddTags.Copy().Merge(config.Conf.Tags).Merge(cf.Tags).Merge(query_obj.Tags)

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

		for _, value := range query.Tags {
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

	return results
}

func c_win_wmi(states interface{}) (collectorlib.MultiDataPoint, error) {
	var md collectorlib.MultiDataPoint
	var st state_win_wmi

	if states != nil {
		st = states.(state_win_wmi)
	}

	for _, query := range st.queries {

		fields := get_fields_of_query(query)

		results := WmiQueryWithFields(query.Query, fields)
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
