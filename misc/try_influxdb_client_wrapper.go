package main

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"github.com/oliveagle/hickwall/backends"
	"log"
	"math/rand"
	"time"
)

func mockWrite() client.BatchPoints {
	return client.BatchPoints{
		Database:        "metrics",
		RetentionPolicy: "p1",
		Points:          mockDataPoint(),
	}
}

func mockDataPoint() (points []client.Point) {

	rand.Seed(time.Now().UTC().UnixNano())
	// tags := map[string]string{"bu": "hotel", "global": "tag", "host": "oliveaglec841"}
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("data%d", i)
		fields := map[string]interface{}{}

		fields["value"] = rand.Intn(100)

		point := client.Point{
			Name: name,
			// Tags:      tags,
			Timestamp: time.Now(),
			Fields:    fields,
			Precision: "s",
		}
		log.Println("point: ", point)
		points = append(points, point)
	}
	return
}

func try_090() {
	iclient, err := backends.NewInfluxdbClient(map[string]interface{}{
		"URL":       "http://192.168.59.103:8086/write",
		"Username":  "root",
		"Password":  "root",
		"UserAgent": "",
	}, "0.9.0-rc7")
	fmt.Println("InfluxdbClient:", iclient, err)

	t, v, err := iclient.Ping()
	fmt.Println(t, v, err)

	// ----------- write ---------------
	write := mockWrite()
	// pretty.Println(write)
	fmt.Println(len(write.Points))

	res, err := iclient.Write(write)
	fmt.Println(res, err)

	// -------------- query -------------------
	res, err = iclient.Query(client.Query{
		Command:  "select * from data1",
		Database: "metrics",
	})
	if err != nil {
		fmt.Println("err", err)
		return
	}
	if res != nil {
		// fmt.Println("len res.Results: ", len(res.Results))
		for _, r := range res.Results {
			for _, s := range r.Series {
				for _, v := range s.Values {
					fmt.Println(v)
				}
			}
		}
	} else {
		// fmt.Println("res == nil ")
	}
}

func try_088() {
	iclient, err := backends.NewInfluxdbClient(map[string]interface{}{
		"Host":     "192.168.59.103:8086",
		"Username": "root",
		"Password": "root",
		"Database": "metrics",
	}, "0.8.8")
	if err != nil {
		fmt.Println("InfluxdbClient: Error: ", err)
		return
	}
	fmt.Println("InfluxdbClient:", iclient, err)

	t, v, err := iclient.Ping()
	fmt.Println(t, v, err)

	// ----------- write ---------------
	write := mockWrite()
	// pretty.Println(write)
	fmt.Println(len(write.Points))

	res, err := iclient.Write(write)
	fmt.Println(res, err)

	// -------------- query -------------------
	res, err = iclient.Query(client.Query{
		Command:  "select * from data1",
		Database: "metrics",
	})
	if err != nil {
		fmt.Println("err", err)
		return
	}
	if res != nil {
		// fmt.Println("len res.Results: ", len(res.Results))
		for _, r := range res.Results {
			for _, s := range r.Series {
				for _, v := range s.Values {
					fmt.Println(v)
				}
			}
		}
	} else {
		// fmt.Println("res == nil ")
	}
}
func main() {
	fmt.Println("---")
	// try_090()
	try_088()

}
