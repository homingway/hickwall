package command

// import (
// 	// "fmt"
// 	log "github.com/cihub/seelog"
// 	"github.com/codegangsta/cli"
// 	// "github.com/oliveagle/hickwall/config"
// 	"github.com/oliveagle/hickwall/servicelib"
// 	"time"
// )

// func CmdKeepAlive(c *cli.Context) {
// 	// log.Info("CmdKeepAlive")
// 	tick := time.Tick(time.Second * time.Duration(5))
// 	for {
// 		select {
// 		case <-tick:
// 			go func() {
// 				state, err := service.Status()
// 				if err != nil {
// 					log.Errorf("CmdServiceStatus: %v", err)
// 					return
// 				}
// 				if state == servicelib.Stopped {
// 					log.Error("service is stopped! trying to start service again")

// 					err := service.StartService()
// 					if err != nil {
// 						log.Error("start service failed: ", err)
// 					} else {
// 						log.Info("service started. ")
// 					}
// 				} else {
// 					log.Info("Serivce state: ", servicelib.StateToString(state))
// 				}
// 			}()
// 		}
// 	}

// }
