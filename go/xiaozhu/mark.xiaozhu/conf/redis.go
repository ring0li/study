package conf

import (
	"github.com/go-redis/redis"
	"github.com/golang/glog"
)

func init() {
	conf := DataHandle.Conf.Redis
	opt := redis.Options{
		Addr:     conf.Host,
		Password: conf.Pwd, // no password set
		DB:       0,
	}
	RedisClient := redis.NewClient(&opt)
	//
	_, err := RedisClient.Ping().Result()
	if err != nil {
		glog.Exitf("redis conn fail err = %+v", err)
	}
	DataHandle.RedisClient = RedisClient
}
