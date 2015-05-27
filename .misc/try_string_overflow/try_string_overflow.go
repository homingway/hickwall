package main

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/utils"
	"time"

	// "strings"
	"math/rand"
)

func main() {

	utils.HttpPprofServe(7070)

	go func() {
		for {
			rand.Seed(time.Now().UnixNano())

			tpl := "{{.Key}}{{.Tags}}"
			key := "meteric"
			tags := map[string]string{
				"bu":    "hotel",
				"value": fmt.Sprintf("%d", rand.Intn(10000)),
			}

			collectorlib.FlatMetricKeyAndTags(tpl, key, tags)
			// metric, _ := collectorlib.FlatMetricKeyAndTags(tpl, key, tags)
			// fmt.Println(metric)

			time.Sleep(time.Millisecond * time.Duration(1))
		}
	}()

	time.Sleep(time.Minute * time.Duration(10))

}
