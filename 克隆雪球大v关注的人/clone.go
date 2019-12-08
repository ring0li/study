package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)
import "github.com/gocolly/colly"

type GuanzhuList struct {
	Count   int                `json:"count"`
	Page    int                `json:"page"`
	MaxPage int                `json:"maxPage"`
	Users   []GuanzhuListUsers `json:"users"`
}

type GuanzhuListUsers struct {
	Following   bool   `json:"following"`
	Screen_name string `json:"screen_name"`
	Id          int    `json:"id"`
}

type Userinfo struct {
	User []UserinfoUser `json:"user"`
}

type UserinfoUser struct {
	Screen_name string `json:"screen_name"`
}

type GroupList struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GuanzhuSession struct {
	Session_token string `json:"session_token"`
}

type GuanzhuReturn struct {
	Success           bool   `json:"success"`
	Error_code        string `json:"error_code"` // {"error_description":"你已经把此用户加入黑名单，加关注前请先解除","error_uri":"/friendships/create/3945809122.json","error_data":null,"error_code":"22008"}
	Error_description string `json:"error_description"`
}

// 实际中应该用更好的变量名
var (
	h      bool
	uid    int
	cookie string
)

func init() {
	flag.IntVar(&uid, "uid", 0, "大v的uid,比如不明真相的群众的uid是1955602780")
	flag.StringVar(&cookie, "cookie", "", "自己雪球的cookie，获取方法参考此目录下的图片")

	// 改变默认的 Usage，flag包中的Usage 其实是一个函数类型。这里是覆盖默认函数实现，具体见后面Usage部分的分析
	flag.Usage = usage
}

