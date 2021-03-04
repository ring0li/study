package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"xingqi/model"
)

import _ "github.com/go-sql-driver/mysql"

var TimeZone, _ = time.LoadLocation("Asia/Shanghai")

func main() {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True", "root", "", "127.0.0.1", "xingqi"))
	//db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=%s", "root", "", "127.0.0.1", "xingqi", "Asia%2FShanghai"))
	defer db.Close()
	checkErr(err)
	//设置连接池
	db.LogMode(true)
	db.DB().SetMaxIdleConns(1) //闲置的连接数
	db.DB().SetMaxOpenConns(1) //最大打开的连接数
	db.SingularTable(true)     // 全局禁用表名复数

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")   // optionally look for config in the working directory
	err = viper.ReadInConfig() // Find and read the config file
	checkErr(err)

	chufangId := viper.GetInt("last_chufangid") //192312
	fmt.Println("上次处方id", chufangId)

	today, _ := time.ParseInLocation("2006-01-02", time.Now().In(TimeZone).Format("2006-01-02"), TimeZone)
	i := 0
	for {
		chufangId++

		//if chufangId == 209000 {
		//	os.Exit(0)
		//}

		//库里存在跳过
		order := model.Order{}
		db.Where("jiuzhenid = ? ", chufangId).First(&order)
		if order != (model.Order{}) {
			viper.Set("last_chufangid", chufangId)
			viper.WriteConfig()

			continue
		}

		//查询接口
		order = Order(db, chufangId)
		if order != (model.Order{}) {
			//今天22点查询昨天的数据
			kaidanDay, _ := time.ParseInLocation("2006-01-02", order.Kaidanshijian.Format("2006-01-02"), TimeZone)
			if today.Unix() == kaidanDay.Unix() {
				fmt.Println("开单时间是今天，终止")
				os.Exit(1)
			}

			i = 0
		} else {
			i++

			if i == 10 {
				fmt.Println("查询10次没有数据，终止")
				os.Exit(0)
			}
		}

		time.Sleep(3 * time.Second)
	}
}

func Order(db *gorm.DB, chufangId int) (order model.Order) {
	url := "https://healthcare.xqyk024.com/consult/prescription/getPrescriptionInfoNew?vaa07=" + strconv.Itoa(chufangId)
	req, err := http.NewRequest("GET", url, nil)
	checkErr(err)
	req.Header.Set("token", viper.GetString("token"))
	req.Header.Set("Content-Type", "application/json")
	fmt.Println(url)
	resp, err := (&http.Client{}).Do(req)
	defer resp.Body.Close()
	checkErr(err)
	s, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	// query
	//createTime := time.Now().In(TimeZone).Format("2006-01-02 15:04:00")
	//createDate1 := time.Now().In(TimeZone).Format("2006-01-02 00:00:00")
	//createMinute, _ := time.ParseInLocation("2006-01-02 15:04:00", createTime, TimeZone)
	//createDate, _ := time.ParseInLocation("2006-01-02 00:00:00", createDate1, TimeZone)

	var ret Ret
	json.Unmarshal(s, &ret)
	fmt.Println(ret)
	if (ret.Success == false) {
		fmt.Println("登陆失败")
		os.Exit(1)
	}

	if len(ret.Data) == 0 {
		return order //查询数据为空，可能真是空，也可能不存在这条记录
	}

	for _, v := range ret.Data {
		fmt.Println(v)

		kaidanshijian, _ := time.Parse("2006-01-02 15:04:05", v.Kaidanshijian)
		danjia, err := strconv.ParseFloat(v.Danjia, 32)
		checkErr(err)

		shuliang, _ := strconv.ParseFloat(v.Shuliang, 32)
		zongjine, _ := strconv.ParseFloat(v.Zongjine, 32)
		jiuzhenid, _ := strconv.Atoi(v.Jiuzhenid)
		order = model.Order{Zhuangtai: v.Zhuangtai, Danjia: float32(danjia), Shuliang: float32(shuliang),
			Yizhu: v.Yizhu, Zongjine: float32(zongjine), Xingming: v.Xingming, Jiuzhenid: jiuzhenid, Kaidanshijian: kaidanshijian}
		db.Create(&order)
	}

	viper.Set("last_chufangid", chufangId)
	viper.WriteConfig()

	return order
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Ret struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    []Data `json:"data"`
}

type Data struct {
	Zhuangtai     string `json:"状态"`
	Danjia        string `json:"单价"`
	Shuliang      string `json:"数量"`
	Yizhu         string `json:"医嘱"`
	Zongjine      string `json:"总金额"`
	Xingming      string `json:"名字"`
	Jiuzhenid     string `json:"就诊id"`
	Kaidanshijian string `json:"开单时间"`
}
