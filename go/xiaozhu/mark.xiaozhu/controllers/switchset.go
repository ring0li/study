package controllers

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"strconv"
	"xiaozhu/mark.xiaozhu/conf"
	pb "xiaozhu/protos"
	"xiaozhu/utils/common"
)

type switchset struct{}

var Switchset *switchset

type ele struct {
	FriendOpenId string `json:"friendOpenId"`
	FriendName   string `json:"friendName"`
	HeadImgUrl   string `json:headImgUrl`
	Switchval    bool   `json:"switchval"`
}

func (this *switchset) GetMsgReceiveSwitch(c echo.Context) error {

	openId := c.FormValue("openId")

	req := &pb.GetUserByOpenidRequest{OpenId: openId}
	userInfo, err := conf.GrpcServer.UserClient.GetUserByOpenid(context.Background(), req)

	//
	var data = map[string]bool{}
	if err == nil {
		var _switchval bool = true
		if userInfo.IsBlock == 1 {
			_switchval = false
		}
		//
		data = map[string]bool{
			"switchval": _switchval,
		}
	}
	o := common.Succ(data)

	return common.Output(c, o)
}

func (this *switchset) SetMsgReceiveSwitch(c echo.Context) error {

	openId := c.FormValue("openId")
	switchval := c.FormValue("switchval")
	//
	var _switchval int32 = 0
	if switchval == "true" {
		_switchval = 1
	}
	//
	req := &pb.AllRequest{OpenId: openId, Switchval: _switchval}
	_, err := conf.GrpcServer.BlockClient.All(context.Background(), req)

	//
	o := common.Succ("succ")
	if err != nil {
		o = common.Fail(101, err.Error())
	}
	return common.Output(c, o)
}

func (this *switchset) GetFriendList(c echo.Context) error {

	openId := c.FormValue("openId")
	//
	req := &pb.GetMyFriendsRequest{OpenId: openId}
	result, err := conf.GrpcServer.BlockClient.GetMyFriends(context.Background(), req)

	//
	// o := common.Fail(101, err.Error())
	var data []ele
	//
	o := common.Succ(data)
	if err == nil {
		for _, _ele := range result.List {
			_switchval := true
			if _ele.IsBlock == 0 {
				_switchval = false
			}

			data = append(data, ele{
				FriendOpenId: _ele.OpenId,
				FriendName:   _ele.Nickname,
				HeadImgUrl:   _ele.HeadImgUrl,
				Switchval:    _switchval,
			})
		}
		o = common.Succ(data)
	}
	return common.Output(c, o)
}

func (this *switchset) SetFriendMsgReceiveSwitch(c echo.Context) error {

	openId := c.FormValue("openId")
	friendOpenId := c.FormValue("friendOpenId")
	switchval := c.FormValue("switchval")
	//
	var _switchval int32 = 0
	if switchval == "true" {
		_switchval = 1
	}

	//
	req := &pb.FriendRequest{OpenId: openId, FriendOpenId: friendOpenId, Switchval: _switchval}
	_, err := conf.GrpcServer.BlockClient.Friend(context.Background(), req)
	//
	o := common.Succ("succ")
	if err != nil {
		o = common.Fail(101, err.Error())
	}
	return common.Output(c, o)
}

func (this *switchset) DelFriend(c echo.Context) error {

	openId := c.FormValue("openId")
	friendOpenId := c.FormValue("friendOpenId")
	fmt.Println(openId, friendOpenId)

	o := common.Succ("")
	return common.Output(c, o)
}

func (this *switchset) GetWaterSwitch(c echo.Context) error {

	openId := c.FormValue("openId")
	//
	req := &pb.GetRequest{OpenId: openId}
	rst, err := conf.GrpcServer.WaterClient.Get(context.Background(), req)
	//
	var data = map[string]int32{}
	if err == nil {
		data = map[string]int32{
			"Waterval": rst.Waterval,
		}
	}
	o := common.Succ(data)
	//
	return common.Output(c, o)
}

func (this *switchset) SetWaterSwitch(c echo.Context) error {

	openId := c.FormValue("openId")
	waterval := c.FormValue("waterval")
	_waterval, _ := strconv.Atoi(waterval)
	//
	req := &pb.SetRequest{OpenId: openId, Waterval: int32(_waterval)}
	_, err := conf.GrpcServer.WaterClient.Set(context.Background(), req)
	//
	o := common.Succ("succ")
	if err != nil {
		o = common.Fail(101, err.Error())
	}
	return common.Output(c, o)
}
