package middleware

import (
	"flag"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/guregu/null"
	"github.com/syyongx/php2go"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"xiaozhu/mark.xiaozhu/conf"
	"xiaozhu/user.xiaozhu/datamodels"
	"xiaozhu/user.xiaozhu/model"
	"xiaozhu/utils/common"
)

//设置打卡统计信息
func SetContinueClockDays(userId int, clockType int, clockTime time.Time) (continueClockDays int, lastClockTime time.Time) {

	clockStat := model.ClockStat{}
	conf.DataHandle.MainDb.Where("userId = ? and clockType = ?", userId, clockType).First(&clockStat)

	//和上次打卡时间做比较
	currentClockTimeDay := GetClockDate(clockTime)
	lastClockTimeDay, _ := time.ParseInLocation("2006-01-02 15:04:05", clockStat.ClockTime.Time.Format("2006-01-02 00:00:00"), common.TimeZone)
	subDays := currentClockTimeDay.Sub(lastClockTimeDay).Hours() / 24

	//变成累计打卡，一直+1即可。
	//作废：今天打过卡;昨天打卡了，今天第一次打卡;昨天没打卡，重置连续打卡天数=1
	if subDays == 0 {
		;
	} else {
		clockStat.ContinueClockDays += 1
		clockStat.ClockTime = null.TimeFrom(clockTime)
		clockStat.ClockDate = null.TimeFrom(clockTime)
		clockStat.AvgClockTime = getAvgClockTime(clockType, clockStat.UserID)
		conf.DataHandle.MainDb.Save(&clockStat)
	}

	return clockStat.ContinueClockDays, clockStat.ClockTime.Time
}

//统计多少人参与，超过多少人
func GetUserStat() (total int, clockRank int) {
	onlineTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-10-25 17:38:00", common.TimeZone)
	total = (int)(time.Now().Unix()/86400 + (time.Now().Unix()-onlineTime.Unix())/60) //1分钟增加一个人

	clockRankKey := "clock_rank_" + time.Now().In(common.TimeZone).Format("20060102")
	clockRank, _ = strconv.Atoi(conf.DataHandle.RedisClient.Get(clockRankKey).Val())
	if clockRank == 0 {
		clockRank = php2go.Rand(1, 1000)
	} else {
		clockRank++;
	}
	conf.DataHandle.RedisClient.Set(clockRankKey, strconv.Itoa(clockRank), 86400*time.Second)

	return
}

//早起图片
func GetUpImg(headImgUrl string, continueClockDays int, clockTime string, total int, num int, d string, Ym string, qrcodeUrl string) (string, string) {
	m, err := imaging.Open(circle(saveImg(headImgUrl)))
	fmt.Println(saveImg(headImgUrl))
	if err != nil {
		fmt.Printf("open file failed")
	}

	bm, err := imaging.Open(getGetUpBgImg())
	if err != nil {
		fmt.Printf("open file failed")
	}

	qrcode, err := imaging.Open(saveImg(qrcodeUrl))
	if err != nil {
		fmt.Printf("open file failed," + err.Error())
	}

	//写头像
	dst := imaging.Resize(m, 120, 120, imaging.Lanczos)     // 图片按比例缩放
	result := imaging.Overlay(bm, dst, image.Pt(60, 60), 1) // 将图片粘贴到背景图的固定位置
	//写二维码
	qrcodeResize := imaging.Resize(qrcode, 120, 120, imaging.Lanczos)
	result = imaging.Overlay(result, qrcodeResize, image.Pt(780, 780), 1)
	//写文字
	writeOnImage(result, 60, 270, "累计早起", 8)
	writeOnImage(result, 60, 360, strconv.Itoa(continueClockDays), 22)
	if continueClockDays <= 9 {
		writeOnImage(result, 120, 350, "天", 8)
	} else if continueClockDays <= 99 {
		writeOnImage(result, 180, 350, "天", 8)
	} else {
		writeOnImage(result, 240, 350, "天", 8)
	}

	writeOnImage(result, 60, 430, "今日早起", 8)
	writeOnImage(result, 60, 520, clockTime, 22)
	writeOnImage(result, 60, 580, "———————", 8)
	writeOnImage(result, 60, 630, strconv.Itoa(total)+"人正在参与", 7)
	writeOnImage(result, 60, 670, "比"+strconv.Itoa(num)+"人起得早", 7)
	writeOnImage(result, 840, 120, d, 18)
	writeOnImage(result, 830, 160, Ym, 7)

	text := getGetUpText()
	firstText := text[1]
	secondText := text[2]
	if firstText != "" {
		firstText = strings.Repeat("　", 15-utf8.RuneCountInString(firstText)) + firstText //用Printf("%s 15s,")报错
		writeOnImage(result, 290, 810, firstText, 9)
	}
	if secondText != "" {
		secondText = strings.Repeat("　", 15-utf8.RuneCountInString(secondText)) + secondText
		writeOnImage(result, 290, 850, secondText, 9)
	}

	writeOnImage(result, 570, 890, "扫码我们互道早安", 7)

	fileName := fmt.Sprintf("./img/clockImg/%s.jpg", php2go.Md5(headImgUrl))
	err = imaging.Save(result, fileName)
	if err != nil {
		return "", ""
	}

	return fileName, text[0]
}

