

// func builtin_win_pdh() <-chan Collector {
//  defer utils.Recover_and_log()
//  config_list := []config.Conf_win_pdh{}

//  queries := []config.Conf_win_pdh_query{}

//  queries = append(queries, config.Conf_win_pdh_query{
//      Query:  "\\System\\Processes",
//      Metric: "win.processes.count"})
//  queries = append(queries, config.Conf_win_pdh_query{
//      Query:  "\\Memory\\Available Bytes",
//      Metric: "win.memory.available.bytes"})
//  queries = append(queries, config.Conf_win_pdh_query{
//      Query:  "\\Process(_Total)\\Working Set",
//      Metric: "win.process.working_set.total.bytes"})
//  queries = append(queries, config.Conf_win_pdh_query{
//      Query:  "\\Memory\\Cache Bytes",
//      Metric: "win.memory.cache.bytes"})

//  config_list = append(config_list, config.Conf_win_pdh{
//      Interval: "2s",
//      Queries:  queries,
//  })

//  //TODO: hickwall internal metrics should not dependend on any third party dll/so.
//  client_perf_queries := []config.Conf_win_pdh_query{}

//  // private memory size
//  client_perf_queries = append(client_perf_queries, config.Conf_win_pdh_query{
//      Query:  "\\Process(hickwall)\\Working Set - Private",
//      Metric: "hickwall.client.mem.private_working_set.bytes"})
//  client_perf_queries = append(client_perf_queries, config.Conf_win_pdh_query{
//      Query:  "\\Process(hickwall)\\Working Set",
//      Metric: "hickwall.client.mem.working_set.bytes"})

//  config_list = append(config_list, config.Conf_win_pdh{
//      Interval: "2s",
//      Queries:  client_perf_queries,
//  })

//  return factory_win_pdh("builtin_win_pdh", config_list)
// }




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