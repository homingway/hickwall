package main

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"github.com/kr/pretty"
	"github.com/oliveagle/hickwall/config"
	"math/rand"
	// "strings"
	"net/url"
	"time"
)

// point --->  &{win.wmi.cpu.name 1425524770 Intel(R) Core(TM) i5-2435M CPU @ 2.40GHz {bu=hotel,global=tag,host=oliveaglec841}} true
// point --->  &{win.wmi.cpu.numberofcores 1425524770 2 {bu=hotel,global=tag,host=oliveaglec841}} true
// point --->  &{win.wmi.cpu.numberoflogicalprocessors 1425524770 2 {bu=hotel,global=tag,host=oliveaglec841}} true
// point --->  &{win.wmi.mem.totalphysicalmemory 1425524770 4294492160 {bu=hotel,global=tag,host=oliveaglec841}} true
// point --->  &{win.wmi.net.domain 1425524770 WORKGROUP {bu=hotel,global=tag,host=oliveaglec841}} true
// point --->  &{win.wmi.fs.size.c.bytes 1425524770 68350373888 {bu=hotel,fs_type=no_value,global=tag,host=oliveaglec841,mount=C}} true
// point --->  &{win.wmi.fs.size.d.bytes 1425524770 9850880 {bu=hotel,fs_type=no_value,global=tag,host=oliveaglec841,mount=D}} true
// point --->  &{win.wmi.os.caption 1425524770 Microsoft Windows Server 2008 R2 Enterprise  {bu=hotel,global=tag,host=oliveaglec841}} true
// point --->  &{win.wmi.os.csdversion 1425524770 Service Pack 1 {bu=hotel,global=tag,host=oliveaglec841}} true
// point --->  &{win.wmi.service.iis.state 1425524770 IIS Not Installed {bu=hotel,global=tag,host=oliveaglec841}} true

var (
	influxdb_host_url, _ = url.Parse("http://192.168.59.103:8086/write")
	influxdb_client_conf = client.Config{
		URL:      *influxdb_host_url,
		Username: "root",
		Password: "root",
	}
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

		for j := 0; j < 3; j++ {
			fields[fmt.Sprintf("field%d", j)] = rand.Intn(100)
		}
		// now, _ := client.EpochToTime(time.Now().UnixNano(), "n")
		point := client.Point{
			Name: name,
			// Tags:      tags,
			Timestamp: time.Now(),
			Fields:    fields,
			Precision: "s",
		}
		// fmt.Println("point: ", point)
		points = append(points, point)
	}
	return
}

func main() {
	pretty.Println("")

	mockDataPoint()

	cli, _ := client.NewClient(influxdb_client_conf)
	fmt.Println(cli.Ping())

	res, err := cli.Query(client.Query{
		Command:  "select * from data1",
		Database: "metrics",
	})
	if err != nil {
		fmt.Println(err)
	}
	if res != nil {
		for _, r := range res.Results {
			for _, s := range r.Series {
				for _, v := range s.Values {
					fmt.Printf("%v\n", v)
				}
				fmt.Println("Count: ", len(s.Values))
			}
		}
	}

	// fmt.Println(res, err)
	write := mockWrite()
	// pretty.Println(write)
	fmt.Println(len(write.Points))

	res, err = cli.Write(write)
	fmt.Println(res, err)

	fmt.Println(err)
}
