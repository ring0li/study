package cron

import (
	"github.com/golang/glog"
	"sync"
	"time"
	"xiaozhu/mark.xiaozhu/conf"
	"xiaozhu/user.xiaozhu/model"
)

type user struct{}

//
func (this *user) UpdateUserinfo() {
	//
	wechatUser := conf.XzWechat.GetUser()
	var waitGroup sync.WaitGroup
	var userList []*model.User
	now := time.Now()
	beforeTime := now.Add(time.Hour * time.Duration(48) * -1)
	//
	conf.DataHandle.MainDb.Where("userLastActiveTime > ? && userUpdateTime < ?", beforeTime, now).Find(&userList).Limit(10)
	if len(userList) > 0 {
		for _, user := range userList {
			waitGroup.Add(1)
			go func(row *model.User) {
				defer waitGroup.Done()

				wechatUserInfo, err := wechatUser.GetUserInfo(row.OpenID)
				if err != nil {
					glog.Errorf("update wechat info err = %+v", err)
					return
				}

				//
				glog.Infof("%+v", row, wechatUserInfo.Headimgurl)
				//
			}(user)
			waitGroup.Wait()
		}
	}
	//
	// user, err := user.GetUserInfo()
}
