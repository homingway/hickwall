package config

import (
	"fmt"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/viper"
	//	"io"
	"io/ioutil"
	"os"
	"path"
)

func LoadRuntimeConfFromPath(filepath string) (*RuntimeConfig, error) {
	file, err := os.Open(filepath)
	if err != nil {
		logging.Critical("failed to load runtime config from file: ", err)
		return nil, err
	}
	defer file.Close()
	return ReadRuntimeConfig(file)
}

func load_group_conf_from_filepath(filepath string) (ccg CollectorConfigGroup, err error) {
	var file *os.File

	defer func() {
		if r := recover(); r != nil {
			logging.Critical("recoverd in load_group_conf_from_filepath: ", r)
			err = fmt.Errorf("load_group_failed: path: %s, err: %v", filepath, err)
		}
	}()

	file, err = os.Open(filepath)
	if err != nil {
		logging.Errorf("failed to open file: %s, %v", filepath, err)
		return
	}
	defer file.Close()

	vp := viper.New()
	vp.SetConfigType("yaml")
	err = vp.ReadConfig(file)
	if err != nil {
		logging.Errorf("load_group_failed: path: %s, err: %v", filepath, err)
		return ccg, fmt.Errorf("load_group_failed: path: %s, err: %v", filepath, err)
	}

	vp.Marshal(&ccg)
	logging.Infof("load_group_conf_from_filepath success: %s", filepath)
	return ccg, nil
}

func LoadRuntimeConfigFromFiles() (rc *RuntimeConfig, err error) {
	if CONF_FILEPATH != "" {
		rc, err = LoadRuntimeConfFromPath(CONF_FILEPATH)
		if err != nil {
			logging.Errorf("load runtime config from files failed: %v", err)
			return rc, fmt.Errorf("cannot load runtime config: %v", err)
		}
	}

	if CONF_GROUP_DIRECTORY != "" {
		files, err := ioutil.ReadDir(CONF_GROUP_DIRECTORY)
		if err == nil {
			for _, f := range files {
				filepath := path.Join(CONF_GROUP_DIRECTORY, f.Name())
				fmt.Println("filepath: ", filepath)
				if ccg, err := load_group_conf_from_filepath(filepath); err == nil {
					rc.Groups = append(rc.Groups, ccg)
				} else {
					fmt.Println("error: ", err)
				}
			}
		}
	}
	return
}