func circle(headImgPath string) (newHeadImgPath string) {
	img, err := imaging.Open(headImgPath)
	if err != nil {
		fmt.Println(err)
	}

	img = makeCircleSmooth(img, 1)

	newHeadImgPath = "./img/headImg/" + php2go.Md5(headImgPath) + "_circle.png"
	err = imaging.Save(img, newHeadImgPath)
	if err != nil {
		fmt.Println(err)
	}
	return
}

//早睡图片
func SleepImg(headImgUrl string, continueClockDays int, clockTime string, total int, num int, d string, Ym string, qrcodeUrl string) (string, string) {
	m, err := imaging.Open(circle(saveImg(headImgUrl)))
	if err != nil {
		fmt.Printf("open file failed")
	}

	bm, err := imaging.Open(getSleepBgImg())
	if err != nil {
		fmt.Printf("open file failed")
	}

	qrcode, err := imaging.Open(saveImg(qrcodeUrl))
	if err != nil {
		fmt.Printf("open file failed")
	}

	//写头像
	dst := imaging.Resize(m, 120, 120, imaging.Lanczos)     // 图片按比例缩放
	result := imaging.Overlay(bm, dst, image.Pt(60, 60), 1) // 将图片粘贴到背景图的固定位置
	//写二维码
	dst = imaging.Resize(qrcode, 120, 120, imaging.Lanczos)
	result = imaging.Overlay(result, dst, image.Pt(780, 780), 1)
	//写文字
	writeOnImage(result, 60, 270, "坚持早睡", 8)
	writeOnImage(result, 60, 360, strconv.Itoa(continueClockDays), 22)

	if continueClockDays <= 9 {
		writeOnImage(result, 120, 350, "天", 8)
	} else if continueClockDays <= 99 {
		writeOnImage(result, 180, 350, "天", 8)
	} else {
		writeOnImage(result, 240, 350, "天", 8)
	}

	writeOnImage(result, 60, 430, "今日入睡", 8)
	writeOnImage(result, 60, 520, clockTime, 22)
	writeOnImage(result, 60, 580, "———————", 8)
	writeOnImage(result, 60, 630, strconv.Itoa(total)+"人正在参与", 7)
	writeOnImage(result, 60, 670, "比"+strconv.Itoa(num)+"人睡得早", 7)
	writeOnImage(result, 840, 120, d, 18)
	writeOnImage(result, 830, 160, Ym, 7)

	text := getSleepText()
	fmt.Println(text)
	firstText := text[1]
	secondText := text[2]
	if firstText != "" {
		firstText = strings.Repeat("　", 15-utf8.RuneCountInString(firstText)) + firstText //fmt.Println(fmt.Sprintf("%　15s", "你好啊"))，中文空格报错
		writeOnImage(result, 290, 810, firstText, 9)
	}
	if secondText != "" {
		secondText = strings.Repeat("　", 15-utf8.RuneCountInString(secondText)) + secondText
		writeOnImage(result, 290, 850, secondText, 9)
	}

	writeOnImage(result, 570, 890, "扫码我们互道晚安", 7)

	fileName := fmt.Sprintf("./img/clockImg/%s.jpg", php2go.Md5(headImgUrl))
	err = imaging.Save(result, fileName)
	if err != nil {
		return "", ""
	}

	return fileName, text[0]
}

var dpi = flag.Float64("dpi", 256, "screen resolution")

func writeOnImage(target *image.NRGBA, x int, y int, text string, fontSize int) {
	c := freetype.NewContext()

	c.SetDPI(*dpi)
	c.SetClip(target.Bounds())
	c.SetDst(target)
	c.SetHinting(font.HintingFull)

	// 设置文字颜色、字体、字大小
	c.SetSrc(image.NewUniform(color.RGBA{R: 220, G: 220, B: 220, A: 220}))
	c.SetFontSize(float64(fontSize)) //默认4
	fontFam, err := getFontFamily()
	if err != nil {
		fmt.Println("get font family error")
	}
	c.SetFont(fontFam)

	pt := freetype.Pt(x, y)
	_, err = c.DrawString(text, pt)

	if err != nil {
		fmt.Printf("draw error: %v \n", err)
	}
}

