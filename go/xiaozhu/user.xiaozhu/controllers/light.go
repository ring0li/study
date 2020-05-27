package controllers

import (
	"fmt"
	"golang.org/x/net/context"
	"xiaozhu/mark.xiaozhu/conf"
	pb "xiaozhu/protos"
	"xiaozhu/user.xiaozhu/middleware"
	"xiaozhu/user.xiaozhu/model"
)

type LightServer struct{}

//生成头像标签
func (self *LightServer) AddTag(c context.Context, req *pb.AddTagRequest) (*pb.AddTagResponse, error){
	
	user := model.User{}
	conf.DataHandle.MainDb.Where("openId = ?", req.OpenId).First(&user)

	if user == (model.User{}) {
		return nil, fmt.Errorf("用户不存在")
	}

	imgPath := middleware.LightTagImg(user.HeadImgURL, req.TypeName) 
	if imgPath==""{
		return nil, fmt.Errorf("用户头像下载失败")
	}
	sharePath := ""
	if req.Qrcode != "" {
		sharePath = middleware.LightShareImg(req.Qrcode) 
	}

	result := &pb.AddTagResponse{Imgurl:imgPath, Shareimgurl:sharePath}
	//
	return result, nil
}

