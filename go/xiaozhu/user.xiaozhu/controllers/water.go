package controllers

import (
	"fmt"
	"golang.org/x/net/context"
	"xiaozhu/mark.xiaozhu/conf"
	pb "xiaozhu/protos"
	"xiaozhu/user.xiaozhu/middleware"
	"xiaozhu/user.xiaozhu/model"
)

type WaterServer struct{}

func (self *WaterServer) GetWaterText(c context.Context, req *pb.GetWaterTextRequest) (*pb.GetWaterTextResponse, error) {

	waterList := middleware.GetWaterUsers(int(req.LastUserId))
	list := []*pb.List{}
	for _, v := range waterList {
		waterText := middleware.GetDrinkWaterText(int(v.UserId), v.Nickname)
		if waterText != "" {
			tmp := pb.List{}
			tmp.OpenId = v.OpenId
			tmp.UserId = v.UserId
			tmp.WaterText = waterText
			list = append(list, &tmp)
		}
	}

	result := &pb.GetWaterTextResponse{List: list}

	return result, nil
}

func (self *WaterServer) Set(c context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {

	user := middleware.GetUserByOpenId(req.OpenId)
	if user == (model.User{}) {
		return nil, fmt.Errorf("用户不存在")
	}

	water := middleware.GetWaterByUserId(user.UserID)
	if water == (model.Water{}) {
		water.UserID = user.UserID
		water.Option = int(req.Waterval)
		conf.DataHandle.MainDb.Create(&water)
	} else {
		water.Option = int(req.Waterval)
		conf.DataHandle.MainDb.Save(&water)
	}

	result := &pb.SetResponse{}
	return result, nil
}

func (self *WaterServer) Get(c context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	user := middleware.GetUserByOpenId(req.OpenId)
	if user == (model.User{}) {
		return nil, fmt.Errorf("用户不存在")
	}

	water := middleware.GetWaterByUserId(user.UserID)

	result := &pb.GetResponse{Waterval: int32(water.Option)}
	return result, nil
}
