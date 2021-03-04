package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/syyongx/php2go"
)

type GuanzhuController struct {
	beego.Controller
}

type GuanzhuRet struct {
	TotalCount   int
}

func (this *GuanzhuController) Get() {
	o := orm.NewOrm()

	var orderDayMoney []OrderDayMoney
	o.Raw("SELECT DATE(kaidanshijian) `day`,sum(round(zongjine)) `money` FROM `order` WHERE zhuangtai='已发送或停止'  " +
		"GROUP BY DATE(kaidanshijian) ORDER BY DATE(kaidanshijian) DESC LIMIT 90").QueryRows(&orderDayMoney)

	var x3, y3 []interface{}
	for _, v3 := range orderDayMoney {
		x3 = append(x3, v3.Day)
		y3 = append(y3, v3.Money)
	}
	php2go.ArrayReverse(x3)
	php2go.ArrayReverse(y3)

	var orderMonthMoney []OrderMonthMoney
	o.Raw("SELECT DATE_FORMAT(kaidanshijian,'%Y年%m月') `month`,sum(round(zongjine)) `money` FROM `order` " +
		"WHERE zhuangtai='已发送或停止' GROUP BY DATE_FORMAT(kaidanshijian,'%Y年%m月') " +
		"ORDER BY DATE_FORMAT(kaidanshijian,'%Y年%m月') DESC LIMIT 90;").QueryRows(&orderMonthMoney)

	var x4, y4 []interface{}
	for _, v4 := range orderMonthMoney {
		x4 = append(x4, v4.Month)
		y4 = append(y4, v4.Money)
	}
	php2go.ArrayReverse(x4)
	php2go.ArrayReverse(y4)

	this.Data["x3"] = x3
	this.Data["y3"] = y3
	this.Data["x4"] = x4
	this.Data["y4"] = y4
	this.TplName = "order.tpl"
}
