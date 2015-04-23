// +build windows
// have to put pdh.go codez into this file to do cross compile from linux.
// don't know way

package collectors

import (
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	"time"

	log "github.com/oliveagle/seelog"

	"syscall"
	"unsafe"
)

/*
typedef long LONG;
typedef unsigned long DWORD;
typedef struct _PDH_FMT_COUNTERVALUE_DOUBLE
{
    DWORD CStatus;
    double DoubleValue;
}PDH_FMT_COUNTERVALUE_DOUBLE;
*/
import "C"

// PDH error codes, which can be returned by all Pdh* functions. Taken from mingw-w64 pdhmsg.h
const (
	PDH_CSTATUS_VALID_DATA                     = 0x00000000 // The returned data is valid.
	PDH_CSTATUS_NEW_DATA                       = 0x00000001 // The return data value is valid and different from the last sample.
	PDH_CSTATUS_NO_MACHINE                     = 0x800007D0 // Unable to connect to the specified computer, or the computer is offline.
	PDH_CSTATUS_NO_INSTANCE                    = 0x800007D1
	PDH_MORE_DATA                              = 0x800007D2 // The PdhGetFormattedCounterArray* function can return this if there's 'more data to be displayed'.
	PDH_CSTATUS_ITEM_NOT_VALIDATED             = 0x800007D3
	PDH_RETRY                                  = 0x800007D4
	PDH_NO_DATA                                = 0x800007D5 // The query does not currently contain any counters (for example, limited access)
	PDH_CALC_NEGATIVE_DENOMINATOR              = 0x800007D6
	PDH_CALC_NEGATIVE_TIMEBASE                 = 0x800007D7
	PDH_CALC_NEGATIVE_VALUE                    = 0x800007D8
	PDH_DIALOG_CANCELLED                       = 0x800007D9
	PDH_END_OF_LOG_FILE                        = 0x800007DA
	PDH_ASYNC_QUERY_TIMEOUT                    = 0x800007DB
	PDH_CANNOT_SET_DEFAULT_REALTIME_DATASOURCE = 0x800007DC
	PDH_CSTATUS_NO_OBJECT                      = 0xC0000BB8
	PDH_CSTATUS_NO_COUNTER                     = 0xC0000BB9 // The specified counter could not be found.
	PDH_CSTATUS_INVALID_DATA                   = 0xC0000BBA // The counter was successfully found, but the data returned is not valid.
	PDH_MEMORY_ALLOCATION_FAILURE              = 0xC0000BBB
	PDH_INVALID_HANDLE                         = 0xC0000BBC
	PDH_INVALID_ARGUMENT                       = 0xC0000BBD // Required argument is missing or incorrect.
	PDH_FUNCTION_NOT_FOUND                     = 0xC0000BBE
	PDH_CSTATUS_NO_COUNTERNAME                 = 0xC0000BBF
	PDH_CSTATUS_BAD_COUNTERNAME                = 0xC0000BC0 // Unable to parse the counter path. Check the format and syntax of the specified path.
	PDH_INVALID_BUFFER                         = 0xC0000BC1
	PDH_INSUFFICIENT_BUFFER                    = 0xC0000BC2
	PDH_CANNOT_CONNECT_MACHINE                 = 0xC0000BC3
	PDH_INVALID_PATH                           = 0xC0000BC4
	PDH_INVALID_INSTANCE                       = 0xC0000BC5
	PDH_INVALID_DATA                           = 0xC0000BC6 // specified counter does not contain valid data or a successful status code.
	PDH_NO_DIALOG_DATA                         = 0xC0000BC7
	PDH_CANNOT_READ_NAME_STRINGS               = 0xC0000BC8
	PDH_LOG_FILE_CREATE_ERROR                  = 0xC0000BC9
	PDH_LOG_FILE_OPEN_ERROR                    = 0xC0000BCA
	PDH_LOG_TYPE_NOT_FOUND                     = 0xC0000BCB
	PDH_NO_MORE_DATA                           = 0xC0000BCC
	PDH_ENTRY_NOT_IN_LOG_FILE                  = 0xC0000BCD
	PDH_DATA_SOURCE_IS_LOG_FILE                = 0xC0000BCE
	PDH_DATA_SOURCE_IS_REAL_TIME               = 0xC0000BCF
	PDH_UNABLE_READ_LOG_HEADER                 = 0xC0000BD0
	PDH_FILE_NOT_FOUND                         = 0xC0000BD1
	PDH_FILE_ALREADY_EXISTS                    = 0xC0000BD2
	PDH_NOT_IMPLEMENTED                        = 0xC0000BD3
	PDH_STRING_NOT_FOUND                       = 0xC0000BD4
	PDH_UNABLE_MAP_NAME_FILES                  = 0x80000BD5
	PDH_UNKNOWN_LOG_FORMAT                     = 0xC0000BD6
	PDH_UNKNOWN_LOGSVC_COMMAND                 = 0xC0000BD7
	PDH_LOGSVC_QUERY_NOT_FOUND                 = 0xC0000BD8
	PDH_LOGSVC_NOT_OPENED                      = 0xC0000BD9
	PDH_WBEM_ERROR                             = 0xC0000BDA
	PDH_ACCESS_DENIED                          = 0xC0000BDB
	PDH_LOG_FILE_TOO_SMALL                     = 0xC0000BDC
	PDH_INVALID_DATASOURCE                     = 0xC0000BDD
	PDH_INVALID_SQLDB                          = 0xC0000BDE
	PDH_NO_COUNTERS                            = 0xC0000BDF
	PDH_SQL_ALLOC_FAILED                       = 0xC0000BE0
	PDH_SQL_ALLOCCON_FAILED                    = 0xC0000BE1
	PDH_SQL_EXEC_DIRECT_FAILED                 = 0xC0000BE2
	PDH_SQL_FETCH_FAILED                       = 0xC0000BE3
	PDH_SQL_ROWCOUNT_FAILED                    = 0xC0000BE4
	PDH_SQL_MORE_RESULTS_FAILED                = 0xC0000BE5
	PDH_SQL_CONNECT_FAILED                     = 0xC0000BE6
	PDH_SQL_BIND_FAILED                        = 0xC0000BE7
	PDH_CANNOT_CONNECT_WMI_SERVER              = 0xC0000BE8
	PDH_PLA_COLLECTION_ALREADY_RUNNING         = 0xC0000BE9
	PDH_PLA_ERROR_SCHEDULE_OVERLAP             = 0xC0000BEA
	PDH_PLA_COLLECTION_NOT_FOUND               = 0xC0000BEB
	PDH_PLA_ERROR_SCHEDULE_ELAPSED             = 0xC0000BEC
	PDH_PLA_ERROR_NOSTART                      = 0xC0000BED
	PDH_PLA_ERROR_ALREADY_EXISTS               = 0xC0000BEE
	PDH_PLA_ERROR_TYPE_MISMATCH                = 0xC0000BEF
	PDH_PLA_ERROR_FILEPATH                     = 0xC0000BF0
	PDH_PLA_SERVICE_ERROR                      = 0xC0000BF1
	PDH_PLA_VALIDATION_ERROR                   = 0xC0000BF2
	PDH_PLA_VALIDATION_WARNING                 = 0x80000BF3
	PDH_PLA_ERROR_NAME_TOO_LONG                = 0xC0000BF4
	PDH_INVALID_SQL_LOG_FORMAT                 = 0xC0000BF5
	PDH_COUNTER_ALREADY_IN_QUERY               = 0xC0000BF6
	PDH_BINARY_LOG_CORRUPT                     = 0xC0000BF7
	PDH_LOG_SAMPLE_TOO_SMALL                   = 0xC0000BF8
	PDH_OS_LATER_VERSION                       = 0xC0000BF9
	PDH_OS_EARLIER_VERSION                     = 0xC0000BFA
	PDH_INCORRECT_APPEND_TIME                  = 0xC0000BFB
	PDH_UNMATCHED_APPEND_COUNTER               = 0xC0000BFC
	PDH_SQL_ALTER_DETAIL_FAILED                = 0xC0000BFD
	PDH_QUERY_PERF_DATA_TIMEOUT                = 0xC0000BFE
)

