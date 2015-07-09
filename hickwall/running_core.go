package hickwall

import (
	"fmt"
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"net/http"
	_ "net/http/pprof"
)

var (
	the_core      newcore.PublicationSet
	the_rconf     *config.RuntimeConfig
	done          = make(chan error)
	pprof_serving = false
)

// create the topology of our running core.
//
//  (collector -> subscription)s -> merged subscription -> fanout -> publications(backends)
//
func create_running_core_hooked(rconf *config.RuntimeConfig, ishook bool) (newcore.PublicationSet, *newcore.HookBackend, error) {
	var hook *newcore.HookBackend
	var subs []newcore.Subscription
	var heartbeat_exists bool

	if rconf == nil {
		return nil, nil, fmt.Errorf("RuntimeConfig is nil")
	}

	// create backends --------------------------------------------------------------------
	bks, err := backends.UseConfigCreateBackends(rconf)
	if err != nil {
		return nil, nil, err
	}

	if len(bks) <= 0 {
		return nil, nil, fmt.Errorf("no backends configured. program will do nothing.")
	}

	for _, bk := range bks {
		logging.Debugf("loaded backend: %s", bk.Name())
		logging.Tracef("loaded backend: %s -> %+v", bk.Name(), bk)
	}

	// create collectors ------------------------------------------------------------------
	clrs, err := collectors.UseConfigCreateCollectors(rconf)
	if err != nil {
		return nil, nil, err
	}

	// make sure heartbeat collector always created.
	for _, c := range clrs {
		if c.Name() == "heartbeat" {
			heartbeat_exists = true
		}
	}
	if heartbeat_exists == false {
		clrs = append(clrs, collectors.NewHeartBeat(rconf.Client.HeartBeat_Interval))
	}

	// create collector subscriptions.
	for _, c := range clrs {
		subs = append(subs, newcore.Subscribe(c, nil))
	}

	// create other subscriptions, such as kafka consumer
	_subs, err := collectors.UseConfigCreateSubscription(rconf)
	if err != nil {
		return nil, nil, err
	}
	subs = append(subs, _subs...)

	for _, s := range subs {
		logging.Debugf("loaded subscription: %s", s.Name())
		logging.Tracef("loaded subscription: %s -> %+v", s.Name(), s)
	}

	merged := newcore.Merge(subs...)

	if ishook == true {
		// the only reason to create a hooked running core is to do unittesting. it's not a good idea though.
		hook = newcore.NewHookBackend()
		bks = append(bks, hook)
		fset := newcore.FanOut(merged, bks...)
		return fset, hook, nil
	} else {
		fset := newcore.FanOut(merged, bks...)
		return fset, nil, nil
	}
}

// Update RunningCore with provided RuntimeConfig.
func UpdateRunningCore(rconf *config.RuntimeConfig) error {
	logging.Debug("UpdateRunningCore")
	if rconf == nil {
		return fmt.Errorf("rconf is nil")
	}
	core, _, err := create_running_core_hooked(rconf, false)

	// http pprof
	// https://github.com/golang/go/issues/4674
	// we can only open http pprof, cannot close it.
	if pprof_serving == false && rconf.Client.Pprof_enabled == true {
		if rconf.Client.Pprof_listen == "" {
			rconf.Client.Pprof_listen = ":6060"
		}
		go func() {
			pprof_serving = true
			logging.Infof("http pprof is listen and served on: %v", rconf.Client.Pprof_listen)
			err := http.ListenAndServe(rconf.Client.Pprof_listen, nil)
			logging.Errorf("pprof ListenAndServe Error: %v", err)
			pprof_serving = false
		}()
	}

	// if registry give us an empty config. agent should also reflect this change.
	close_core()

	if err != nil {
		return err
	}

	//	close_core()
	the_core = core
	the_rconf = rconf
	logging.Debug("UpdateRunningCore Finished")
	return nil
}

func close_core() {
	logging.Debugf("closing the core")
	if the_core != nil {
		the_core.Close()
	}
	the_core = nil
	the_rconf = nil
	logging.Debugf("the_core now closed")
}

func Running() bool {
	if the_core != nil {
		return true
	}
	return false
}

func Start() error {

	// start api serve once.
	if !api_srv_running && config.CoreConf.Enable_http_api {
		go serve_api()
	}

	if Running() == true {
		return logging.SError("one core is already running. stop it first!")
	}
	logging.Info("Starting the core.")

	switch config.ValidStrategy(config.CoreConf.Config_strategy) {
	case config.ETCD:
		logging.Info("use etcd config strategy")
		if len(config.CoreConf.Etcd_machines) <= 0 {
			return logging.SCritical("EtcdMachines is empty!!")
		}
		if config.CoreConf.Etcd_path == "" {
			return logging.SCritical("EtcdPath is empty!!")
		}
		go new_core_from_etcd(config.CoreConf.Etcd_machines, config.CoreConf.Etcd_path, done)
	case config.REGISTRY:
		logging.Info("use registry config strategy")
		if len(config.CoreConf.Registry_urls) <= 0 {
			return logging.SCritical("RegistryURLS is empty!!")
		}
		go new_core_from_registry(done)
	default:
		logging.Info("[default] use file config strategy")
		_, err := new_core_from_file()
		if err != nil {
			return logging.SError(err)
		}
	}
	return nil
}

func Stop() error {
	if Running() {
		switch config.CoreConf.Config_strategy {
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
