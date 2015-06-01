package hickwall

import (
	"fmt"
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/newcore"
	log "github.com/oliveagle/seelog"
)

var (
	core *newcore.PublicationSet
)

func create_running_core_hooked(rconf *config.RuntimeConfig, ishook bool) (newcore.PublicationSet, *newcore.HookBackend, error) {
	var hook *newcore.HookBackend
	var subs []newcore.Subscription
	var heartbeat_exists bool

	if rconf == nil {
		log.Error("RuntimeConfig is nil")
		return nil, nil, fmt.Errorf("RuntimeConfig is nil")
	}

	bks, err := backends.UseConfigCreateBackends(rconf)
	if err != nil {
		log.Error("UseConfigCreateBackends failed: ", err)
		return nil, nil, err
	}
	fmt.Println("bks: ", bks)

	clrs, err := collectors.UseConfigCreateCollectors(rconf)
	if err != nil {
		log.Error("UseConfigCreateCollectors failed: ", err)
		return nil, nil, err
	}

	for _, c := range clrs {
		if c.Name() == "heartbeat" {
			heartbeat_exists = true
		}
		//		subs = append(subs, newcore.Subscribe(c, nil))
	}

	if heartbeat_exists == false {
		log.Debugf("heartbeat_exists == false: len(subs): %d", len(subs))
		clrs = append(clrs, collectors.NewHeartBeat(rconf.Client.HeartBeatInterval))
	}

	fmt.Println("collectors: ", clrs)

	for _, c := range clrs {
		subs = append(subs, newcore.Subscribe(c, nil))
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
	log.Debug("running_core.CreateRunningCore")
	core, _, err := create_running_core_hooked(rconf, false)
	if err != nil {
		log.Error("running_core.CreateRunningCore: ", err)
		return nil, err
	}
	return core, nil
}

func update_core(c *newcore.PublicationSet) {
	if c != nil {
		if core != nil {
			close_core()
		}
	}
	core = c
}

func close_core() {
	(*core).Close()
	core = nil
}

func IsRunning() bool {
	if core != nil {
		return true
	}
	return false
}

func Start() error {
	if !config.IsCoreConfigLoaded() {
		err := config.LoadCoreConfig()
		if err != nil {
			return fmt.Errorf("Faild to load CoreConfig: ", err)
		}
		fmt.Println("CoreConfig Loaded")
	}

	if IsRunning() == true {
		return fmt.Errorf("one core is already running. stop it first!")
	}

	switch config.CoreConf.Config_Strategy {
	case config.ETCD:
		fmt.Println("use etcd strategy")
	case config.REGISTRY:
		fmt.Println("use registry strategy")
	default:
		fmt.Println("[default] use file strategy")
		core, err := LoadConfigFromFileAndRun()
		if err != nil {
			return fmt.Errorf("failed to create core from file: %v", err)
		}
		update_core(&core)
	}
	return nil
}

func Stop() error {
	close_core()
	return nil
}
