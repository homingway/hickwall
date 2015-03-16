package collectorlib

import (
	// "fmt"
	"regexp"
	"strings"
)

var (
	nmetric_pat_s, _      = regexp.Compile(`\s+|-+|/+|\\+|_+|,|;`)
	nmetric_pat_remove, _ = regexp.Compile(`[^\w-_.]`)
)

func NormalizeMetricKey(metric string) string {

	metric = Re_Trim.ReplaceAllString(metric, "")

	metric = nmetric_pat_s.ReplaceAllString(metric, "_")

	metric = nmetric_pat_remove.ReplaceAllString(metric, "")

	metric = nmetric_pat_s.ReplaceAllString(metric, "_")

	metric = strings.Trim(metric, "_")

	metric = strings.ToLower(metric)

	// this line will coz heap growing slowly.
	// metric = strings.ToLower(strings.Trim(metric, "_"))

	// fmt.Printf("*%s*\n", metric)
	return metric
}
