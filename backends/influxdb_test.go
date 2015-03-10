package backends

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"math/rand"
	"testing"
	"time"
)

func mockWrite(n int) client.Write {
	return client.Write{
		Database:        "metrics",
		RetentionPolicy: "p1",
		Points:          mockDataPoint(n),
	}
}

func mockDataPoint(n int) (points []client.Point) {
	if n <= 0 {
		n = 10
	}

	rand.Seed(time.Now().UTC().UnixNano())
	// tags := map[string]string{"bu": "hotel", "global": "tag", "host": "oliveaglec841"}
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("data%d", i)
		fields := map[string]interface{}{}

		for j := 0; j < 3; j++ {
			fields[fmt.Sprintf("field%d", j)] = rand.Intn(100)
		}
		// now, _ := client.EpochToTime(time.Now().UnixNano(), "n")
		point := client.Point{
			Name: name,
			// Tags:      tags,
			Timestamp: client.Timestamp(time.Now()),
			Fields:    fields,
			Precision: "s",
		}
		// fmt.Println("point: ", point)
		points = append(points, point)
	}
	return
}

func Test_rebatchWrites(t *testing.T) {

	w1 := mockWrite(3)
	w2 := mockWrite(4)
	t.Log("w1 points: ", len(w1.Points))
	t.Log("w2 points: ", len(w2.Points))

	ws, tail, err := rebatchWrites(2, w1, w2)
	t.Log(err)

	t.Log("ws count: ", len(ws))
	for _, w := range ws {
		t.Log("points: ", len(w.Points))
	}
	if len(ws) != 3 {
		t.Error("rebatch ws count is not 3")
	}
	t.Log("tail points: ", len(tail.Points))
	if len(tail.Points) != 1 {
		t.Error("tail count is not 1")
	}

	// t.Error("----")
}
