package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {

	metric := "  win.wmi.fs.D:.CDFS.free space.bytes 中文"

	pat_w, _ := regexp.Compile("\\w+(\\s+\\w+)?")
	pat_s, _ := regexp.Compile("\\s+")

	parts := pat_w.FindAllString(metric, -1)

	for _, part := range parts {
		fmt.Println(part)
	}
	metric = strings.Join(parts, ".")
	fmt.Println(metric)

	metric = pat_s.ReplaceAllString(metric, "_")
	fmt.Println(metric)

	metric = strings.ToLower(metric)
	fmt.Println(metric)
}
