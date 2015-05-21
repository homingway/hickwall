// +build windows
package collectors

// func builtin_win_pdh() <-chan Collector {
// 	defer utils.Recover_and_log()
// 	config_list := []config.Conf_win_pdh{}

// 	queries := []config.Conf_win_pdh_query{}

// 	queries = append(queries, config.Conf_win_pdh_query{
// 		Query:  "\\System\\Processes",
// 		Metric: "win.processes.count"})
// 	queries = append(queries, config.Conf_win_pdh_query{
// 		Query:  "\\Memory\\Available Bytes",
// 		Metric: "win.memory.available.bytes"})
// 	queries = append(queries, config.Conf_win_pdh_query{
// 		Query:  "\\Process(_Total)\\Working Set",
// 		Metric: "win.process.working_set.total.bytes"})
// 	queries = append(queries, config.Conf_win_pdh_query{
// 		Query:  "\\Memory\\Cache Bytes",
// 		Metric: "win.memory.cache.bytes"})

// 	config_list = append(config_list, config.Conf_win_pdh{
// 		Interval: "2s",
// 		Queries:  queries,
// 	})

// 	//TODO: hickwall internal metrics should not dependend on any third party dll/so.
// 	client_perf_queries := []config.Conf_win_pdh_query{}

// 	// private memory size
// 	client_perf_queries = append(client_perf_queries, config.Conf_win_pdh_query{
// 		Query:  "\\Process(hickwall)\\Working Set - Private",
// 		Metric: "hickwall.client.mem.private_working_set.bytes"})
// 	client_perf_queries = append(client_perf_queries, config.Conf_win_pdh_query{
// 		Query:  "\\Process(hickwall)\\Working Set",
// 		Metric: "hickwall.client.mem.working_set.bytes"})

// 	config_list = append(config_list, config.Conf_win_pdh{
// 		Interval: "2s",
// 		Queries:  client_perf_queries,
// 	})

// 	return factory_win_pdh("builtin_win_pdh", config_list)
// }
