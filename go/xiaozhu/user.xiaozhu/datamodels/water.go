package datamodels

type WaterTextList struct {
	Nickname string `gorm:"column:nickname" json:"nickname"`
	OpenId   string `gorm:"column:openId" json:"openId"`
	Text     string `gorm:"column:text" json:"text"`
	Option   int32  `gorm:"column:option" json:"option"`
	UserId   int32  `gorm:"column:userId" json:"userId"`
}
