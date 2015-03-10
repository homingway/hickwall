package backends

import (
	"encoding/json"
	// "fmt"
	"github.com/oliveagle/boltq"
	"github.com/oliveagle/go-collectors/datapoint"
)

func PushMd(q *boltq.BoltQ, md datapoint.MultiDataPoint) error {
	dump, err := json.Marshal(md)
	if err != nil {
		// fmt.Println("EnqueueBatch Error: ", err)
		return err
	}
	err = q.Push(dump)
	if err != nil {
		// fmt.Println("EnqueueBatch Error: ", err)
		return err
	}
	return err
}

func PopMd(q *boltq.BoltQ) (md datapoint.MultiDataPoint, err error) {
	dump, err := q.Pop()
	if err != nil {
		// fmt.Println("DequeueBatch Error: ", err)
		return nil, err
	}

	err = json.Unmarshal(dump, &md)
	return
}

func PopBottomMd(q *boltq.BoltQ) (md datapoint.MultiDataPoint, err error) {
	dump, err := q.PopBottom()
	if err != nil {
		// fmt.Println("DequeueBatch Error: ", err)
		return nil, err
	}

	err = json.Unmarshal(dump, &md)
	return
}
