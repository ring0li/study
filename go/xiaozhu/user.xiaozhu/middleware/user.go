package middleware

import (
	"xiaozhu/mark.xiaozhu/conf"
	"xiaozhu/user.xiaozhu/datamodels"
	"xiaozhu/user.xiaozhu/model"
)

func GetUserByOpenId(openId string) (user model.User) {
	conf.DataHandle.MainDb.Where("openId = ?", openId).First(&user)
	return
}

func GetFriendsByUserId(userId int) ([]*datamodels.BlockFriend) {
	//下线的好友，多个人
	blockFriends := []datamodels.BlockFriend{}
	conf.DataHandle.MainDb.Raw("SELECT f.userId,friendUserId,openId,nickname,sex,headImgUrl FROM `friend` f "+
		" inner join `user` u on f.friendUserId = u.userId"+
		" where f.userId = ?", userId).Find(&blockFriends)

	list := []*datamodels.BlockFriend{}
	for _, v := range blockFriends {
		tmp := datamodels.BlockFriend{}
		tmp.UserID = int(v.FriendUserId)
		tmp.OpenID = v.OpenID
		tmp.Nickname = v.Nickname
		tmp.Sex = int(v.Sex)
		tmp.HeadImgUrl = v.HeadImgUrl
		tmp.IsBlock = MyFriendIsBlock(userId, v.FriendUserId)
		list = append(list, &tmp)
	}

	//上线的好友，1个人
	blockFriend := datamodels.BlockFriend{}
	conf.DataHandle.MainDb.Raw("SELECT f.userId,friendUserId,openId,nickname,sex,headImgUrl FROM `friend` f "+
		" inner join `user` u on f.userId = u.userId"+
		" where f.friendUserId = ?", userId).Find(&blockFriend)

	if blockFriend != (datamodels.BlockFriend{}) {
		blockFriend.IsBlock = MyFriendIsBlock(blockFriend.FriendUserId, blockFriend.UserID)
		list = append(list, &blockFriend)
	}

	return list
}
