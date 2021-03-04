package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"xingqi/models"
	_ "xingqi/routers"
)

func init() {
	// set default database
	orm.RegisterDataBase("default", "mysql", "root:@tcp(127.0.0.1:3306)/xingqi?charset=utf8&loc=Asia%2FShanghai")
	//orm.RegisterDataBase("default", "mysql", "root:@tcp(127.0.0.1:3306)/xingqi?charset=utf8&loc=Asia%2FShanghai")
	// 需要在init中注册定义的model
	orm.RegisterModel(new(models.Guahao), new(models.GuahaoLog))

	orm.RunSyncdb("default", false, true)
}

func main() {
	orm.Debug = true

	orm.RunCommand() //接收命令参数

	//o := orm.NewOrm()
	//o.Using("default") // 默认使用 default，你可以指定为其他数据库

	//StaticDir["/static"] = "static"

	beego.Run(":30573")
}