func getFontFamily() (*truetype.Font, error) {
	// 这里需要读取中文字体，否则中文文字会变成方格
	fontBytes, err := ioutil.ReadFile("./img/bg/msyh.ttc")
	//fontBytes, err := ioutil.ReadFile("./img/bg/FZY1JW.TTF")
	if err != nil {
		fmt.Println("read file error:", err)
		return &truetype.Font{}, err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println("parse font error:", err)
		return &truetype.Font{}, err
	}

	return f, err
}

func saveImg(url string) (imgPath string) {
	imgPath = "./img/headImg/" + php2go.Md5(url) + ".jpg"
	img, _ := os.Create(imgPath)
	defer img.Close()

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	b, _ := io.Copy(img, resp.Body)

	fmt.Println("File size: ", b)
	fmt.Println("File Path: ", imgPath)

	return imgPath
}

func makeCircleSmooth(src image.Image, factor float64) image.Image {
	d := src.Bounds().Dx()
	if src.Bounds().Dy() < d {
		d = src.Bounds().Dy()
	}
	dst := imaging.CropCenter(src, d, d)
	r := float64(d) / 2
	center := r - 0.5
	for x := 0; x < d; x++ {
		for y := 0; y < d; y++ {
			xf := float64(x)
			yf := float64(y)
			delta := math.Sqrt((xf-center)*(xf-center)+(yf-center)*(yf-center)) + factor - r
			switch {
			case delta > factor:
				dst.SetNRGBA(x, y, color.NRGBA{0, 0, 0, 0})
			case delta > 0:
				m := 1 - delta/factor
				c := dst.NRGBAAt(x, y)
				c.A = uint8(float64(c.A) * m)
				dst.SetNRGBA(x, y, c)
			}
		}
	}
	return dst
}

//未打卡的用户列表
func GetNoClockUsers(lastUserId int, clockType int, startTime string, endTime string) (users []datamodels.BlockFriend) {
	conf.DataHandle.MainDb.Raw("  SELECT u.userId,openId,nickname,sex,IF(blockId IS NOT NULL,1,0) isBlock FROM `user` u"+
		" INNER JOIN clock_stat cs USING(userId) "+
		" LEFT JOIN block b on u.userId= b.userId and blockType=1 "+
		" WHERE clockType = ? "+
		" 	AND clockDate != ? "+
		" 	AND avgClockTime >= ? "+
		" 	AND avgClockTime < ? "+
		" 	AND `u`.userId > ? "+
		" 	AND userLastActiveTime >= ? "+
		" 	AND blockId is null "+
		" ORDER BY `u`.userId ASC  "+
		" LIMIT 100",
		clockType, time.Now().Format("2006-01-02"), startTime, endTime,
		lastUserId, time.Now().AddDate(0, 0, -2).Format("2006-01-02 15:04:05")).Find(&users)
	return
}

func getGetUpBgImg() (bgImgPath string) {
	today, _ := strconv.Atoi(time.Now().Format("20060102"))
	bgImgPath = "./img/bg/getup/" + strconv.Itoa(today%14+1) + ".png"
	return bgImgPath
}

func getSleepBgImg() (bgImgPath string) {
	today, _ := strconv.Atoi(time.Now().Format("20060102"))
	bgImgPath = "./img/bg/sleep/" + strconv.Itoa(today%20+1) + ".png"
	return bgImgPath
}

func getGetUpText() []string {

	text := [][]string{
		//{"一二三四五六七八九十一二三四五", "", ""},
		{"有梦就有希望，为梦想加油！", "", "有梦就有希望，为梦想加油"},
		{"存好心，做好事，说好话，做好人。", "存好心，做好事", "说好话，做好人"},
		{"发奋努力的背后，必有加倍的赏赐。", "", "发奋努力的背后，必有加倍的赏赐"},
		{"没有伞的孩子，必须努力奔跑！", "", "没有伞的孩子，必须努力奔跑"},
		{"最大的破产是绝望，最大的资产是希望。", "最大的破产是绝望", "最大的资产是希望"},
		{"想象力比知识更重要。", "", "想象力比知识更重要"},
		{"信仰是伟大的情感，一种创造力量。", "", "信仰是伟大的情感，一种创造力量"},
		{"真正的才智是刚毅的志向。", "", "真正的才智是刚毅的志向"},
		{"含泪播种的人一定能含笑收获。", "", "含泪播种的人一定能含笑收获"},
		{"心若阳光，无谓悲伤。", "", "心若阳光，无谓悲伤"},
		{"像蚂蚁一样工作，像蝴蝶一样生活。", "", "像蚂蚁一样工作，像蝴蝶一样生活"},
		{"失败是因为放弃得太快，加油！", "", "失败是因为放弃得太快，加油"},
		{"值得做的事情，都值得把它做好。", "", "值得做的事情，都值得把它做好"},
		{"保持善良，拥有远方。", "", "保持善良，拥有远方"},
		{"向前奔跑，只为了心中的美好。", "", "向前奔跑，只为了心中的美好"},
		{"每天进步一点点，希望的火苗不熄灭。", "每天进步一点点", "希望的火苗不熄灭"},
		{"不怕路远，就怕志短。", "", "不怕路远，就怕志短。"},
	}

	return text[php2go.Rand(0, len(text)-1)]
}

