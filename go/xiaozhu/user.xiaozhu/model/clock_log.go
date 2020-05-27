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

type ClockLog struct {
	LogID     int       `gorm:"column:logId;primary_key" json:"logId"`
	UserID    int       `gorm:"column:userId" json:"userId"`
	ClockType int       `gorm:"column:clockType" json:"clockType"`
	ClockTime null.Time `gorm:"column:clockTime" json:"clockTime"`
	ClockDate null.Time `gorm:"column:clockDate" json:"clockDate"`
}

// TableName sets the insert table name for this struct type
func (c *ClockLog) TableName() string {
	return "clock_log"
}
