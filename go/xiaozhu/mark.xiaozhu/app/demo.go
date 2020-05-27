package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"xiaozhu/mark.xiaozhu/conf"
	"xiaozhu/mark.xiaozhu/controllers"

	//xiaozhuMiddleware "xiaozhu/app/user/middleware"
	//"xiaozhu/config"
)

func main() {
	//conf.InitConf("../app.yaml")
	fmt.Println(conf.DataHandle.Conf.Redis.Host)
	//
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://wx.iiwoo.com", "http://192.168.8.102:8080","http://127.0.01"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())

	da := e.Group("/wechat")
	//da.Use(xiaozhuMiddleware.SystemUserAuth)
	da.GET("/pushmsg", controllers.Work.MessageHandle)

	e.Logger.Fatal(e.Start(":8080"))
}
