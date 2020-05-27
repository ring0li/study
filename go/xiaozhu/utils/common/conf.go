package common

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

var TimeZone, _ = time.LoadLocation("Asia/Shanghai")

func ReadConf(confPath string, yamlStruct interface{}) error {
	//
	data, err := ioutil.ReadFile(confPath)
	//
	if PathExists(confPath) == false {
		return errors.New("conf file not Exists")
	}
	//
	err = yaml.Unmarshal(data, yamlStruct)
	//
	if err != nil {
		return errors.New("conf file not Exists")
	}

	return nil
}
