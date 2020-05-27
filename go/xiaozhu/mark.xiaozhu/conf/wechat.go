package conf

import (
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/cache"
)

var XzWechat *wechat.Wechat

func init() {
	//配置微信参数
	config := &wechat.Config{
		AppID:          DataHandle.Conf.Wechat.AppID,
		AppSecret:      DataHandle.Conf.Wechat.AppSecret,
		Token:          DataHandle.Conf.Wechat.Token,
		EncodingAESKey: DataHandle.Conf.Wechat.EncodingAESKey,
	}

	conf := DataHandle.Conf.Redis
	opts := &cache.RedisOpts{
		Host:        conf.Host,
		Password:    conf.Pwd,
		Database:    0,
		IdleTimeout: 2,
	}
	redisCache := cache.NewRedis(opts)
	config.Cache = redisCache
	//
	XzWechat = wechat.NewWechat(config)
}
