package cgo_pdh

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lpdh
#include "c\cgo_pdh.c"
*/
import "C"

import (
	"errors"
	"fmt"
	"math"
	"unsafe"
)

func getcpuload() {
	fmt.Println(C.getcpuload())
}

type pdhCollector struct {
	hQuery   C.PDH_HQUERY
	counters []C.PDH_HCOUNTER
	queries  []string
}

func NewPdhCollector() *pdhCollector {
	p := &pdhCollector{}
	p.OpenQuery()
	return p
}

func (p *pdhCollector) OpenQuery() (err error) {
	var hStatus C.PDH_STATUS
	hStatus = C.PdhOpenQuery(nil, 0, &p.hQuery)
	err = ValidateCStatus(hStatus)
	if err != nil {
		return fmt.Errorf("failed to open query: %v", err)
	}
	return nil
}

func (p *pdhCollector) Close() (err error) {
	var hStatus C.PDH_STATUS
	hStatus = C.PdhCloseQuery(p.hQuery)
	err = ValidateCStatus(hStatus)
	if err != nil {
		return fmt.Errorf("failed to close query: %v", err)
	}
	return nil
}

func (p *pdhCollector) AddCounter(query string) (err error) {
	var hStatus C.PDH_STATUS
	var hCounter C.PDH_HCOUNTER

	query_str := (*C.CHAR)(C.CString(query))
	defer C.free(unsafe.Pointer(query_str))

	hStatus = C.PdhAddCounter(p.hQuery, query_str, 0, &hCounter)
	err = ValidateCStatus(hStatus)
	if err != nil {
		return fmt.Errorf("failed to AddCounter: %v", err)
	}

	p.counters = append(p.counters, hCounter)
	p.queries = append(p.queries, query)
	return nil
}

func (p *pdhCollector) CollectAllDouble() (res []float64, err error) {
	var hStatus C.PDH_STATUS
	var value = C.double(0.0)

	C.PdhCollectQueryData(p.hQuery)
	hStatus = C.PdhCollectQueryData(p.hQuery)
	err = ValidateCStatus(hStatus)
	if err != nil {
		return res, err
	}

	for _, h := range p.counters {
		hStatus = C.GetDoubleCounterValue(h, &value)
		err = ValidateCStatus(hStatus)
		if err != nil {
			res = append(res, math.NaN())
		} else {
			res = append(res, (float64)(value))
		}

	}
	return res, nil
}

func (p *pdhCollector) GetQueryByIdx(idx int) string {
	return p.queries[idx][:]
}

