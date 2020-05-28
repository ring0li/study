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

type Stock struct {
	ID          int       `gorm:"column:id;primary_key" json:"id"`
	GoodsType   string    `gorm:"column:goods_type" json:"goods_type"`
	StyleID     string    `gorm:"column:style_id" json:"style_id"`
	GoodsID     string    `gorm:"column:goods_id" json:"goods_id"`
	GoodsName   string    `gorm:"column:goods_name" json:"goods_name"`
	GoodsColor  string    `gorm:"column:goods_color" json:"goods_color"`
	GoodsPrice  string    `gorm:"column:goods_price" json:"goods_price"`
	GoodsHref   string    `gorm:"column:goods_href" json:"goods_href"`
	GoodsSize   string    `gorm:"column:goods_size" json:"goods_size"`
	GoodsStock  string    `gorm:"column:goods_stock" json:"goods_stock"`
	GoodsRemark string    `gorm:"column:goods_remark" json:"goods_remark"`
	CreateTime  null.Time `gorm:"column:create_time" json:"create_time"`
}

// TableName sets the insert table name for this struct type
func (s *Stock) TableName() string {
	return "stock"
}