// Formatting options for GetFormattedCounterValue().
const (
	PDH_FMT_DOUBLE = 0x00000200 // Return data as a double precision floating point real.
)

type (
	HANDLE uintptr // query handle
)

// const (
// 	PDH_FMT_COUNTERVALUE = C.PDH_FMT_COUNTERVALUE
// )

// // Union specialization for double values
type PDH_FMT_COUNTERVALUE_DOUBLE struct {
	CStatus     uint32
	DoubleValue float64
}

// // Union specialization for double values, used by PdhGetFormattedCounterArrayDouble()
type PDH_FMT_COUNTERVALUE_ITEM_DOUBLE struct {
	SzName   *uint16 // pointer to a string
	FmtValue PDH_FMT_COUNTERVALUE_DOUBLE
}

var (
	// Library
	libpdhDll *syscall.DLL

	// Functions
	pdh_AddCounterW               *syscall.Proc
	pdh_AddEnglishCounterW        *syscall.Proc
	pdh_CloseQuery                *syscall.Proc
	pdh_CollectQueryData          *syscall.Proc
	pdh_GetFormattedCounterValue  *syscall.Proc
	pdh_GetFormattedCounterArrayW *syscall.Proc
	pdh_OpenQuery                 *syscall.Proc
	pdh_ValidatePathW             *syscall.Proc
)

