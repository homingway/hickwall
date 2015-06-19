package hickwall

import (
	"fmt"
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
)

var (
	the_core  newcore.PublicationSet
	the_rconf *config.RuntimeConfig
	done      = make(chan error)
)

func create_running_core_hooked(rconf *config.RuntimeConfig, ishook bool) (newcore.PublicationSet, *newcore.HookBackend, error) {
	var hook *newcore.HookBackend
	var subs []newcore.Subscription
	var heartbeat_exists bool

	if rconf == nil {
		logging.Error("RuntimeConfig is nil")
		return nil, nil, fmt.Errorf("RuntimeConfig is nil")
	}

	bks, err := backends.UseConfigCreateBackends(rconf)
	if err != nil {
		logging.Errorf("UseConfigCreateBackends failed: %v", err)
		return nil, nil, err
	}

	if len(bks) <= 0 {
		err := fmt.Errorf("no backends configured. program will do nothing.")
		logging.Error(err)
		return nil, nil, err
	}

	for _, bk := range bks {
		logging.Debugf("backend: %s", bk.Name())
		logging.Tracef("backend: %s -> %+v", bk.Name(), bk)
	}

	clrs, err := collectors.UseConfigCreateCollectors(rconf)
	if err != nil {
		logging.Errorf("UseConfigCreateCollectors failed: %v", err)
		return nil, nil, err
	}

	for _, c := range clrs {
		if c.Name() == "heartbeat" {
			heartbeat_exists = true
		}
		//		subs = append(subs, newcore.Subscribe(c, nil))
	}

	if heartbeat_exists == false {
		logging.Debugf(" heartbeat_exists == false: len(subs): %d", len(subs))
		clrs = append(clrs, collectors.NewHeartBeat(rconf.Client.HeartBeat_Interval))
	}

	// fmt.Printf("collectors: %+v", clrs)

	for _, c := range clrs {
		logging.Debugf("collector: %s", c.Name())
		logging.Tracef("collector: %s -> %+v", c.Name(), c)
		subs = append(subs, newcore.Subscribe(c, nil))
	}

	// create other subscriptions, such as kafka
	_subs, err := collectors.UseConfigCreateSubscription(rconf)
	if err != nil {
		logging.Error("UseConfigCreateSubscription failed: ", err)
		return nil, nil, err
	}
	subs = append(subs, _subs...)

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
	logging.Debug("running_core.CreateRunningCore")
	core, _, err := create_running_core_hooked(rconf, false)
	if err != nil {
		logging.Errorf("running_core.CreateRunningCore: %v", err)
		return nil, err
	}
	return core, nil
}

func replace_core(c newcore.PublicationSet, rconf *config.RuntimeConfig) {
	// do nothing if nil interface
	if c == nil {
		return
	}

	// close first
	if the_core != nil {
		close_core()
	}

	the_core = c
	the_rconf = rconf
}

func close_core() {
	if the_core != nil {
		the_core.Close()
	}
	the_core = nil
	the_rconf = nil
}

func IsRunning() bool {
	if the_core != nil {
		return true
	}
	return false
}

func Start() error {
	if IsRunning() == true {
		err := fmt.Errorf("one core is already running. stop it first!")
		logging.Errorf("failed to start hickwall core: %v", err)
		return err
	}

	switch config.CoreConf.Config_Strategy {
	case config.ETCD:
		logging.Info("use etcd config strategy")
		go LoadConfigStrategyEtcd(done)
	case config.REGISTRY:
		logging.Info("use registry config strategy")
	default:
		logging.Info("[default] use file config strategy")
		core, p_rconf, err := LoadConfigStrategyFile()
		if err != nil {
			logging.Errorf("faile to create running core from file: %v", err)
			return err
		}
		replace_core(core, p_rconf)
	}
	return nil
}

func Stop() error {
	if IsRunning() {
		switch config.CoreConf.Config_Strategy {
		case config.ETCD:
			logging.Trace("Stopping etcd strategy")
			done <- nil
		case config.REGISTRY:
			logging.Trace("Stopping registry strategy")
		default:
			logging.Trace("Stopping default file strategy")
			close_core()
		}
	}
	logging.Info("core stopped")
	return nil
}
