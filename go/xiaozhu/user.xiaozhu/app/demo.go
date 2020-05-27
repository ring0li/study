package main

import (
	"fmt"
	"github.com/syyongx/php2go"
	"strconv"
	"time"
	"xiaozhu/mark.xiaozhu/conf"
	"xiaozhu/utils/common"
)

func main() {
	conf.InitConf("../app.yaml")

	onlineTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-10-25 17:38:00", common.TimeZone)
	total := (int)(time.Now().Unix()/86400 + (time.Now().Unix()-onlineTime.Unix())/60) //1分钟增加一个人

	today := time.Now().In(common.TimeZone).Format("20060102")
	fmt.Println(conf.DataHandle.RedisClient)
	clockRank, _ := strconv.Atoi(conf.DataHandle.RedisClient.Get("clock_rank_" + today).Val())
	fmt.Println(total, clockRank, today)
	if clockRank == 0 {
		clockRank = php2go.Rand(1, 1000)
	} else {
		clockRank++;
	}
	a, err := conf.DataHandle.RedisClient.Set("clock_rank_"+today, strconv.Itoa(clockRank), 86400*time.Second).Result()

	fmt.Println(total, clockRank, today, a, err)
	return

}