func getSleepText() []string {
	//因为世界很圆，所以终会遇见。
	//爱是陪我们行走一生的行李。
	//记忆迷了路，用爱找回来。
	//时间, 让爱更了解爱。
	//你忘记的, 我都记得。
	//说不出你的好，但就是谁都代替不了。
	//你来了，我的心就满了。
	//注定要去的地方，不管多晚都有光。
	//过期的旧书，不过期的求知欲。
	//在现实断裂的地方，梦，汇成了海洋。
	//每个人都是生活的导演。
	//最好的道别就是明天见。
	//用心，是生活最好的质感。

	text := [][]string{
		//{"一二三四五六七八九十一二三四五", "", ""},
		{"不为失败找理由，只为成功找方法。", "", "不为失败找理由，只为成功找方法"},
		{"含泪播种的人一定能含笑收获。", "", "含泪播种的人一定能含笑收获"},
		{"人海中，认定了你，这便是我的执着。", "人海中，认定了你", "这便是我的执着"},
		{"最好的日子是你在闹，我在笑，如此温暖过一生。", "最好的日子是你在闹，我在笑", "如此温暖过一生"},
		{"永远不要放弃你真正想要的东西，等待虽难，但后悔更甚。", "永远不要放弃你真正想要的东西", "等待虽难，但后悔更甚"},
		{"无论走多远，家都是起点。", "", "无论走多远，家都是起点"},
		{"爱在日常，才不寻常。", "", "爱在日常，才不寻常"},
		{"有前进的梦想，也有回家的方向。", "", "有前进的梦想，也有回家的方向"},
		{"世界再大，不过你我之间。", "", "世界再大，不过你我之间"},
	}
	return text[php2go.Rand(0, len(text)-1)]
}

func getAvgClockTime(clockType int, userId int) (avgClockTime string) {
	clockStat := model.ClockStat{}

	//取消ONLY_FULL_GROUP_BY
	//conf.DataHandle.MainDb.Raw("set @@sql_mode='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'")

	if clockType == 1 {
		conf.DataHandle.MainDb.Raw("select if(avgClockTime > '07:00:00','07:00:00',avgClockTime) avgClockTime from ("+
			"SELECT ifnull(FROM_unixtime(avg(unix_timestamp(clockTime)),'%H:%i:%s'),'07:00:00') avgClockTime FROM ("+
			"	SELECT FROM_unixtime(unix_timestamp(any_value(clockTime)),'2000-01-01 %H:%i:%s') clockTime"+
			"		FROM ("+
			"			SELECT * FROM clock_log	WHERE clockType=1 AND clockDate >= ? AND userId = ? ORDER BY clockTime ASC"+
			"		) t1 "+
			"	GROUP BY clockDate) t2"+
			") t3",
			time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04:05"), userId).Find(&clockStat)

	} else {
		conf.DataHandle.MainDb.Raw("select if(avgClockTime > '2000-01-01 20:30:00','20:30:00',FROM_unixtime(unix_timestamp(avgClockTime),'%H:%i:%s')) avgClockTime from ("+
			"SELECT ifnull(FROM_unixtime(avg(unix_timestamp(clockTime)),'%Y-%m-%d %H:%i:%s'),'2000-01-01 20:30:00') avgClockTime FROM ("+
			"	SELECT IF(FROM_unixtime(unix_timestamp(any_value(clockTime)),'%H')< 4 ,FROM_unixtime(unix_timestamp(any_value(clockTime)),'2000-01-02 %H:%i:%s'),	FROM_unixtime(unix_timestamp(any_value(clockTime)),'2000-01-01 %H:%i:%s')) clockTime"+
			" 		FROM ("+
			"			SELECT * FROM clock_log	WHERE clockType=2 AND clockDate >= ? AND userId = ? ORDER BY clockTime ASC"+
			"		) t1"+
			"	GROUP BY clockDate) t2"+
			") t3",
			time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04:05"), userId).Find(&clockStat)
	}
	return clockStat.AvgClockTime
}

func GetClockDate(clockTime time.Time) time.Time {
	if clockTime.Hour() < 4 {
		clockTime = clockTime.AddDate(0, 0, -1)
	}
	clockDate, _ := time.Parse("2006-01-02 15:04:05", clockTime.Format("2006-01-02 00:00:00"))
	return clockDate
}
