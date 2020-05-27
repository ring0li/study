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

type XzUser struct {
	UID     int    `gorm:"column:uid;primary_key" json:"uid"`
	Openid  string `gorm:"column:openid" json:"openid"`
	UnionID string `gorm:"column:unionId" json:"unionId"`
	Name    string `gorm:"column:name" json:"name"`
	ImgURL  string `gorm:"column:imgUrl" json:"imgUrl"`
}

// TableName sets the insert table name for this struct type
func (x *XzUser) TableName() string {
	return "xz_user"
}
