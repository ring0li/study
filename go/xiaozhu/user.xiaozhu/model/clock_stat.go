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

type ClockStat struct {
	StatID            int       `gorm:"column:statId;primary_key" json:"statId"`
	UserID            int       `gorm:"column:userId" json:"userId"`
	ClockType         int       `gorm:"column:clockType" json:"clockType"`
	ClockDate         null.Time `gorm:"column:clockDate" json:"clockDate"`
	ClockTime         null.Time `gorm:"column:clockTime" json:"clockTime"`
	AvgClockTime      string    `gorm:"column:avgClockTime" json:"avgClockTime"`
	ContinueClockDays int       `gorm:"column:continueClockDays" json:"continueClockDays"`
}

// TableName sets the insert table name for this struct type
func (c *ClockStat) TableName() string {
	return "clock_stat"
}
