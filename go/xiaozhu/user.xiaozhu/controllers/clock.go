package controllers

import (
	"fmt"
	"github.com/guregu/null"
	"golang.org/x/net/context"
	"strings"
	"time"
	"unicode/utf8"
	"xiaozhu/mark.xiaozhu/conf"
	pb "xiaozhu/protos"
	"xiaozhu/user.xiaozhu/middleware"
	"xiaozhu/user.xiaozhu/model"
	"xiaozhu/utils/common"
)

type ClockServer struct{}

//早起卡
func (self *ClockServer) GetUp(c context.Context, req *pb.GetUpRequest) (*pb.GetUpResponse, error) {
	user := model.User{}
	conf.DataHandle.MainDb.Where("openId = ?", req.OpenId).First(&user)

	if user == (model.User{}) {
		return nil, fmt.Errorf("用户不存在")
	}

	//打卡
	clockTime, _ := time.ParseInLocation("2006-01-02 15:04:05", req.ClockTime, common.TimeZone)
	clock := model.ClockLog{UserID: user.UserID, ClockType: 1, ClockTime: null.TimeFrom(clockTime), ClockDate: null.TimeFrom(middleware.GetClockDate(clockTime))}
	conf.DataHandle.MainDb.Create(&clock)

	//更新最后活跃时间
	user.UserLastActiveTime = null.TimeFrom(time.Now())
	conf.DataHandle.MainDb.Save(&user)

	//设置连续打卡天数
	continueClockDays, lastClockTime := middleware.SetContinueClockDays(user.UserID, 1, clockTime)
	total, num := middleware.GetUserStat()
	imgPath, text := middleware.GetUpImg(user.HeadImgURL, continueClockDays, lastClockTime.Format("15:04"),
		total, num, clockTime.Format("02"), clockTime.Format("2006.01"), req.QrcodeUrl)

	//喝水信息
	waterText := middleware.GetWaterGetUpText(user.UserID, user.Nickname)
	result := &pb.GetUpResponse{ImgUrl: imgPath, Text: text, ContinueClockDays: int32(continueClockDays), WaterText: waterText}

	return result, nil
}

//早睡卡
func (self *ClockServer) Sleep(c context.Context, req *pb.SleepRequest) (*pb.SleepResponse, error) {
	//conf.InitConf("../app.yaml")

	user := model.User{}
	conf.DataHandle.MainDb.Where("openId = ?", req.OpenId).First(&user)

	if user == (model.User{}) {
		return nil, fmt.Errorf("用户不存在")
	}

	//打卡
	clockTime, _ := time.ParseInLocation("2006-01-02 15:04:05", req.ClockTime, common.TimeZone)
	clock := model.ClockLog{UserID: user.UserID, ClockType: 2, ClockTime: null.TimeFrom(clockTime), ClockDate: null.TimeFrom(middleware.GetClockDate(clockTime))}
	conf.DataHandle.MainDb.Create(&clock)

	//更新最后活跃时间
	user.UserLastActiveTime = null.TimeFrom(time.Now())
	conf.DataHandle.MainDb.Save(&user)

	//设置连续打卡天数
	continueClockDays, lastClockTime := middleware.SetContinueClockDays(user.UserID, 2, clockTime)
	total, num := middleware.GetUserStat()
	imgPath, text := middleware.SleepImg(user.HeadImgURL, continueClockDays, lastClockTime.Format("15:04"),
		total, num, clockTime.Format("02"), clockTime.Format("2006.01"), req.QrcodeUrl)

	//喝水信息
	waterText := middleware.GetWaterSleepText(user.UserID, user.Nickname)

	result := &pb.SleepResponse{ImgUrl: imgPath, Text: text, ContinueClockDays: int32(continueClockDays), WaterText: waterText}
	return result, nil
}

//未打卡的用户
func (self *ClockServer) NoClockUsers(c context.Context, req *pb.NoClockUsersRequest) (*pb.NoClockUsersResponse, error) {
	users := middleware.GetNoClockUsers(int(req.LastUserId), int(req.ClockType), req.StartTime, req.EndTime)
	list := []*pb.UserList{}
	for _, v := range users {
		tmp := pb.UserList{}
		tmp.UserId = int32(v.UserID)
		tmp.OpenId = v.OpenID
		tmp.Nickname = v.Nickname
		tmp.Sex = int32(v.Sex)
		tmp.IsBlock = int32(v.IsBlock)
		list = append(list, &tmp)
	}

	result := &pb.NoClockUsersResponse{List: list}
	return result, nil
}

//未打卡的好友列表
func (self *ClockServer) GetNoClockFriends(c context.Context, req *pb.GetNoClockFriendsRequest) (*pb.GetNoClockFriendsResponse, error) {
	user := middleware.GetUserByOpenId(req.OpenId)
	if user == (model.User{}) {
		return nil, fmt.Errorf("用户不存在")
	}

	friends := middleware.GetFriendsByUserId(user.UserID)

	list := []*pb.UserList{}
	for _, v := range friends {
		if v.IsBlock == 0 {
			tmp := pb.UserList{}
			tmp.UserId = int32(v.UserID)
			tmp.OpenId = v.OpenID
			tmp.Nickname = v.Nickname
			tmp.Sex = int32(v.Sex)
			tmp.IsBlock = int32(v.IsBlock)
			list = append(list, &tmp)
		}
	}

	result := &pb.GetNoClockFriendsResponse{List: list}

	return result, nil
}

//早起卡
func (self *ClockServer) Test(c context.Context, req *pb.GetUpRequest) (*pb.GetUpResponse, error) {
	user := middleware.GetUserByOpenId(req.OpenId)
	if user == (model.User{}) {
		return nil, fmt.Errorf("用户不存在")
	}

	firstText := "12312312adfasdfadddsdfasdfasdfadfadsf"
	if firstText != "" {
		firstText = strings.Repeat("　", 15-utf8.RuneCountInString(firstText)) + firstText //用Printf("%s 15s,")报错
	}
	result := &pb.GetUpResponse{ImgUrl: "", Text: ""}

	return result, nil
}

//获得打卡天数,如果当天没有打卡，要+1（类似于当天打卡了）
func (self *ClockServer) GetTotalClockDays(c context.Context, req *pb.GetTotalClockDaysRequest) (*pb.GetTotalClockDaysResponse, error) {
	user := middleware.GetUserByOpenId(req.OpenId)
	if user == (model.User{}) {
		return nil, fmt.Errorf("用户不存在")
	}

	clockStat := model.ClockStat{}
	conf.DataHandle.MainDb.Where("userId = ? and clockType = ?", user.UserID, req.ClockType).First(&clockStat)
	totalClockDays := 0
	//fmt.Println(clockStat.ClockDate.Time, time.Now())
	if clockStat == (model.ClockStat{}) {
		totalClockDays = 1
	} else if clockStat.ClockDate.Time.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		totalClockDays = clockStat.ContinueClockDays
	} else {
		totalClockDays = clockStat.ContinueClockDays + 1
	}

	result := &pb.GetTotalClockDaysResponse{TotalClockDays: int32(totalClockDays)}

	return result, nil
}
