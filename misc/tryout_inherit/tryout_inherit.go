package tryout_inherit

import (
	"fmt"
	"time"
)

type Collector interface {
	Run()
	GetName() string
	SetName(name string)
}

type IntervalCollector struct {
	Interval int
	Name     string
	config   string
}

func (c *IntervalCollector) Run() {
	done := time.After(time.Second * 3)
	tick := time.Tick(time.Second * 1)
loop:
	for {
		select {
		case <-done:
			break loop
		case <-tick:
			// fmt.Println("I'm runing: ", c)
			// fmt.Println(c.GetName())
			x := (Collector)(c)
			// fmt.Println("Coverted to Collector: ", x)
			fmt.Println(x.GetName())
		}
	}
}

func (c *IntervalCollector) GetName() string {
	return c.Name
}

func (c *IntervalCollector) SetName(name string) {
	c.Name = name
	c.config = name
}

type C_win_pdh struct {
	IntervalCollector
}

func (c *C_win_pdh) GetName() string {
	// C_win_pdh 只override GetName 方法。
	//	如果从 Collector.GetName 来调用，调用的是C_win_pdh.GetName, 因为该Collecotr对应的实现类是C_Win_pdh
	//  如果从 IntervalCollector.Run中来调用 GetName, 则直接调用的是IntervalCollector.GetName

	return "you will not see this if called from parent"
	// return fmt.Sprintf("this is from C_win_pdh.GetName: %s", c.Name)
}
