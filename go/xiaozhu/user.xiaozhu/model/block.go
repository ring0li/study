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

type Block struct {
	BlockID      int `gorm:"column:blockId;primary_key" json:"blockId"`
	BlockType    int `gorm:"column:blockType" json:"blockType"`
	UserID       int `gorm:"column:userId" json:"userId"`
	FriendUserID int `gorm:"column:friendUserId" json:"friendUserId"`
}

// TableName sets the insert table name for this struct type
func (b *Block) TableName() string {
	return "block"
}
