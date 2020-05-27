package conf

import (
	"xiaozhu/utils/common"
	//
	"github.com/jinzhu/gorm"
)

type Deploy struct {
	RemoteUrl string `yaml:"RemoteUrl"`
	DeployDir string `yaml:"DeployDir"`
}

type HostConf struct {
	Host        string `yaml:"Host"`
	CreateShell string `yaml:"CreateShell"`
	DeployShell string `yaml:"DeployShell"`
}

type MysqlChild struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Pwd      string `yaml:"pwd"`
	Dbname   string `yaml:"dbname"`
	MinConns int    `yaml:"minConns"`
	MaxConns int    `yaml:"maxConns"`
}

//
type YamlConf struct {
	Port         string     `yaml:"Port"`
	Debug        bool       `yaml:"Debug"`
	DbType       string     `yaml:"DbType"`
	Mysql        MysqlChild `yaml:"Mysql"`
	LogsRootPath string     `yaml:"LogsRootPath"`
}

//
type dataHandle struct {
	Conf   *YamlConf
	MainDb *gorm.DB
}

var DataHandle = &dataHandle{
	Conf: &YamlConf{},
	//MainDb :&gorm.DB{},
}

//
func InitConf(path string) {
	//
	//err := libs.ReadConf(path,  DataHandle.Conf)
	common.ReadConf(path, DataHandle.Conf)
	//
	DataHandle.MainDb = InitDb(DataHandle.Conf.Mysql)
}
