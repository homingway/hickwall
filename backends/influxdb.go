package backends

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	// "github.com/oliveagle/hickwall/config"
)

type Writes []client.Write

// rebatchWrites(2, 3points, 4points) => 2points, 2points, 2points, 1point
// rebatchWrties(200, 134points, 235points) => 200points, 169points
func rebatchWrites(batchSize int, writes ...client.Write) (ws Writes, tail client.Write, err error) {
	var database = ""
	var retentionpolicy = ""
	var points = []client.Point{}

	for _, w := range writes {
		if database == "" {
			database = w.Database
		} else if database != w.Database {
			err = fmt.Errorf("cannot merge writes to different database")
			return
		}

		if retentionpolicy == "" {
			retentionpolicy = w.RetentionPolicy
		} else if retentionpolicy != w.RetentionPolicy {
			err = fmt.Errorf("cannot merge writes which have different RetentionPolicy")
			return
		}

		for _, point := range w.Points {
			points = append(points, point)
			if len(points) >= batchSize {
				ws = append(ws, client.Write{
					Database:        database,
					RetentionPolicy: retentionpolicy,
					Points:          points,
				})
				points = nil
			}
		}
	}

	if len(points) > 0 {
		tail = client.Write{
			Database:        database,
			RetentionPolicy: retentionpolicy,
			Points:          points,
		}
		points = nil
	}
	// return ws, tail, nil
	return
}
