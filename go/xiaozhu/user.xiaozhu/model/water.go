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

type Water struct {
	WaterID int `gorm:"column:waterId;primary_key" json:"waterId"`
	UserID  int `gorm:"column:userId" json:"userId"`
	Option  int `gorm:"column:option" json:"option"`
}

// TableName sets the insert table name for this struct type
func (w *Water) TableName() string {
	return "water"
}