func init() {
	// Library
	libpdhDll = syscall.MustLoadDLL("pdh.dll")

	// Functions
	pdh_AddCounterW = libpdhDll.MustFindProc("PdhAddCounterW")
	pdh_AddEnglishCounterW, _ = libpdhDll.FindProc("PdhAddEnglishCounterW") // XXX: only supported on versions > Vista.
	pdh_CloseQuery = libpdhDll.MustFindProc("PdhCloseQuery")
	pdh_CollectQueryData = libpdhDll.MustFindProc("PdhCollectQueryData")
	pdh_GetFormattedCounterValue = libpdhDll.MustFindProc("PdhGetFormattedCounterValue")
	pdh_GetFormattedCounterArrayW = libpdhDll.MustFindProc("PdhGetFormattedCounterArrayW")
	pdh_OpenQuery = libpdhDll.MustFindProc("PdhOpenQuery")
	pdh_ValidatePathW = libpdhDll.MustFindProc("PdhValidatePathW")
}

// Examples of szFullCounterPath (in an English version of Windows):
//
//  \\Processor(_Total)\\% Idle Time
//  \\Processor(_Total)\\% Processor Time
//  \\LogicalDisk(C:)\% Free Space

// The typeperf command may also be pretty easy. To find all performance counters, simply execute:
//
//  typeperf -qx

// Adds the specified language-neutral counter to the query. See the PdhAddCounter function. This function only exists on
// Windows versions higher than Vista.
func PdhAddEnglishCounter(hQuery uintptr, szFullCounterPath string, dwUserData uintptr, phCounter *uintptr) uint32 {
	//TODO: ERROR_INVALID_FUNCTION
	// if pdh_AddEnglishCounterW == nil {
	// 	return ERROR_INVALID_FUNCTION
	// }

	ptxt, _ := syscall.UTF16PtrFromString(szFullCounterPath)
	ret, _, _ := pdh_AddEnglishCounterW.Call(
		uintptr(hQuery),
		uintptr(unsafe.Pointer(ptxt)),
		dwUserData,
		uintptr(unsafe.Pointer(phCounter)))

	return uint32(ret)
}

// Closes all counters contained in the specified query, closes all handles related to the query,
// and frees all memory associated with the query.
func PdhCloseQuery(hQuery uintptr) uint32 {
	ret, _, _ := pdh_CloseQuery.Call(uintptr(hQuery))

	return uint32(ret)
}

// The PdhCollectQueryData will return an error in the first call because it needs two values for
// displaying the correct data for the processor idle time. The second call will have a 0 return code.
func PdhCollectQueryData(hQuery uintptr) uint32 {
	ret, _, _ := pdh_CollectQueryData.Call(uintptr(hQuery))

	return uint32(ret)
}

