package hickwall

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kr/pretty"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"github.com/oliveagle/hickwall/utils"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	_               = pretty.Sprintf("")
	_               = fmt.Sprint("")
	api_srv_running = false
)

var unsigner utils.Unsigner

func load_unsigner() error {
	if unsigner == nil && (config.CoreConf.SecureAPIWrite || config.CoreConf.SecureAPIRead) {
		s, err := utils.LoadPublicKeyFromPath(config.CoreConf.ServerPubKeyPath)
		if err != nil {
			return err
		}
		unsigner = s
		return nil
	} else {
		return fmt.Errorf("unsigner already exists")
	}
}

func reload_unsigner() error {
	if unsigner != nil {
		unsigner = nil
	}
	return load_unsigner()
}

func protect_read(h httprouter.Handle, expire time.Duration) httprouter.Handle {
	return protect(h, expire, "read")
}

func protect_write(h httprouter.Handle, expire time.Duration) httprouter.Handle {
	return protect(h, expire, "write")
}

func protect(h httprouter.Handle, expire time.Duration, trigger string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		secure := false
		if trigger == "read" && config.CoreConf.SecureAPIRead {
			secure = true
		} else if trigger == "write" && config.CoreConf.SecureAPIWrite {
			secure = true
		}
		logging.Infof("trigger: %s, secure: %v, write: %v, read: %v\n", trigger, secure, config.CoreConf.SecureAPIWrite, config.CoreConf.SecureAPIRead)

		if secure {
			hostname := r.URL.Query().Get("hostname")
			if strings.ToLower(hostname) != newcore.GetHostname() {
				logging.Errorf("hostname mismatch: %v", hostname)
				http.Error(w, "hostname mismatch", 500)
				return
			}

			time_str := r.URL.Query().Get("time")
			tm, err := utils.UTCTimeFromUnixStr(time_str)
			if err != nil {
				logging.Errorf("invalid time: %v", time_str)
				http.Error(w, "Invalid Time", 500)
				return
			}

			if time.Now().Sub(tm) > expire {
				// expired reqeust
				logging.Errorf("expired request: %v", time.Now().Sub(tm))
				http.Error(w, "expired request", 500)
				return
			}

			// we need to verify request.
			// request should put signature of this agent hostname into header HICKWALL_ADMIN_SIGN
			load_unsigner()

			signed_str := r.Header.Get("HICKWALL_ADMIN_SIGN")
			signed, err := base64.StdEncoding.DecodeString(signed_str)
			if err != nil {
				logging.Error("cannot decode sign")
				http.Error(w, "cannot decode sign", 500)
				return
			}

			toSign := fmt.Sprintf("%s%s", hostname, time_str)
			logging.Trace("unsign started")
			err = unsigner.Unsign([]byte(toSign), signed)
			logging.Trace("unsign finished")
			if err != nil {
				logging.Errorf("-> invalid signature: %v <-", string(signed))
				http.Error(w, "invalid signature", 500)
				return
			}
		}

		h(w, r, ps)
	}
}

func serveSysInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// pretty.Println(r)
	logging.Debugf("api /sys_info called from %s, query: %s", r.RemoteAddr, r.URL.RawQuery)

	sys_info, err := GetSystemInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get sys info: %v", err), 500)
		return
	}

	dump, err := json.Marshal(sys_info)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot marshal response: %v", err), 500)
		return
	}
	fmt.Fprint(w, string(dump))
}

//func serveRegistryGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	logging.Debugf("api /registry called from %s", r.RemoteAddr)
//
//	fmt.Println("ps: %v", ps)
//	fmt.Fprint(w, "Welcome!, \n")
//}
//
//func serveRegistryAccept(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	logging.Debugf("api /registry/accept called from %s", r.RemoteAddr)
//
//	fmt.Println("ps: %v", ps)
//	fmt.Fprint(w, "Welcome!, \n")
//}

func serveRegistryRevoke(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logging.Debugf("api /registry/revoke called from %s", r.RemoteAddr)

	err := Stop() // stop hickwall first
	if err != nil {
		http.Error(w, "Failed to Stop agent", 500)
		return
	}
	// delete registration file
	err = os.Remove(config.REGISTRY_FILEPATH)
	if err != nil {
		http.Error(w, "Failed to Delete Registration File", 500)
		return
	}
	err = Start() // restart hickwall. if we can pass registration process.
	if err != nil {
		http.Error(w, "Failed to Start agent", 500)
		return
	}
	logging.Info("agent started again.")
	return
}

//func serveRegistryRenew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	logging.Debugf("api /registry/renew called from %s", r.RemoteAddr)
//
//	fmt.Println("ps: %v", ps)
//	fmt.Fprint(w, "Welcome!, \n")
//}
//
//func serveConfigRenew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	logging.Debugf("api /config/renew called from %s", r.RemoteAddr)
//
//	fmt.Println("ps: %v", ps)
//	fmt.Fprint(w, "Welcome!, \n")
//}

func serve_api() {
	router := httprouter.New()
	router.GET("/sys_info", protect_read(serveSysInfo, time.Second))

	router.DELETE("/registry/revoke", protect_write(serveRegistryRevoke, time.Second))

	//	router.GET("/registry", serveRegistryGet)
	//	router.POST("/registry/accept", serveRegistryAccept)

	//	router.PUT("/registry/renew", serveRegistryRenew)
	//	router.PUT("/config/renew", serveConfigRenew)

	addr := ":3031"
	if config.CoreConf.ListenPort > 0 {
		addr = fmt.Sprintf(":%d", config.CoreConf.ListenPort)
	}
	logging.Infof("api served at: %s", addr)

	api_srv_running = true
	err := http.ListenAndServe(addr, router)
	api_srv_running = false
	if err != nil {
		logging.Criticalf("api server is not running!: %v", err)
	} else {
		logging.Info("api server stopped")
	}
}
