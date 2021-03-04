package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"time"
	"xingqi/models"
)

import _ "github.com/go-sql-driver/mysql"

var TimeZone, _ = time.LoadLocation("Asia/Shanghai")

func main() {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=%s", "root", "", "127.0.0.1", "xingqi", "Asia%2FShanghai"))
	//db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=%s", "root", "", "127.0.0.1", "xingqi", "Asia%2FShanghai"))
	defer db.Close()
	checkErr(err)
	//设置连接池
	db.LogMode(true)
	db.DB().SetMaxIdleConns(1) //闲置的连接数
	db.DB().SetMaxOpenConns(1) //最大打开的连接数
	db.SingularTable(true)     // 全局禁用表名复数

	Guahao(db)

	Guanzhu(db)
}

func Guahao(db *gorm.DB) {
	resp, err := http.Get("https://healthcare.xqyk024.com/consult/diagnosis/doctors?pageNumber=1&pageSize=100&order=0")
	if err != nil {
		panic(err)

	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	// query
	//createTime := time.Now().In(TimeZone).Format("2006-01-02 15:04:00")
	createDate1 := time.Now().In(TimeZone).Format("2006-01-02 00:00:00")
	//createMinute, _ := time.ParseInLocation("2006-01-02 15:04:00", createTime, TimeZone)
	createDate, _ := time.ParseInLocation("2006-01-02 00:00:00", createDate1, TimeZone)

	var ret Ret
	json.Unmarshal(s, &ret)
	fmt.Println(ret.Data.Records)
	for _, v := range ret.Data.Records {
		fmt.Println(v)

		//入日志表
		//guahaoLog :=models.GuahaoLog{}
		//db.Where("create_time = ? and employee_id = ?", createMinute, v.EmployeeId).First(&guahaoLog)
		//if guahaoLog == (models.GuahaoLog{}) {
		//	guahaoLog := models.GuahaoLog{OrgId: v.OrgId, DeptId: v.DeptId, EmployeeId: v.EmployeeId, EmployeeName: v.EmployeeName, ConsultCount: v.ConsultCount, CreateTime: null.TimeFrom(createMinute)}
		//	db.Create(&guahaoLog)
		//}

		//入每天汇总表
		guahao := models.Guahao{}
		db.Where("create_date = ? and employee_id = ?", createDate, v.EmployeeId).First(&guahao)
		if guahao == (models.Guahao{}) {
			guahao = models.Guahao{OrgId: v.OrgId, DeptId: v.DeptId, EmployeeId: v.EmployeeId, EmployeeName: v.EmployeeName, ConsultCount: v.ConsultCount, CreateDate: createDate}
			db.Create(&guahao)
		} else {
			guahao.ConsultCount = v.ConsultCount
			db.Save(&guahao)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Ret struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}

type Data struct {
	Records []List `json:"records"`
}

type List struct {
	OrgId        int    `json:"orgId"`
	DeptId       int    `json:"deptId"`
	EmployeeId   int    `json:"employeeId"`
	EmployeeName string `json:"employeeName"`
	ConsultCount int    `json:"consultCount"`
}

type GuanzhuRet struct {
	Totalcount int `json:"totalcount"`
}

func Guanzhu(db *gorm.DB) {
	req, err := http.NewRequest("GET", "https://xueqiu.com/recommend/pofriends.json?type=1&code=SZ300573&start=0&count=0", nil)
	checkErr(err)
	req.Header.Set("Cookie", "aliyungf_tc=AQAAAIdSYwVjHwoAE1l5e+n0CgwxBLne; __utmc=1; _ga=GA1.2.34798019.1550112771; snbim_minify=true; stock-ad-remove=1; xq_is_login.sig=J3LxgPVPUzbBg3Kee_PquUfih7Q; xqat.sig=hXkoUMdoEfo4G4ibg3_Cb4FLZO4; xq_a_token.sig=3qSja2_fUSM5UP_o6UJKMLWaiUg; xq_r_token.sig=PvRTb_HW8e1cPeenTBSzrCZPAic; s=d511o8e6i7; bid=47f4d825d4afb001a045b48961584c74_k6lm3cir; device_id=e707bdd7b35a9f10915ee359b95c0497; Hm_lvt_fe218c11eab60b6ab1b6f84fb38bcc4a=1586407929; Hm_lpvt_fe218c11eab60b6ab1b6f84fb38bcc4a=1586407929; __utma=1.34798019.1550112771.1586435362.1586516277.1918; Hm_lvt_1db88642e346389874251b5a1eded6e3=1598969689; Hm_lpvt_1db88642e346389874251b5a1eded6e3=1598972241; acw_tc=2760820416109601078392372ef58b44e7ebe619cdb43d7f401a58b1ef2dc4; xq_a_token=176b14b3953a7c8a2ae4e4fae4c848decc03a883; xqat=176b14b3953a7c8a2ae4e4fae4c848decc03a883; xq_r_token=2c9b0faa98159f39fa3f96606a9498edb9ddac60; xq_id_token=eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ1aWQiOi0xLCJpc3MiOiJ1YyIsImV4cCI6MTYxMzQ0MzE3MSwiY3RtIjoxNjEwOTYwMTAxMjEzLCJjaWQiOiJkOWQwbjRBWnVwIn0.SOHTFWyLvBNYjoZ3uPpt84XEDyjbU-rYZgzLhi9mJ0qCmmJdSnYCFbhqGisi2y4Y1uypMPrWEU51kAdBhwMxvvlnUsPLvdjF8PJFt5j1m1BFe7NpRWpg9P_6hrEUj8yJJ0cyiXjza7EUI0ulXl0Iy-eBA8kep8yBE0sRw4mrU5RQeYecRTrXNlT34p-tSbvIyVUGs_ykUfsaQHD_03BmDLlBucPT7s6htoZoZPgJBRPJMjgJeSHuJH1gOUnxIbj_1ZjKZ9jYph3SHtNa-L8sAPt-7yiA1W3qFoeRjxQn0bVkQld0xY9YGnR76eiJhkl8rNPBUwbh2CR9P-cISR6ZrA; u=701610960107920; is_overseas=0")
	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	createDateStr := time.Now().In(TimeZone).Format("2006-01-02 00:00:00")
	createDate, _ := time.ParseInLocation("2006-01-02 00:00:00", createDateStr, TimeZone)

	var ret GuanzhuRet
	json.Unmarshal(s, &ret)
	guanzhu := models.Guanzhu{}
	db.Where("date = ?", createDate).First(&guanzhu)
	if guanzhu == (models.Guanzhu{}) {
		guanzhu = models.Guanzhu{Date: createDate, Count: ret.Totalcount}
		db.Create(&guanzhu)
	} else {
		guanzhu.Count = ret.Totalcount
		db.Save(&guanzhu)
	}
}
