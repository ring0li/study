package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

type User struct {
	UserID             int       `gorm:"column:userId;primary_key" json:"userId"`
	OpenID             string    `gorm:"column:openId" json:"openId"`
	UnionID            string    `gorm:"column:unionId" json:"unionId"`
	Nickname           string    `gorm:"column:nickname" json:"nickname"`
	Sex                int       `gorm:"column:sex" json:"sex"`
	City               string    `gorm:"column:city" json:"city"`
	Province           string    `gorm:"column:province" json:"province"`
	Country            string    `gorm:"column:country" json:"country"`
	HeadImgURL         string    `gorm:"column:headImgUrl" json:"headImgUrl"`
	UserStatus         int       `gorm:"column:userStatus" json:"userStatus"`
	UserLastActiveTime null.Time `gorm:"column:userLastActiveTime" json:"userLastActiveTime"`
	UserCreateTime     null.Time `gorm:"column:userCreateTime" json:"userCreateTime"`
	UserUpdateTime     null.Time `gorm:"column:userUpdateTime" json:"userUpdateTime"`
}

// TableName sets the insert table name for this struct type
func (u *User) TableName() string {
	return "user"
}
