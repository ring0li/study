package main

import (
	"fmt"
	"gin/acnestudios/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"time"
)

func main() {

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=%s", "root", "", "127.0.0.1", "acnestudios", "Asia%2FShanghai"))
	defer db.Close()
	checkErr(err)
	db.LogMode(true)
	//设置连接池
	db.DB().SetMaxIdleConns(5)  //闲置的连接数
	db.DB().SetMaxOpenConns(10) //最大打开的连接数
	//
	db.SingularTable(true) // 全局禁用表名复数

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {

		day := c.DefaultQuery("day", "")
		name := c.DefaultQuery("goods_name", "")
		color := c.DefaultQuery("goods_color", "")
		size := c.DefaultQuery("goods_size", "")
		styleId := c.DefaultQuery("style_id", "")

		sql := "select * from stock where 1";
		if day != "" {
			today, _ := time.Parse("2006-01-02 15:04:05", day+" 00:00:00")
			tomorrow := today.Add(86400 * time.Second)
			sql += " and create_time >= '" + today.Format("2006-01-02 15:04:05") + "' and create_time < '" + tomorrow.Format("2006-01-02 15:04:05") + "'"
		}

		if name != "" {
			sql += " and goods_name like '%" + name + "%' "
		}
		if color != "" {
			sql += " and goods_color = '" + color + "' "
		}
		if size != "" {
			sql += " and goods_size = '" + size + "' "
		}
		if styleId != "" {
			sql += " and style_id = '" + styleId + "' "
		}

		sql += " order by id desc limit 50";
		stocks := []model.Stock{}
		db.Raw(sql).Scan(&stocks)

		c.HTML(http.StatusOK, "index.tmpl", map[string]interface{}{
			"list":        stocks,
			"goods_name":  name,
			"goods_color": color,
			"goods_size":  size,
			"day":         day,
			"style_id":    styleId,
		})
	})

	router.Run(":8080")
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
