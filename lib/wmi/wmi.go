package wmi

import (
	"fmt"
	"github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"
	"github.com/oliveagle/hickwall/logging"
	"regexp"
)

var wmi_service *ole.IDispatch

func init() {
	ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)

	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		logging.Criticalf("oleutil.CreateObject Failed: %v", err)
		panic(err)
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		logging.Criticalf("QueryInterface Failed: %v", err)
		panic(err)
	}
	defer wmi.Release()

	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer")
	if err != nil {
		logging.Criticalf("Connect to Server Failed: %v", err)
		panic(err)
	}
	wmi_service = serviceRaw.ToIDispatch()
}

var win_wmi_pat_field = regexp.MustCompile(`(?i)select\s+(.*)\s+from`)
var win_wmi_pat_split = regexp.MustCompile(`,\s+`)

func parseFieldsFromQuery(query string) ([]string, error) {
	matches := win_wmi_pat_field.FindAllStringSubmatch(query, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		f_str := matches[0][1]
		return win_wmi_pat_split.Split(f_str, -1), nil
	}

	return nil, fmt.Errorf("didn't match any fields")
}

func QueryWmiFields(query string, fields []string) ([]map[string]string, error) {

	if len(fields) == 1 && fields[0] == "*" {
		logging.Errorf("`select * ` not supported, need to address fields explicitly.")
		return nil, fmt.Errorf("`select * ` not supported, need to address fields explicitly.")
	}

	resultRaw, err := oleutil.CallMethod(wmi_service, "ExecQuery", query)
	if err != nil {
		logging.Error("ExecQuery Failed: ", err)
		return nil, fmt.Errorf("ExecQuery Failed: %v", err)
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, err := oleutil.GetProperty(result, "Count")
	if err != nil {
		logging.Errorf("Get result count Failed: %v", err)
		return nil, fmt.Errorf("Get result count Failed: %v", err)
	}
	count := int(countVar.Val)

	resultMap := []map[string]string{}

	for i := 0; i < count; i++ {
		itemMap := make(map[string]string)

		itemRaw, err := oleutil.CallMethod(result, "ItemIndex", i)
		if err != nil {
			return nil, fmt.Errorf("ItemIndex Failed: %v", err)
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
		logging.Tracef("wmi query result: %+v", itemMap)
	}
	logging.Tracef("wmi query result count: %d", len(resultMap))
	return resultMap, nil
}

func QueryWmi(query string) ([]map[string]string, error) {
	fields, err := parseFieldsFromQuery(query)
	if err != nil {
		logging.Error("cannot parse fields from query: %v", err)
		return nil, err
	}
	return QueryWmiFields(query, fields)
}
