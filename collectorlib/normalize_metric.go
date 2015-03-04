package collectorlib

import (
	// "fmt"
	"regexp"
	"strings"
)

var (
	nmetric_pat_s, _      = regexp.Compile(`\s+|-+|/+|\\+`)
	nmetric_pat_remove, _ = regexp.Compile(`[^\w-_.]`)
)

func NormalizeMetricKey(metric string) string {

	metric = Re_Trim.ReplaceAllString(metric, "")

	metric = nmetric_pat_s.ReplaceAllString(metric, "_")
	// fmt.Printf("*%s*\n", metric)

	metric = nmetric_pat_remove.ReplaceAllString(metric, "")
	// fmt.Printf("*%s*\n", metric)

	metric = strings.ToLower(strings.Trim(metric, "_"))
	// fmt.Printf("*%s*\n", metric)
	return metric
}
