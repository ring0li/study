package controllers

import (
	"github.com/guregu/null"
	"github.com/syyongx/php2go"
	"golang.org/x/net/context"
	"time"
	"xiaozhu/user.xiaozhu/middleware"
	"xiaozhu/utils/common"

	pb "xiaozhu/protos"

	"xiaozhu/mark.xiaozhu/conf"
	"xiaozhu/user.xiaozhu/model"
)

type UserServer struct{}

//根据openId查询用户信息
func (self *UserServer) GetUserByOpenid(c context.Context, req *pb.GetUserByOpenidRequest) (*pb.GetUserByOpenidResponse, error) {
	user := middleware.GetUserByOpenId(req.OpenId)

	isBlock := 0
	if user != (model.User{}) {
		isBlock = middleware.IsBlockAll(user.UserID)
	}

	result := &pb.GetUserByOpenidResponse{
		UserId:     int32(user.UserID),
		OpenId:     user.OpenID,
		Nickname:   user.Nickname,
		Sex:        int32(user.Sex),
		HeadImgUrl: user.HeadImgURL,
		IsBlock:    int32(isBlock),
	}

	return result, nil
}

//根据openId查询用户信息（批量）
func (self *UserServer) GetUsersByOpenids(c context.Context, req *pb.GetUsersByOpenidsRequest) (*pb.GetUsersByOpenidsResponse, error) {
	openids := php2go.Explode(",", req.OpenIds)
	users := []*model.User{}
	conf.DataHandle.MainDb.Where("openId in (?)", openids).Find(&users)

	list := []*pb.GetUserByOpenidResponse{}
	for _, v := range users {
		tmp := pb.GetUserByOpenidResponse{}
		tmp.UserId = int32(v.UserID)
		tmp.OpenId = v.OpenID
		tmp.UnionId = v.UnionID
		tmp.Nickname = v.Nickname
		tmp.Sex = int32(v.Sex)
		tmp.City = v.City
		tmp.Province = v.Province
		tmp.Country = v.Country
		tmp.HeadImgUrl = v.HeadImgURL

		list = append(list, &tmp)
	}

	result := &pb.GetUsersByOpenidsResponse{List: list}

	return result, nil
}

//保存用户信息
func (self *UserServer) SaveUser(c context.Context, req *pb.SaveUserRequest) (*pb.SaveUserResponse, error) {
	user := model.User{}
	conf.DataHandle.MainDb.Where("openId = ?", req.OpenId).First(&user)

	if user == (model.User{}) {
		user.OpenID = req.OpenId
		user.UnionID = req.UnionId
		user.Nickname = req.Nickname
		user.Sex = int(req.Sex)
		user.City = req.City
		user.Province = req.Province
		user.Country = req.Country
		user.HeadImgURL = req.HeadImgUrl
		user.UserStatus = 0
		user.UserLastActiveTime = null.TimeFrom(time.Now())
		user.UserCreateTime = null.TimeFrom(time.Now())

		conf.DataHandle.MainDb.Create(&user)

		//初始化clock_stat表
		initClockDate, _ := time.ParseInLocation("2006-01-02 15:04:05", "2000-01-01 00:00:00", common.TimeZone)
		clockStat := model.ClockStat{}
		clockStat.UserID = user.UserID
		clockStat.ClockType = 1
		clockStat.ClockDate = null.TimeFrom(initClockDate)
		clockStat.AvgClockTime = "07:00:00"
		conf.DataHandle.MainDb.Create(&clockStat)

		clockStat1 := model.ClockStat{}
		clockStat1.UserID = user.UserID
		clockStat1.ClockType = 2
		clockStat1.ClockDate = null.TimeFrom(initClockDate)
		clockStat1.AvgClockTime = "20:30:00"
		conf.DataHandle.MainDb.Create(&clockStat1)

		//初始化8杯水
		water := model.Water{}
		water.UserID = user.UserID
		water.Option = 8
		conf.DataHandle.MainDb.Create(&water)
	} else {
		user.UnionID = req.UnionId
		user.Nickname = req.Nickname
		user.Sex = int(req.Sex)
		user.City = req.City
		user.Province = req.Province
		user.Country = req.Country
		user.HeadImgURL = req.HeadImgUrl
		user.UserStatus = 0
		user.UserLastActiveTime = null.TimeFrom(time.Now())
		user.UserUpdateTime = null.TimeFrom(time.Now())

		conf.DataHandle.MainDb.Save(&user)
	}

	result := &pb.SaveUserResponse{UserId: int32(user.UserID)}

	//如果传邀请人的openId，需要绑定好友关系
	if req.FriendOpenId != "" {
		friendUser := model.User{}
		conf.DataHandle.MainDb.Where("openId = ?", req.FriendOpenId).First(&friendUser)

		//return nil, fmt.Errorf("邀请人不存在")

		if friendUser != (model.User{}) { //邀请人存在
			friend := model.Friend{}
			conf.DataHandle.MainDb.Where("userId = ? and friendUserId = ?", friendUser.UserID, user.UserID).First(&friend)

			if friend == (model.Friend{}) {
				friend.UserID = friendUser.UserID
				friend.FriendUserID = user.UserID
				conf.DataHandle.MainDb.Create(&friend)
			}
		}
	}

	return result, nil
}

//取消关注
func (self *UserServer) Unsubscribe(c context.Context, req *pb.UnsubscribeRequest) (*pb.UnsubscribeResponse, error) {
	user := model.User{}
	conf.DataHandle.MainDb.Where("openId = ?", req.OpenId).First(&user)

	if user == (model.User{}) {
		result := &pb.UnsubscribeResponse{Success: 0}

		return result, nil
	}

	user.UserStatus = 1
	user.UserUpdateTime = null.TimeFrom(time.Now())
	conf.DataHandle.MainDb.Save(&user)

	result := &pb.UnsubscribeResponse{Success: 1}

	return result, nil
}