// Formats the given hCounter using a 'double'. The result is set into the specialized union struct pValue.
// This function does not directly translate to a Windows counterpart due to union specialization tricks.
func PdhGetFormattedCounterValueDouble(hCounter uintptr, lpdwType uint32, pValue *PDH_FMT_COUNTERVALUE_DOUBLE) uint32 {
	var pv C.PDH_FMT_COUNTERVALUE_DOUBLE
	ret, _, _ := pdh_GetFormattedCounterValue.Call(
		uintptr(hCounter),
		uintptr(PDH_FMT_DOUBLE),
		uintptr(unsafe.Pointer(&lpdwType)),
		uintptr(unsafe.Pointer(&pv)))
	// fmt.Println(ret, pValue)
	// return 0, PDH_FMT_COUNTERVALUE_DOUBLE{}
	if C.ulong(pv.CStatus) == 0 {
		// return uint32(ret), nil
		pValue.CStatus = uint32(C.ulong(pv.CStatus))
		pValue.DoubleValue = float64(C.double(pv.DoubleValue))
	}
	return uint32(ret)
}

// Creates a new query that is used to manage the collection of performance data.
// szDataSource is a null terminated string that specifies the name of the log file from which to
// retrieve the performance data. If 0, performance data is collected from a real-time data source.
// dwUserData is a user-defined value to associate with this query. To retrieve the user data later,
// call PdhGetCounterInfo and access dwQueryUserData of the PDH_COUNTER_INFO structure. phQuery is
// the handle to the query, and must be used in subsequent calls. This function returns a PDH_
// constant error code, or ERROR_SUCCESS if the call succeeded.
func PdhOpenQuery(szDataSource uintptr, dwUserData uintptr, phQuery *uintptr) uint32 {
	ret, _, _ := pdh_OpenQuery.Call(
		szDataSource,
		dwUserData,
		uintptr(unsafe.Pointer(phQuery)))
	return uint32(ret)
}

// Validates a path. Will return ERROR_SUCCESS when ok, or PDH_CSTATUS_BAD_COUNTERNAME when the path is
// erroneous.
func PdhValidatePath(path string) uint32 {
	ptxt, _ := syscall.UTF16PtrFromString(path)
	ret, _, _ := pdh_ValidatePathW.Call(uintptr(unsafe.Pointer(ptxt)))

	return uint32(ret)
}

func UTF16PtrToString(s *uint16) string {
	if s == nil {
		return ""
	}
	return syscall.UTF16ToString((*[1 << 29]uint16)(unsafe.Pointer(s))[0:])
}

// --------------- wrapper code ------

type PdhDataPoint struct {
	Query     string
	Timestamp time.Time
	Value     float64
}

type PdhCollector struct {
	handle   uintptr
	counters map[string]uintptr
}

func NewPdhCollector() *PdhCollector {
	var handle uintptr
	PdhOpenQuery(0, 0, &handle)

	return &PdhCollector{
		handle:   handle,
		counters: make(map[string]uintptr),
	}
}

func (p *PdhCollector) GetHandle() uintptr {
	return p.handle
}

func (p *PdhCollector) Close() {
	PdhCloseQuery(p.handle)
}

func (p *PdhCollector) AddEnglishCounter(query string) {
	var handle uintptr
	PdhAddEnglishCounter(p.handle, query, 0, &handle)
	p.counters[query] = handle
}

func valid_pdh_cstatus(cs uint32) bool {
	if cs == uint32(PDH_CSTATUS_VALID_DATA) || cs == uint32(PDH_CSTATUS_NEW_DATA) {
		return true
	}
	return false
}

func (p *PdhCollector) CollectData() []*PdhDataPoint {
	PdhCollectQueryData(p.handle)
	data := []*PdhDataPoint{}

	var perf PDH_FMT_COUNTERVALUE_DOUBLE
	for key, chandle := range p.counters {
		cstatus := PdhValidatePath(key)
		if valid_pdh_cstatus(cstatus) == true {
			PdhGetFormattedCounterValueDouble(chandle, 0, &perf)
			if valid_pdh_cstatus(perf.CStatus) == true {
				pd := PdhDataPoint{
					Query:     key,
					Timestamp: time.Now(),
					Value:     perf.DoubleValue,
				}
				data = append(data, &pd)
			} else {
				log.Debugf("invalid data: CSTATUS: %x", perf.CStatus)
			}
		} else {
			log.Debugf("invalid path: CSTATUS: %x Path: %s", cstatus, key)
		}
	}
	return data
}

