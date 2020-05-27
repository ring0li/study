package main

import (
	_ "xiaozhu/mark.xiaozhu/conf"
	//
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"xiaozhu/mark.xiaozhu/controllers"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://wx.iiwoo.com", "http://127.0.0.1:8080"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())

	//静态目录
	e.Static("/", "../static")

	//u
	da := e.Group("/wechat")
	//da.Use(xiaozhuMiddleware.SystemUserAuth)
	da.POST("/pushmsg", controllers.Work.MessageHandle)
	da.GET("/pushmsg", func(c echo.Context) error {
		return c.String(http.StatusOK, "you see see you, one day day de…\n")
	})
	da.GET("/setmenu", controllers.Work.SetMenu)
	da.GET("/delmenu", controllers.Work.DelMenu)

	//set switch 消息开关设置
	sw := e.Group("/switchset")
	sw.POST("/GetMsgReceiveSwitch", controllers.Switchset.GetMsgReceiveSwitch)
	sw.POST("/SetMsgReceiveSwitch", controllers.Switchset.SetMsgReceiveSwitch)
	sw.POST("/GetFriendList", controllers.Switchset.GetFriendList)
	sw.POST("/SetFriendMsgReceiveSwitch", controllers.Switchset.GetMsgReceiveSwitch)
	sw.POST("/DelFriend", controllers.Switchset.DelFriend)
	sw.POST("/GetWaterSwitch", controllers.Switchset.GetWaterSwitch)
	sw.POST("/SetWaterSwitch", controllers.Switchset.SetWaterSwitch)

	//wx认证文件
	web := e.Group("/web")
	web.Static("/", "web")
	//da.Static("/pushmsg/MP_verify_05lPsU8wCPSFWjJt.txt", "../static/MP_verify_05lPsU8wCPSFWjJt.txt")

	//build task
	//go controllers.BuildTask.Build()

	//
	//e.HideBanner = false
	//e.Server.Addr = ":"+config.AppConf.Port
	//e.Server.Addr = ":80"

	// 实时消息推送
	go func() { controllers.Work.Push() }()
	// 计划任务
	go func() { controllers.CronWork.Run() }()
	//e.Logger.Fatal(gracehttp.Serve(e.Server)) //windows不支持
	e.Logger.Fatal(e.Start(":80"))
}
