package backends

import (
	"encoding/json"
	"fmt"
	"github.com/oliveagle/boltq"
	"github.com/oliveagle/go-collectors/datapoint"
)

func MdPush(q *boltq.BoltQ, md datapoint.MultiDataPoint) error {
	if len(md) == 0 {
		return nil
	}

	// TODO: 内存有溢出，json.Marshal(md)
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

func MdPop(q *boltq.BoltQ) (md datapoint.MultiDataPoint, err error) {
	dump, err := q.Pop()
	if err != nil {
		// fmt.Println("DequeueBatch Error: ", err)
		return nil, err
	}

	err = json.Unmarshal(dump, &md)
	return
}

func MdPopBottom(q *boltq.BoltQ) (md datapoint.MultiDataPoint, err error) {
	dump, err := q.PopBottom()
	if err != nil {
		// fmt.Println("DequeueBatch Error: ", err)
		return nil, err
	}

	err = json.Unmarshal(dump, &md)
	return
}

func MdPopMany(q *boltq.BoltQ, points_cnt_limit int) (md datapoint.MultiDataPoint, err error) {

	cnt := 0
	err = q.PopMany(func(v []byte) bool {
		tmp_md := datapoint.MultiDataPoint{}
		err = json.Unmarshal(v, &tmp_md)
		if err != nil {
			return false
		}
		fmt.Printf("cnt: %d , len(tmp_md): %d\n", cnt, len(tmp_md))
		if cnt+len(tmp_md) <= points_cnt_limit {
			for _, p := range tmp_md {
				md = append(md, p)
			}
			cnt += len(tmp_md)
		} else {
			return false
		}
		return true
	})
	if err != nil {
		// fmt.Println("DequeueBatch Error: ", err)
		return nil, err
	}
	return
}

func MdPopManyBottom(q *boltq.BoltQ, points_cnt_limit int) (md datapoint.MultiDataPoint, err error) {

	cnt := 0
	err = q.PopManyBottom(func(v []byte) bool {
		tmp_md := datapoint.MultiDataPoint{}
		err = json.Unmarshal(v, &tmp_md)
		if err != nil {
			return false
		}

		fmt.Printf("cnt: %d , len(tmp_md): %d\n", cnt, len(tmp_md))
		if cnt+len(tmp_md) <= points_cnt_limit {
			for _, p := range tmp_md {
				md = append(md, p)
			}
			cnt += len(tmp_md)
		} else {
			return false
		}
		return true
	})
	if err != nil {
		// fmt.Println("DequeueBatch Error: ", err)
		return nil, err
	}
	return
}
