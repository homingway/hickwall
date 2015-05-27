package hickwall

import (
	"fmt"
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/newcore"
)

func create_running_core_hooked(rconf *config.RuntimeConfig, ishook bool) (newcore.PublicationSet, *newcore.HookBackend, error) {
	var hook *newcore.HookBackend
	var subs []newcore.Subscription
	var heartbeat_exists bool

	if rconf == nil {
		return nil, nil, fmt.Errorf("RuntimeConfig is nil")
	}

	bks, err := backends.UseConfigCreateBackends(rconf)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("bks: ", bks)

	clrs, err := collectors.UseConfigCreateCollectors(rconf)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("collectors: ", clrs)

	for _, c := range clrs {
		if c.Name() == "heartbeat" {
			heartbeat_exists = true
		}
		subs = append(subs, newcore.Subscribe(c, nil))
	}

	if heartbeat_exists == false {
		subs = append(subs, newcore.Subscribe(collectors.NewHeartBeat(rconf.Client.HeartBeatInterval), nil))
	}

	merge := newcore.Merge(subs...)

	if ishook == true {
		hook = newcore.NewHookBackend()
		bks = append(bks, hook)
		fset := newcore.FanOut(merge, bks...)
		return fset, hook, nil
	} else {
		fset := newcore.FanOut(merge, bks...)
		return fset, nil, nil
	}
}

func CreateRunningCore(rconf *config.RuntimeConfig) (newcore.PublicationSet, error) {
	core, _, err := create_running_core_hooked(rconf, false)
	if err != nil {
		return nil, err
	}
	return core, nil
}
