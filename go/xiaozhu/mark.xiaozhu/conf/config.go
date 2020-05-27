package conf

import (
	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"xiaozhu/utils/common"
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

type RedisChild struct {
	Host     string `yaml:"host"`
	Pwd      string `yaml:"pwd"`
	Dbname   string `yaml:"dbname"`
	MinConns int    `yaml:"minConns"`
	MaxConns int    `yaml:"maxConns"`
}

type UserServer struct {
	Host string `yaml:"host"`
}

// Config for user
type WechatConf struct {
	AppID          string `yaml:"appId"`
	AppSecret      string `yaml:"appSecret"`
	Token          string `yaml:"token"`
	EncodingAESKey string `yaml:"encodingAESKey"`
	PayMchID       string `yaml:"payMchID"`     //支付 - 商户 ID
	PayNotifyURL   string `yaml:"payNotifyURL"` //支付 - 接受微信支付结果通知的接口地址
	PayKey         string `yaml:"payKey"`       //支付 - 商户后台设置的支付 key
	Cache          string `yaml:"cache"`
}

//
type YamlConf struct {
	Port           string     `yaml:"Port"`
	Debug          bool       `yaml:"Debug"`
	DbType         string     `yaml:"DbType"`
	Mysql          MysqlChild `yaml:"Mysql"`
	glogsRootPath  string     `yaml:"LogsRootPath"`
	GRPCUserServer UserServer `yaml:"GRPCUserServer"`
	Redis          RedisChild `yaml:"Redis"`
	Wechat         WechatConf `yaml:"Wechat"`
}

//
type dataHandle struct {
	Conf        *YamlConf
	MainDb      *gorm.DB
	RedisClient *redis.Client
}

var DataHandle = &dataHandle{
	Conf: &YamlConf{},
	//MainDb :&gorm.DB{},
}

//
func init() {
	err := common.ReadConf("../app.yaml", DataHandle.Conf)
	if err != nil {
		glog.Exitf("init yaml err = %+v", err)
	}
	//
	//DataHandle.RedisClient = InitRedis(DataHandle.Conf.Redis)
	////
	//GrpcServer.conn(DataHandle.Conf.GRPCUserServer.Host)
}
