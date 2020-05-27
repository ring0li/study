package controllers

import (
	"fmt"
	"golang.org/x/net/context"
	"xiaozhu/mark.xiaozhu/conf"
	pb "xiaozhu/protos"
	"xiaozhu/user.xiaozhu/middleware"
	"xiaozhu/user.xiaozhu/model"
)

type BlockServer struct{}

//屏蔽所有消息
func (self *BlockServer) All(c context.Context, req *pb.AllRequest) (*pb.AllResponse, error) {
	//conf.InitConf("../app.yaml")

	if req.Switchval != 0 && req.Switchval != 1 {
		return nil, fmt.Errorf("Switchval错误")
	}

	user := middleware.GetUserByOpenId(req.OpenId)

	if user == (model.User{}) {
		return nil, fmt.Errorf("用户不存在")
	}

	block := model.Block{}
	conf.DataHandle.MainDb.Where("blockType = 1 and userId = ?", user.UserID).First(&block)

	//
	if req.Switchval == 1 {
		if block != (model.Block{}) {
			conf.DataHandle.MainDb.Delete(&block)
		}
	} else {
		if block == (model.Block{}) {
			block.BlockType = 1
			block.UserID = user.UserID

			conf.DataHandle.MainDb.Create(&block)
		}
	}

	result := &pb.AllResponse{Code: 200, Message: "OK"}

	return result, nil
}

//屏蔽好友消息
func (self *BlockServer) Friend(c context.Context, req *pb.FriendRequest) (*pb.AllResponse, error) {
	//conf.InitConf("../app.yaml")

	user := middleware.GetUserByOpenId(req.OpenId)
	friendUser := middleware.GetUserByOpenId(req.FriendOpenId)

	if user == (model.User{}) || friendUser == (model.User{}) {
		return nil, fmt.Errorf("用户不存在")
	}

	block := model.Block{}
	conf.DataHandle.MainDb.Where("blockType = 2 and userId = ? and friendUserId = ?", user.UserID, friendUser.UserID).First(&block)

	if req.Switchval == 1 {
		if block == (model.Block{}) {
			block.BlockType = 2
			block.UserID = int(user.UserID)
			block.FriendUserID = int(friendUser.UserID)

			conf.DataHandle.MainDb.Create(&block)
		}
	} else {
		conf.DataHandle.MainDb.Delete(&block)
	}

	result := &pb.AllResponse{Code: 200, Message: "OK"}
	return result, nil
}

func (self *BlockServer) GetMyFriends(c context.Context, req *pb.GetMyFriendsRequest) (*pb.GetMyFriendsResponse, error) {
	user := middleware.GetUserByOpenId(req.OpenId)
	if user == (model.User{}) {
		return nil, fmt.Errorf("用户不存在")
	}

	users := middleware.GetFriendsByUserId(user.UserID)

	list := []*pb.GetMyFriendsResponseList{}

	for _, v := range users {
		tmp := pb.GetMyFriendsResponseList{}
		tmp.UserId = int32(v.UserID)
		tmp.OpenId = v.OpenID
		tmp.Nickname = v.Nickname
		tmp.Sex = int32(v.Sex)
		tmp.HeadImgUrl = v.HeadImgUrl
		tmp.IsBlock = int32(v.IsBlock)

		list = append(list, &tmp)
	}

	result := &pb.GetMyFriendsResponse{List: list}

	return result, nil
}