func usage() {
	filename := ""
	switch runtime.GOOS {
	case "darwin":
		filename = "./clone"
	case "windows":
		filename = "clone.exe"
	case "linux":
		filename = "./clone"
	}
	fmt.Fprintf(os.Stderr, `克隆大v关注的人 version 1.0
用法: `+filename+` -cookie "cookie" -uid 8160302897

参数说明:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if h || uid == 0 || cookie == "" {
		flag.Usage()
		os.Exit(1)
	}

	for page := 1; ; page++ {
		guanzhuList := guanzhuList(uid, page) //电总uid

		for _, v := range guanzhuList.Users {
			if v.Following == false {
				ret := guanzhu(v.Id)
				if ret.Success == true {
					fmt.Println("关注成功", v.Screen_name)
				} else {
					if ret.Error_code == "22802" {
						fmt.Println(ret.Error_description)
						os.Exit(1)
					} else {
						fmt.Println("关注失败", v.Screen_name, "失败原因：", ret.Error_description)
					}
				}
				//group(v.Id, gid)

				time.Sleep(1 * time.Second)
			} else {
				fmt.Println("已关注过", v.Screen_name, )
			}

		}

		//最后一页
		if guanzhuList.Page == guanzhuList.MaxPage {
			break;
		}
	}

}

func guanzhuList(uid int, page int) GuanzhuList {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"),

	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", cookie)
		r.Headers.Set("Content-Type", "application/json;charset=UTF-8")
		//fmt.Println("Visiting", r.URL, r.Body)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("返回值", string(r.Body))
		check(err)
	})

	guanzhuList := GuanzhuList{}
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println("返回值", string(r.Body))

		check(json.Unmarshal(r.Body, &guanzhuList))
	})

	c.Visit("https://xueqiu.com/friendships/groups/members.json?uid=" + strconv.Itoa(uid) + "&page=" + strconv.Itoa(page) + "&gid=0&count=20&_=" + rand())
	return guanzhuList
}

//设置用户组
func group(uid int, gid int) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"),

	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", cookie)
		//fmt.Println("Visiting", r.URL, r.Body)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("返回值", string(r.Body))
		check(err)
	})

	c.OnResponse(func(r *colly.Response) {
		//fmt.Println("返回值", string(r.Body))
	})

	form := map[string]string{"uid": strconv.Itoa(uid), "gid": strconv.Itoa(gid)}
	c.Post("https://xueqiu.com/friendships/groups/members/update.json", form)
}

func guanzhu(uid int) (guanzhuReturn GuanzhuReturn) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"),

	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", cookie)
		//fmt.Println("Visiting", r.URL, r.Body)
	})

	c.OnError(func(r *colly.Response, err error) {
		check(json.Unmarshal(r.Body, &guanzhuReturn))
	})

	c.OnResponse(func(r *colly.Response) {
		check(json.Unmarshal(r.Body, &guanzhuReturn))
	})

	c.Post("https://xueqiu.com/friendships/create/"+strconv.Itoa(uid)+".json?_="+rand(), nil)
	return
}

func remark(uid int) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"),

	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", cookie)
		//fmt.Println("Visiting", r.URL, r.Body)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("返回值", string(r.Body))
		check(err)
	})

	c.OnResponse(func(r *colly.Response) {
		//fmt.Println("返回值", string(r.Body))
	})

	c.Post("https://xueqiu.com/friendships/create/"+strconv.Itoa(uid)+".json", nil)
}

func userinfo(uid int) (userinfo Userinfo) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"),

	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", cookie)
		r.Headers.Set("Content-Type", "application/json;charset=UTF-8")
		//fmt.Println("Visiting", r.URL, r.Body)
	})

	c.OnError(func(r *colly.Response, err error) {
		//fmt.Println("返回值", string(r.Body))
		check(json.Unmarshal(r.Body, &userinfo))
	})

	c.OnResponse(func(r *colly.Response) {
		//fmt.Println("返回值", string(r.Body))

		check(json.Unmarshal(r.Body, &userinfo))
	})

	c.Visit("https://xueqiu.com/statuses/original/show.json?user_id=" + strconv.Itoa(uid))
	return
}

func groupList(uid int) (groupList *[]GroupList) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"),

	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "device_id=de5a2ccddb9f4b115c6c8c299334c936; _ga=GA1.2.34798019.1550112771; s=er15y021au; bid=47f4d825d4afb001a045b48961584c74_js425sdm; snbim_minify=true; Hm_lvt_fe218c11eab60b6ab1b6f84fb38bcc4a=1571549042; Hm_lpvt_fe218c11eab60b6ab1b6f84fb38bcc4a=1573608618; xq_token_expire=Tue%20Dec%2024%202019%2010%3A21%3A52%20GMT%2B0800%20(China%20Standard%20Time); u=5343943747; xq_is_login=1; xq_is_login.sig=J3LxgPVPUzbBg3Kee_PquUfih7Q; u.sig=zovzbwBLzBHTWWu9NKY23rv7D9c; xqat=8dfc4bc2788208c1ed12d9bf33aebc5a09dfefa5; xqat.sig=V0ATlxvrewcfyBI2hlwAQ7YwR5w; xq_a_token=8dfc4bc2788208c1ed12d9bf33aebc5a09dfefa5; xq_a_token.sig=N2T6nQ55MfLHCkIms9-UaNbbtKQ; xq_r_token=90d3d24c1574654cb8efdab1265ee628c0c6066a; xq_r_token.sig=52mUa9vuJgBQjakhPNrMrx1DPh0; Hm_lvt_1db88642e346389874251b5a1eded6e3=1575430820,1575431742,1575431751,1575683776; Hm_lpvt_1db88642e346389874251b5a1eded6e3=1575689980")
		r.Headers.Set("Content-Type", "application/json;charset=UTF-8")
		fmt.Println("Visiting", r.URL, r.Body)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("返回值", string(r.Body))
		check(err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("返回值", string(r.Body))

		check(json.Unmarshal(r.Body, &groupList))
	})

	c.Visit("https://xueqiu.com/friendships/groups.json?tid=" + strconv.Itoa(uid) + "&_=1575712298029")
	return
}
func check(err error) {
	if err != nil {
		fmt.Println("错误", err.Error())
		os.Exit(1)
	}
}

func rand() string {
	return strconv.Itoa(int(time.Now().Unix())) + "000"
}

/**
1，查询大v的关注列表，循环
2，自己是否关注这个人，没有关注就关注
*/
