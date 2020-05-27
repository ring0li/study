package datamodels

type UserInfo struct {
	UID int `gorm:"column:uid" json:"uid"`
}

type UserInfo1 struct {
	UID int ` json:"uid"`
}
