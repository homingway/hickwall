package main

import (
	"fmt"
	"github.com/oliveagle/go-collectors/util"
	"regexp"
	"strings"
	"time"
)

var (
	pat_format, _ = regexp.Compile("\\/format:\\w+(.xsl)?")
	pat_get, _    = regexp.Compile("\\bget\\b")
)

func WmiQueryCmdLineV4(query string) []map[string]string {

	// fmt.Println("query: ", query)
	the_query := strings.ToLower(query)

	if pat_format.MatchString(the_query) == true {
		the_query = pat_format.ReplaceAllString(the_query, "/format:textvaluelist")
	} else {
		if pat_get.MatchString(the_query) == true {
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
		// 	fmt.Println(p)
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

func perf() {
	tick := time.Tick(time.Millisecond * 10)
	done := time.After(time.Second * 1000)
loop:
	for {
		select {
		case <-tick:
			// for performance, should not use go func() here
			util.ReadCommand(func(line string) error {
				// fmt.Println(line)
				fmt.Printf(".")
				return nil
			}, "wmic", "cpu", "get", "/format:textvaluelist.xsl")
		case <-done:
			break loop
		}
	}
}

func main() {
	// perf()

	// should use with limitations.

	// results := WmiQueryCmdLineV4("wmic cpu")
	// results := WmiQueryCmdLineV4("wmic cpu get Name,NumberOfCores,NumberOfLogicalProcessors /format:rawxml")
	// results := WmiQueryCmdLineV4("wmic logicaldisk get name, filesystem, size, FreeSpace ")
	// results := WmiQueryCmdLineV4("wmic logicaldisk ")
	// results := WmiQueryCmdLineV4("wmic nicconfig where 'IPEnabled=TRUE' get MACAddress, DefaultIPGateway, IPAddress, IPSubnet, DNSHostName, DNSDomain")

	// results := WmiQueryCmdLineV4("wmic logicaldisk where 'mediatype=12' get name, filesystem, size, FreeSpace ")
	results := WmiQueryCmdLineV4("wmic logicaldisk where 'mediatype=11 or mediatype=12' get name, filesystem, size, FreeSpace ")

	for _, record := range results {
		fmt.Println("--------")
		fmt.Println(record)
		// pretty.Println(record)
		// for _, item := range record.Properties {
		// 	fmt.Println(item)
		// }
	}
}
