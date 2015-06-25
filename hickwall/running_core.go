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
		return nil, nil, fmt.Errorf("RuntimeConfig is nil")
	}

	bks, err := backends.UseConfigCreateBackends(rconf)
	if err != nil {
		return nil, nil, err
	}

	if len(bks) <= 0 {
		return nil, nil, fmt.Errorf("no backends configured. program will do nothing.")
	}

	for _, bk := range bks {
		logging.Debugf("backend: %s", bk.Name())
		logging.Tracef("backend: %s -> %+v", bk.Name(), bk)
	}

	clrs, err := collectors.UseConfigCreateCollectors(rconf)
	if err != nil {
		return nil, nil, err
	}

	for _, c := range clrs {
		if c.Name() == "heartbeat" {
			heartbeat_exists = true
		}
	}

	if heartbeat_exists == false {
		logging.Debugf(" heartbeat_exists == false: len(subs): %d", len(subs))
		clrs = append(clrs, collectors.NewHeartBeat(rconf.Client.HeartBeat_Interval))
	}

	for _, c := range clrs {
		logging.Debugf("collector: %s", c.Name())
		logging.Tracef("collector: %s -> %+v", c.Name(), c)
		subs = append(subs, newcore.Subscribe(c, nil))
	}

	// create other subscriptions, such as kafka
	_subs, err := collectors.UseConfigCreateSubscription(rconf)
	if err != nil {
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

// Update RunningCore with provided RuntimeConfig.
func UpdateRunningCore(rconf *config.RuntimeConfig) (newcore.PublicationSet, error) {
	logging.Debug("CreateRunningCore")
	core, _, err := create_running_core_hooked(rconf, false)
	if err != nil {
		return nil, err
	}

	if the_core != nil {
		close_core()
	}
	the_core = core
	the_rconf = rconf

	return core, nil
}

//func replace_core(c newcore.PublicationSet, rconf *config.RuntimeConfig) {
//	// do nothing if nil interface
//	if c == nil {
//		return
//	}
//
//	// close first
//	if the_core != nil {
//		close_core()
//	}
//
//	the_core = c
//	the_rconf = rconf
//}

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
		return fmt.Errorf("one core is already running. stop it first!")
	}
	logging.Info("Starting the core.")

	switch config.ValidStrategy(config.CoreConf.ConfigStrategy) {
	case config.ETCD:
		logging.Info("use etcd config strategy")
		if len(config.CoreConf.EtcdMachines) <= 0 {
			logging.Critical("EtcdMachines is empty!!")
			return fmt.Errorf("EtcdMachines is empty!!")
		}
		if config.CoreConf.EtcdPath == "" {
			logging.Critical("EtcdPath is empty!!")
			return fmt.Errorf("EtcdPath is empty!!")
		}
		go NewCoreFromEtcd(config.CoreConf.EtcdMachines, config.CoreConf.EtcdPath, done)
	case config.REGISTRY:
		logging.Info("use registry config strategy")
		if len(config.CoreConf.RegistryURLs) <= 0 {
			logging.Criticalf("RegistryURLs is empty!!")
			return fmt.Errorf("RegistryURLS is empty!!")
		}
		go RegistryAndRun(done)
	default:
		logging.Info("[default] use file config strategy")
		_, _, err := NewCoreFromFile()
		if err != nil {
			// logging.Errorf("faile to create running core from file: %v", err)
			logging.Error(err)
			return err
		}
	}
	return nil
}

func Stop() error {
	if IsRunning() {
		switch config.CoreConf.ConfigStrategy {
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
