// +build windows

package main

import (
	"fmt"
	"github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"
	"regexp"
	"strings"
	"time"
)

const (
	MaxUint32 = ^uint32(0)
	MaxInt32  = int32(MaxUint32 >> 1)
)

func ParseFieldsFromQuery(query string) []string {
	results := []string{}

	re, _ := regexp.Compile("select (.*) from")

	matched_fields := re.FindStringSubmatch(strings.ToLower(query))

	if len(matched_fields) != 2 {
		return []string{}
	}

	fields := strings.Split(matched_fields[1], ",")
	// for idx, pat := range fields {
	for _, pat := range fields {
		field := strings.Trim(pat, " ")
		// fmt.Println(idx, field)
		results = append(results, field)
	}

	return results
}

func WmiQueryWithFields(query string, fields []string) []map[string]string {
	// init COM, oh yeah
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, _ := oleutil.CreateObject("WbemScripting.SWbemLocator")
	defer unknown.Release()

	wmi, _ := unknown.QueryInterface(ole.IID_IDispatch)
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, _ := oleutil.CallMethod(wmi, "ConnectServer")
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// result is a SWBemObjectSet
	resultRaw, _ := oleutil.CallMethod(service, "ExecQuery", query)

	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, _ := oleutil.GetProperty(result, "Count")
	count := int(countVar.Val)

	//
	resultMap := []map[string]string{}

	fmt.Println("Count: ", count)
	for i := 0; i < count; i++ {
		itemMap := make(map[string]string)

		// item is a SWbemObject, but really a Win32_Process
		itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i)

		item := itemRaw.ToIDispatch()
		defer item.Release()

		for _, field := range fields {
			asString, err := oleutil.GetProperty(item, field)

			fmt.Println(field, asString, err)
			if err == nil {
				// itemMap[field] = asString.ToString()
				itemMap[field] = fmt.Sprintf("%v", asString.Value())
			} else {
				fmt.Println(err)
			}
		}

		// fmt.Println(itemMap)
		resultMap = append(resultMap, itemMap)
	}

	return resultMap
}

/*
return  results, count of results
*/
func WmiQuery(query string) []map[string]string {
	fields := ParseFieldsFromQuery(query)
	// init COM, oh yeah
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, _ := oleutil.CreateObject("WbemScripting.SWbemLocator")
	defer unknown.Release()

	wmi, _ := unknown.QueryInterface(ole.IID_IDispatch)
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, _ := oleutil.CallMethod(wmi, "ConnectServer")
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// result is a SWBemObjectSet
	resultRaw, _ := oleutil.CallMethod(service, "ExecQuery", query)

	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, _ := oleutil.GetProperty(result, "Count")
	count := int(countVar.Val)

	//
	resultMap := []map[string]string{}

	fmt.Println("Count: ", count)
	for i := 0; i < count; i++ {
		itemMap := make(map[string]string)

		// item is a SWbemObject, but really a Win32_Process
		itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i)

		item := itemRaw.ToIDispatch()
		defer item.Release()

		for _, field := range fields {
			asString, err := oleutil.GetProperty(item, field)
			// asString, err := oleutil.GetProperty(item, "NumberOfLogicalProcessors")

			// asString may return invalid pointer: drivetype &{3 0 0 0 8463721106786746371}

			// fmt.Println(field, asString, err, field_string)
			if err == nil {
				// itemMap[field] = asString.ToString()
				itemMap[field] = fmt.Sprintf("%v", asString.Value())
			} else {
				fmt.Println(err)
			}
		}

		// fmt.Println(itemMap)
		resultMap = append(resultMap, itemMap)
	}

	return resultMap
}

func main() {
	// query := "select Caption, FreePhysicalMemory, TotalVirtualMemorySize from win32_operatingsystem"

	// query := "select Caption, NumberOfCores from Win32_Processor"
	// query := "select Caption, NumberOfProcessors, NumberOfLogicalProcessors from Win32_ComputerSystem"

	// don't support *
	// query := "select * from Win32_ComputerSystem"
	// query := "select Caption, FreePhysicalMemory, TotalVirtualMemorySize from win32_operatingsystem"
	// query := "select Name, FileSystem, Size, FreeSpace from Win32_LogicalDisk"
	query := "select Name, FileSystem, FreeSpace, Size from Win32_LogicalDisk where MediaType=11 or mediatype=12"

	results := WmiQuery(query)
	// results := WmiQueryWithFields(query, []string{"Caption", "NumberOfProcessors"})

	for _, item := range results {
		fmt.Println(item)
	}
	time.Sleep(time.Millisecond * 1)

	// 	done := time.After(time.Second * 300)
	// 	delay := time.After(time.Second * 1)
	// 	tick := time.Tick(time.Millisecond * 1000)
	// loop:
	// 	for {
	// 		select {
	// 		// case dp, err := <-ch:
	// 		// case <-ch:
	// 		// fmt.Println(" point ---> ", dp, err)
	// 		// fmt.Println("-------------------")
	// 		case <-tick:
	// 			WmiQuery(query)
	// 		case <-delay:
	// 			// fmt.Println("-------------------")
	// 			// change config on the fly
	// 			// cs[0].Init()
	// 			// cs[0].(*IntervalCollector).SetInterval(time.Millisecond * 200)
	// 		case <-done:
	// 			break loop
	// 		}
	// 	}
}