func ValidateCStatus(hStatus C.PDH_STATUS) error {
	switch uint32(hStatus) {
	case 0x0: //C.ERROR_SUCCESS, PDH_CSTATUS_VALID_DATA
		return nil
	case 0x00000001:
		return errors.New("PDH_CSTATUS_NEW_DATA")
	case 0x800007D0:
		return errors.New("PDH_CSTATUS_NO_MACHINE")
	case 0x800007D1:
		return errors.New("PDH_CSTATUS_NO_INSTANCE")
	case 0x800007D2:
		return errors.New("PDH_MORE_DATA")
	case 0x800007D3:
		return errors.New("PDH_CSTATUS_ITEM_NOT_VALIDATED")
	case 0x800007D4:
		return errors.New("PDH_RETRY")
	case 0x800007D5:
		return errors.New("PDH_NO_DATA")
	case 0x800007D6:
		return errors.New("PDH_CALC_NEGATIVE_DENOMINATOR")
	case 0x800007D7:
		return errors.New("PDH_CALC_NEGATIVE_TIMEBASE")
	case 0x800007D8:
		return errors.New("PDH_CALC_NEGATIVE_VALUE")
	case 0x800007D9:
		return errors.New("PDH_DIALOG_CANCELLED")
	case 0x800007DA:
		return errors.New("PDH_END_OF_LOG_FILE")
	case 0x800007DB:
		return errors.New("PDH_ASYNC_QUERY_TIMEOUT")
	case 0x800007DC:
		return errors.New("PDH_CANNOT_SET_DEFAULT_REALTIME_DATASOURCE")
	case 0xC0000BB8:
		return errors.New("PDH_CSTATUS_NO_OBJECT")
	case 0xC0000BB9:
		return errors.New("PDH_CSTATUS_NO_COUNTER")
	case 0xC0000BBA:
		return errors.New("PDH_CSTATUS_INVALID_DATA")
	case 0xC0000BBB:
		return errors.New("PDH_MEMORY_ALLOCATION_FAILURE")
	case 0xC0000BBC:
		return errors.New("PDH_INVALID_HANDLE")
	case 0xC0000BBD:
		return errors.New("PDH_INVALID_ARGUMENT")
	case 0xC0000BBE:
		return errors.New("PDH_FUNCTION_NOT_FOUND")
	case 0xC0000BBF:
		return errors.New("PDH_CSTATUS_NO_COUNTERNAME")
	case 0xC0000BC0:
		return errors.New("PDH_CSTATUS_BAD_COUNTERNAME")
	case 0xC0000BC1:
		return errors.New("PDH_INVALID_BUFFER")
	case 0xC0000BC2:
		return errors.New("PDH_INSUFFICIENT_BUFFER")
	case 0xC0000BC3:
		return errors.New("PDH_CANNOT_CONNECT_MACHINE")
	case 0xC0000BC4:
		return errors.New("PDH_INVALID_PATH")
	case 0xC0000BC5:
		return errors.New("PDH_INVALID_INSTANCE")
	case 0xC0000BC6:
		return errors.New("PDH_INVALID_DATA")
	case 0xC0000BC7:
		return errors.New("PDH_NO_DIALOG_DATA")
	case 0xC0000BC8:
		return errors.New("PDH_CANNOT_READ_NAME_STRINGS")
	case 0xC0000BC9:
		return errors.New("PDH_LOG_FILE_CREATE_ERROR")
	case 0xC0000BCA:
		return errors.New("PDH_LOG_FILE_OPEN_ERROR")
	case 0xC0000BCB:
		return errors.New("PDH_LOG_TYPE_NOT_FOUND")
	case 0xC0000BCC:
		return errors.New("PDH_NO_MORE_DATA")
	case 0xC0000BCD:
		return errors.New("PDH_ENTRY_NOT_IN_LOG_FILE")
	case 0xC0000BCE:
		return errors.New("PDH_DATA_SOURCE_IS_LOG_FILE")
	case 0xC0000BCF:
		return errors.New("PDH_DATA_SOURCE_IS_REAL_TIME")
	case 0xC0000BD0:
		return errors.New("PDH_UNABLE_READ_LOG_HEADER")
	case 0xC0000BD1:
		return errors.New("PDH_FILE_NOT_FOUND")
	case 0xC0000BD2:
		return errors.New("PDH_FILE_ALREADY_EXISTS")
	case 0xC0000BD3:
		return errors.New("PDH_NOT_IMPLEMENTED")
	case 0xC0000BD4:
		return errors.New("PDH_STRING_NOT_FOUND")
	case 0x80000BD5:
		return errors.New("PDH_UNABLE_MAP_NAME_FILES")
	case 0xC0000BD6:
		return errors.New("PDH_UNKNOWN_LOG_FORMAT")
	case 0xC0000BD7:
		return errors.New("PDH_UNKNOWN_LOGSVC_COMMAND")
	case 0xC0000BD8:
		return errors.New("PDH_LOGSVC_QUERY_NOT_FOUND")
	case 0xC0000BD9:
		return errors.New("PDH_LOGSVC_NOT_OPENED")
	case 0xC0000BDA:
		return errors.New("PDH_WBEM_ERROR")
	case 0xC0000BDB:
		return errors.New("PDH_ACCESS_DENIED")
	case 0xC0000BDC:
		return errors.New("PDH_LOG_FILE_TOO_SMALL")
	case 0xC0000BDD:
		return errors.New("PDH_INVALID_DATASOURCE")
	case 0xC0000BDE:
		return errors.New("PDH_INVALID_SQLDB")
	case 0xC0000BDF:
		return errors.New("PDH_NO_COUNTERS")
	case 0xC0000BE0:
		return errors.New("PDH_SQL_ALLOC_FAILED")
	case 0xC0000BE1:
		return errors.New("PDH_SQL_ALLOCCON_FAILED")
	case 0xC0000BE2:
		return errors.New("PDH_SQL_EXEC_DIRECT_FAILED")
	case 0xC0000BE3:
		return errors.New("PDH_SQL_FETCH_FAILED")
	case 0xC0000BE4:
		return errors.New("PDH_SQL_ROWCOUNT_FAILED")
	case 0xC0000BE5:
		return errors.New("PDH_SQL_MORE_RESULTS_FAILED")
	case 0xC0000BE6:
		return errors.New("PDH_SQL_CONNECT_FAILED")
	case 0xC0000BE7:
		return errors.New("PDH_SQL_BIND_FAILED")
	case 0xC0000BE8:
		return errors.New("PDH_CANNOT_CONNECT_WMI_SERVER")
	case 0xC0000BE9:
		return errors.New("PDH_PLA_COLLECTION_ALREADY_RUNNING")
	case 0xC0000BEA:
		return errors.New("PDH_PLA_ERROR_SCHEDULE_OVERLAP")
	case 0xC0000BEB:
		return errors.New("PDH_PLA_COLLECTION_NOT_FOUND")
	case 0xC0000BEC:
		return errors.New("PDH_PLA_ERROR_SCHEDULE_ELAPSED")
	case 0xC0000BED:
		return errors.New("PDH_PLA_ERROR_NOSTART")
	case 0xC0000BEE:
		return errors.New("PDH_PLA_ERROR_ALREADY_EXISTS")
	case 0xC0000BEF:
		return errors.New("PDH_PLA_ERROR_TYPE_MISMATCH")
	case 0xC0000BF0:
		return errors.New("PDH_PLA_ERROR_FILEPATH")
	case 0xC0000BF1:
		return errors.New("PDH_PLA_SERVICE_ERROR")
	case 0xC0000BF2:
		return errors.New("PDH_PLA_VALIDATION_ERROR")
	case 0x80000BF3:
		return errors.New("PDH_PLA_VALIDATION_WARNING")
	case 0xC0000BF4:
		return errors.New("PDH_PLA_ERROR_NAME_TOO_LONG")
	case 0xC0000BF5:
		return errors.New("PDH_INVALID_SQL_LOG_FORMAT")
	case 0xC0000BF6:
		return errors.New("PDH_COUNTER_ALREADY_IN_QUERY")
	case 0xC0000BF7:
		return errors.New("PDH_BINARY_LOG_CORRUPT")
	case 0xC0000BF8:
		return errors.New("PDH_LOG_SAMPLE_TOO_SMALL")
	case 0xC0000BF9:
		return errors.New("PDH_OS_LATER_VERSION")
	case 0xC0000BFA:
		return errors.New("PDH_OS_EARLIER_VERSION")
	case 0xC0000BFB:
		return errors.New("PDH_INCORRECT_APPEND_TIME")
	case 0xC0000BFC:
		return errors.New("PDH_UNMATCHED_APPEND_COUNTER")
	case 0xC0000BFD:
		return errors.New("PDH_SQL_ALTER_DETAIL_FAILED")
	case 0xC0000BFE:
		return errors.New("PDH_QUERY_PERF_DATA_TIMEOUT")
	}
	return errors.New("Unknown Error")
}
