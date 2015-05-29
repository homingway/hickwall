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
	return ReadRuntimeConfig(file)
}

func load_group_conf(filepath string) (*CollectorConfigGroup, error) {
	var (
		ccg CollectorConfigGroup
		err error
	)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			err = fmt.Errorf("load_group_failed: path: %s, err: %v", filepath, err)
		}
	}()
	//	panic("hahah")

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	vp := viper.New()
	vp.SetConfigType("yaml")
	fmt.Println("---------- nothing wrong ----------")
	err = vp.ReadConfig(file)
	fmt.Println("---------- nothing wrong ----------")
	if err != nil {
		return nil, fmt.Errorf("load_group_failed: path: %s, err: %v", filepath, err)
	}

	vp.Marshal(&ccg)
	return &ccg, nil
}

func LoadRuntimeConfigFromFiles() (rc *RuntimeConfig, err error) {
	if CONF_FILEPATH != "" {
		rc, err = load_runtime_conf(CONF_FILEPATH)
		if err != nil {
			return nil, fmt.Errorf("cannot load runtime config: %v", err)
		}
	}

	fmt.Println("hahah ---------------------- 1")

	if CONF_GROUP_DIRECTORY != "" {
		files, err := ioutil.ReadDir(CONF_GROUP_DIRECTORY)
		if err == nil {
			for _, f := range files {
				filepath := path.Join(CONF_GROUP_DIRECTORY, f.Name())
				fmt.Println("filepath: ", filepath)
				if ccg, err := load_group_conf(filepath); err == nil {
					rc.Groups = append(rc.Groups, ccg)
				} else {
					fmt.Println("error: ", err)
				}
			}
		}
	}
	fmt.Println("hahah ---------------------- 2")
	return
}
