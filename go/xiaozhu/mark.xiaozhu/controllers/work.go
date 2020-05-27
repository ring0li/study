package controllers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/silenceper/wechat/menu"
	"xiaozhu/mark.xiaozhu/services"
)

type work struct{}

var Work *work

func (this *work) MessageHandle(c echo.Context) error {
	//services.WechatService.NewWS()
	err := services.WechatService.WechatHandle(c)
	return err
}

// SetMenu 设置菜单
func (this *work) SetMenu(c echo.Context) error {
	//services.WechatService.NewWS()
	mu := services.WechatService.Wc.GetMenu()

	//buttons := make([]*menu.Button, 2)
	//btn := new(menu.Button)
	////创建click类型菜单
	//btn.SetClickButton("早安打卡", "morning")
	//buttons[0] = btn
	//btn1 := new(menu.Button)
	////创建click类型菜单
	//btn1.SetClickButton("晚安打卡", "evening")
	//buttons[1] = btn1

	//设置btn为二级菜单
	// subBtn := new(menu.Button)
	// subBtn.SetSubButton("早晚安", buttons)

	//
	subBtn := new(menu.Button)
	subBtn.SetClickButton("去打卡", "checkin")
	//
	zgbtn := new(menu.Button)
	zgbtn.SetClickButton("中国加油", "zhongguojiayou")
	whbtn := new(menu.Button)
	whbtn.SetClickButton("武汉加油", "wuhanjiayou")
	jiayouBtns := make([]*menu.Button, 2)
	jiayouBtns[0]=zgbtn
	jiayouBtns[1]=whbtn
	//设置btn为二级菜单
	subBtn2 := new(menu.Button)
	subBtn2.SetSubButton("点亮头像", jiayouBtns)

	//
	subBtns := make([]*menu.Button, 2)
	subBtns[0] = subBtn
	subBtns[1] = subBtn2

	//发送请求
	err := mu.SetMenu(subBtns)
	if err != nil {
		fmt.Printf("err= %v", err)

	} else {
		fmt.Println("推送菜单成功")
	}
	return err
}

// DelMenu 删除微信菜单
func (this *work) DelMenu(c echo.Context) error {
	//services.WechatService.NewWS()
	mu := services.WechatService.Wc.GetMenu()
	err := mu.DeleteMenu()

	if err != nil {
		fmt.Printf("err= %v", err)

	} else {
		fmt.Println("删除菜单成功")
	}
	return err
}

func (this *work) Push() error {
	//services.WechatService.Init()
	services.WechatService.PushPunchCardMsg()
	return nil
}
