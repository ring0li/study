package services

import (
	"xiaozhu/mark.xiaozhu/conf"
	"xiaozhu/user.xiaozhu/model"
)

type userService struct{}

var UserService userService

func (self *userService) GetList() []*model.User {

	userAll := []*model.User{}
	conf.DataHandle.MainDb.Find(&userAll)

	return userAll
}
