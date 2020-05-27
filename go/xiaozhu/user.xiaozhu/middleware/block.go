package middleware

import (
	"xiaozhu/mark.xiaozhu/conf"
	"xiaozhu/user.xiaozhu/model"
)

func IsBlockAll(userId int) ( int) {
	block := model.Block{}
	conf.DataHandle.MainDb.Where("blockType = 1 and userId = ? ", userId).First(&block)

	if block == (model.Block{}) {
		return 0
	}
	return 1
}

func IsBlockFriend(userId int, friendUserId int) ( int) {
	block := model.Block{}
	conf.DataHandle.MainDb.Where("blockType = 2 and userId = ? and friendUserId = ? ", userId, friendUserId).First(&block)

	if block == (model.Block{}) {
		return 0
	}
	return 1
}

//我的好友，是否屏蔽消息了
func MyFriendIsBlock(userId int, friendUserId int) ( int) {
	isBlock := IsBlockAll(friendUserId)
	if isBlock == 1 {
		return 1
	}

	isBlock = IsBlockFriend(friendUserId, userId)
	if isBlock == 1 {
		return 1
	}
	return 0
}
