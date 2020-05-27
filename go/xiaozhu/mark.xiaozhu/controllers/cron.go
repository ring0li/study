package controllers

import (
	"github.com/golang/glog"
	"github.com/robfig/cron"
	"strconv"
	"xiaozhu/mark.xiaozhu/services"
)

type cronWork struct{}

var CronWork *cronWork

func (this *cronWork) Run() error {
	glog.Info("Crontab Run")

	c := cron.New()
	//specPunchCard := "0 */5 * * * *"
	remindStepTime := 2 //通知接口的时间段
	specPunchCard := "0 */" + strconv.Itoa(remindStepTime) + " * * * *"
	c.AddFunc(specPunchCard, func() {
		glog.Info("CronPunchCardRemind Run")
		services.WechatService.CronPunchCardRemind(remindStepTime)
		glog.Info("CronPunchCardRemind End")
	})

	specPunchCardMsgToFriend := "0 */15 * * * *"
	c.AddFunc(specPunchCardMsgToFriend, func() {
		glog.Info("CronPushMsgToFriend Run")
		services.WechatService.CronPushMsgToFriend()
		glog.Info("CronPushMsgToFriend End")
	})

	specPushWaterRemind := "0 0 7,8,9,10,11,12,13,14,15,16,17,18,19,20,21 * * *"
	c.AddFunc(specPushWaterRemind, func() {
		glog.Info("specPushWaterRemind Run")
		services.WechatService.CronPushWaterRemind()
		glog.Info("specPushWaterRemind End")
	})
	c.Start()
	select {}
}