func init() {
	collector_factories["win_pdh"] = factory_win_pdh

	builtin_collectors = append(builtin_collectors, builtin_win_pdh())
}

func builtin_win_pdh() Collector {

	queries := []config.Conf_win_pdh_query{}

	queries = append(queries, config.Conf_win_pdh_query{
		Query:  "\\System\\Processes",
		Metric: "win.processes.count"})
	queries = append(queries, config.Conf_win_pdh_query{
		Query:  "\\Memory\\Available Bytes",
		Metric: "win.memory.available_bytes"})
	queries = append(queries, config.Conf_win_pdh_query{
		Query:  "\\Process(_Total)\\Working Set",
		Metric: "win.process.working_set.total"})
	queries = append(queries, config.Conf_win_pdh_query{
		Query:  "\\Memory\\Cache Bytes",
		Metric: "win.memory.cache_bytes"})

	// private memory size
	queries = append(queries, config.Conf_win_pdh_query{
		Query:  "\\Process(hickwall)\\Working Set - Private",
		Metric: "hickwall.client.mem.private_working_set"})
	queries = append(queries, config.Conf_win_pdh_query{
		Query:  "\\Process(hickwall)\\Working Set",
		Metric: "hickwall.client.mem.working_set"})

	//FIXME: temp fix double process problem.
	// queries = append(queries, config.Conf_win_pdh_query{
	// 	Query:  "\\Process(hickwall#1)\\Working Set - Private",
	// 	Metric: "hickwall.client.mem.private_working_set.1"})
	// queries = append(queries, config.Conf_win_pdh_query{
	// 	Query:  "\\Process(hickwall#1)\\Working Set",
	// 	Metric: "hickwall.client.mem.working_set.1"})

	conf := config.Conf_win_pdh{
		Interval: "2s",
		Queries:  queries,
	}

	return factory_win_pdh("builtin_win_pdh", conf)
}

func factory_win_pdh(name string, conf interface{}) Collector {
	var states state_win_pdh
	var cf config.Conf_win_pdh
	var default_interval = time.Duration(1) * time.Second

	if conf != nil {
		cf = conf.(config.Conf_win_pdh)
		// fmt.Println("factory_win_pdh: ", cf)
		// pretty.Println("factory_win_pdh:", cf)

		interval, err := collectorlib.ParseInterval(cf.Interval)
		if err != nil {
			log.Errorf("cannot parse interval of collector_pdh: %s - %v", cf.Interval, err)
			interval = default_interval
		}
		states.Interval = interval

		states.hPdh = NewPdhCollector()

		states.map_queries = make(map[string]config.Conf_win_pdh_query)

		for _, query_obj := range cf.Queries {
			query := query_obj.Query
			//TODO: validate query

			states.hPdh.AddEnglishCounter(query)

			query_obj.Tags = AddTags.Copy().Merge(config.Conf.Tags).Merge(cf.Tags).Merge(query_obj.Tags)

			states.map_queries[query] = query_obj
		}
	}

	return &IntervalCollector{
		F:        c_win_pdh,
		Enable:   nil,
		name:     name,
		states:   states,
		Interval: states.Interval,
	}
}

type state_win_pdh struct {
	Interval time.Duration

	// internal use only
	hPdh        *PdhCollector
	map_queries map[string]config.Conf_win_pdh_query
}

func c_win_pdh(states interface{}) (collectorlib.MultiDataPoint, error) {
	var md collectorlib.MultiDataPoint
	var st state_win_pdh

	if states != nil {
		st = states.(state_win_pdh)
		// fmt.Println("c_win_pdh states: ", states)
	}

	if st.hPdh != nil {

		data := st.hPdh.CollectData()
		queries := st.map_queries

		for _, pd := range data {
			query := queries[pd.Query]

			Add(&md, query.Metric, pd.Value, query.Tags, "", "", "")
		}
	}
	return md, nil
}
