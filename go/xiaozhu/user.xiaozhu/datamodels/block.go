package datamodels

type BlockFriend struct {
	UserID       int    `gorm:"column:userId;primary_key" json:"userId"`
	FriendUserId int    `gorm:"column:friendUserId" json:"friendUserId"`
	OpenID       string `gorm:"column:openId" json:"openId"`
	Nickname     string `gorm:"column:nickname" json:"nickname"`
	Sex          int    `gorm:"column:sex" json:"sex"`
	IsBlock      int    `gorm:"column:isBlock" json:"isBlock"`
	HeadImgUrl   string `gorm:"column:headImgUrl" json:"headImgUrl"`
}
