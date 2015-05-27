package config

import (
	//	"bytes"
	//	"fmt"
	//	"github.com/oliveagle/hickwall/utils"
	"fmt"
	//	log "github.com/oliveagle/seelog"
	"github.com/oliveagle/viper"
	"io/ioutil"
	"os"
	"path"
)

//func loadRuntimeConfFromFile() <-chan *RespConfig {
//	log.Debug("loadRuntimeConfFromFile")
//
//	var (
//		out           = make(chan *RespConfig, 1)
//		runtime_viper = viper.New()
//	)
//	// runtime_viper.SetConfigFile(config_file)
//	runtime_viper.SetConfigName("config")
//	runtime_viper.SetConfigType("yaml")
//	runtime_viper.AddConfigPath(SHARED_DIR) // packaged distribution
//	runtime_viper.AddConfigPath("../..")    // for hickwall/misc/try_xxx
//	runtime_viper.AddConfigPath(".")        // for hickwall
//	runtime_viper.AddConfigPath("..")       // for hickwall/misc
//
//	go func() {
//		var runtime_conf RuntimeConfig
//
//		err := runtime_viper.ReadInConfig()
//
//		log.Debug("RuntimeConfig File Used: ", runtime_viper.ConfigFileUsed())
//
//		// fmt.Println("RuntimeConfig File Used: ", runtime_viper.ConfigFileUsed())
//
//		if err != nil {
//			log.Error("loadRuntimeConfFromFile error: ", err)
//			out <- &RespConfig{nil, fmt.Errorf("No configuration file loaded. config.yml: %v", err)}
//			return
//		}
//
//		// Marshal values
//		err = runtime_viper.Marshal(&runtime_conf)
//		if err != nil {
//			log.Error("loadRuntimeConfFromFile error: ", err)
//			out <- &RespConfig{nil, fmt.Errorf("Error: unable to parse Configuration: %v\n", err)}
//			return
//		}
//
//		out <- &RespConfig{&runtime_conf, nil}
//		close(out)
//		return
//	}()
//
//	return out
//}
//
//func LoadRuntimeConfFromFileOnce() error {
//	defer log.Flush()
//
//	for resp := range loadRuntimeConfFromFile() {
//		if resp.Err != nil {
//			log.Errorf("cannot load runtime config from file: %v", resp.Err)
//			return fmt.Errorf("cannot load runtime config from file: %v", resp.Err)
//		} else {
//			UpdateRuntimeConf(resp.Config)
//			log.Debug("updated runtime config")
//		}
//	}
//	return nil
//}

func load_runtime_conf(filepath string) (*RuntimeConfig, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rc RuntimeConfig

	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.ReadConfig(file)

	vp.Marshal(&rc)
	return &rc, nil
}

func load_group_conf(filepath string) (*CollectorConfigGroup, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ccg CollectorConfigGroup

	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.ReadConfig(file)

	vp.Marshal(&ccg)
	return &ccg, nil
}

func load_runtime_conf_from_files() (rc *RuntimeConfig, err error) {
	if CONF_FILEPATH != "" {
		rc, err = load_runtime_conf(CONF_FILEPATH)
		if err != nil {
			return nil, fmt.Errorf("cannot load runtime config: %v", err)
		}
	}

	if CONF_GROUP_DIRECTORY != "" {
		files, err := ioutil.ReadDir(CONF_GROUP_DIRECTORY)
		if err == nil {
			for _, f := range files {
				filepath := path.Join(CONF_GROUP_DIRECTORY, f.Name())
				if ccg, err := load_group_conf(filepath); err == nil {
					rc.Groups = append(rc.Groups, ccg)
				}
			}
		}
	}
	return
}
