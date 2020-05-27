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

type Friend struct {
	FriendID     int `gorm:"column:friendId;primary_key" json:"friendId"`
	UserID       int `gorm:"column:userId" json:"userId"`
	FriendUserID int `gorm:"column:friendUserId" json:"friendUserId"`
}

// TableName sets the insert table name for this struct type
func (f *Friend) TableName() string {
	return "friend"
}
