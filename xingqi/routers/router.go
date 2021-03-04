package routers

import (
	"xingqi/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/guahao", &controllers.GuahaoController{})
    beego.Router("/liuli", &controllers.OrderController{})
    //beego.Router("/", &controllers.MainController{})
}
