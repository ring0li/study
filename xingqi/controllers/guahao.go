package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/syyongx/php2go"
)

type GuahaoController struct {
	beego.Controller
}

type Week struct {
	Key   string
	Count int
}

type Day struct {
	Key   string
	Count int
}

type Month struct {
	Key   string
	Count int
}

type Jidu struct {
	Key   string
	Count int
}

type Guanzhu struct {
	Key   string
	Count int
}

func (this *GuahaoController) Get() {
	limit := this.GetString("limit")
	if limit == "" {
		limit = "30"
	}

	o := orm.NewOrm()

	var days []Day
	o.Raw("SELECT create_date `key`,sum(`consult_count`) `count` FROM guahao  GROUP BY create_date ORDER BY create_date desc LIMIT " + limit).
		QueryRows(&days)

	var x1, y1 []interface{}

	for _, v1 := range days {
		x1 = append(x1, v1.Key)
		y1 = append(y1, v1.Count)
	}
	php2go.ArrayReverse(x1)
	php2go.ArrayReverse(y1)

	var weeks []Week

	o.Raw("SELECT DATE_FORMAT(create_date,'%Y年第%u周') `key`,sum(consult_count) `count` FROM guahao" +
		" GROUP BY DATE_FORMAT(create_date,'%Y年第%u周') " +
		" ORDER BY DATE_FORMAT(create_date,'%Y年第%u周')   desc " +
		" LIMIT " + limit).QueryRows(&weeks)

	var x2 []interface{}
	var y2 []interface{}
	firstWeekCount := 0
	for _, v2 := range weeks {
		//第0周，可能属于是上一年年末的星期，进行合并
		if php2go.Strstr(v2.Key, "第00周") != "" {
			firstWeekCount = v2.Count
		} else if php2go.Strstr(v2.Key, "第53周") != "" {
			x2 = append(x2, v2.Key)
			y2 = append(y2, v2.Count+firstWeekCount)
		} else {
			x2 = append(x2, v2.Key)
			y2 = append(y2, v2.Count)
		}
	}
	php2go.ArrayReverse(x2)
	php2go.ArrayReverse(y2)

	var months []Month
	o.Raw("SELECT DATE_FORMAT(create_date,'%Y年%m月') `key`,sum(consult_count) `count` FROM guahao " +
		"GROUP BY DATE_FORMAT(create_date,'%Y年%m月') ORDER BY DATE_FORMAT(create_date,'%Y年%m月')   DESC LIMIT " + limit).
		QueryRows(&months)

	var x3, y3 []interface{}

	for _, v3 := range months {
		x3 = append(x3, v3.Key)
		y3 = append(y3, v3.Count)
	}
	php2go.ArrayReverse(x3)
	php2go.ArrayReverse(y3)

	var jidus []Jidu
	o.Raw("SELECT concat(date_format(create_date, '%Y年'),FLOOR((date_format(create_date, '%m')+2)/3),'季度') `key`," +
		" sum(consult_count) `count` FROM guahao " +
		" GROUP BY concat(date_format(create_date, '%Y年'),FLOOR((date_format(create_date, '%m')+2)/3),'季度')" +
		" ORDER BY concat(date_format(create_date, '%Y年'),FLOOR((date_format(create_date, '%m')+2)/3),'季度')  DESC" +
		" LIMIT " + limit).QueryRows(&jidus)

	var x4, y4 []interface{}

	for _, v4 := range jidus {
		x4 = append(x4, v4.Key)
		y4 = append(y4, v4.Count)
	}
	php2go.ArrayReverse(x4)
	php2go.ArrayReverse(y4)

	var guanzhus []Guanzhu
	o.Raw("SELECT `date` `key`,sum(`count`) `count` FROM guanzhu  GROUP BY `DATE` ORDER BY `date` desc LIMIT " + limit).
		QueryRows(&guanzhus)

	var guanzhuX, guanzhuY []interface{}

	for _, guanzhuV := range guanzhus {
		guanzhuX = append(guanzhuX, guanzhuV.Key)
		guanzhuY = append(guanzhuY, guanzhuV.Count)
	}
	php2go.ArrayReverse(guanzhuX)
	php2go.ArrayReverse(guanzhuY)

	this.Data["x1"] = x1
	this.Data["y1"] = y1
	this.Data["x2"] = x2
	this.Data["y2"] = y2
	this.Data["x3"] = x3
	this.Data["y3"] = y3
	this.Data["x4"] = x4
	this.Data["y4"] = y4
	this.Data["guanzhuX"] = guanzhuX
	this.Data["guanzhuY"] = guanzhuY

	this.TplName = "guahao.tpl"
}
