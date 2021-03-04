package models

import (
	"time"
)

type Guanzhu struct {
	Id           int       `orm:"auto"`
	Date   time.Time       `orm:"type(date)"`
	Count        int
}