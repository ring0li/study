package middleware

import (
	"time"
	"xiaozhu/mark.xiaozhu/conf"
	"xiaozhu/user.xiaozhu/datamodels"
	"xiaozhu/user.xiaozhu/model"
	"xiaozhu/utils/common"
)

func GetWaterByUserId(userId int) (water model.Water) {
	conf.DataHandle.MainDb.Where("userId = ?", userId).First(&water)
	return
}

func GetWaterUsers(lastUserId int) (waterTextList []datamodels.WaterTextList) {
	conf.DataHandle.MainDb.Raw("  SELECT u.userId,openId,nickname,ifnull(`option`,0) `option` FROM `user` u"+
		" LEFT JOIN water w USING(userId)"+
		" WHERE u.userId > ?"+
		" AND userLastActiveTime >= ? "+
		" AND `option` != 0 "+
		" ORDER BY userId ASC"+
		" LIMIT 100", lastUserId,
		time.Now().In(common.TimeZone).AddDate(0, 0, -2).Format("2006-01-02 15:04:05")).
		Find(&waterTextList)
	return
}

func GetDrinkWaterText(userId int, nickname string) (text string) {
	option := GetWaterOption(userId)
	if option == 0 {
		return ""
	}

	//没打早卡，返回空
	//clockStat := model.ClockStat{}
	//conf.DataHandle.MainDb.Where("userId = ? and clockType = 1 and clockDate = ?",
	//	userId, time.Now().In(common.TimeZone).Format("2006-01-02")).First(&clockStat)
	//if clockStat == (model.ClockStat{}) {
	//	return ""
	//}

	//打过晚卡，返回空
	//conf.DataHandle.MainDb.Where("userId = ? and clockType = 2 and clockDate = ?",
	//	userId, time.Now().In(common.TimeZone).Format("2006-01-02")).First(&clockStat)
	//if clockStat != (model.ClockStat{}) {
	//	return ""
	//}

	hour := time.Now().Hour()
	if option == 4 {
		if hour == 12 {
			return "「" + nickname + "」小主，中午吃过午饭之后，来杯淡盐水既补水又消炎，小小的眯一下，可以让下午精神百倍，精神百倍。"
		} else if hour == 15 {
			return "「" + nickname + "」小主，下午喝杯茶提神养颜哦，如果有失眠的状态那就改喝温水哦〜"
		}
	} else if option == 6 {
		if hour == 9 {
			return "「" + nickname + "」小主，喝杯温水补补水哦💦，也可以通过其他的方式给自己补水，例如：加湿器、保湿喷雾等。"
		} else if hour == 12 {
			return "「" + nickname + "」小主，中午吃过午饭之后，来杯淡盐水既补水又消炎，小小的眯一下，可以让下午精神百倍，精神百倍。"
		} else if hour == 15 {
			return "「" + nickname + "」小主，下午喝杯茶提神养颜哦，如果有失眠的状态那就改喝温水哦〜"
		} else if hour == 17 {
			return "「" + nickname + "」小主，喝杯清水助消化，好准备吃晚饭啦。"
		}

	} else if option == 8 {
		if hour == 8 {
			return "「" + nickname + "」小主，水是生命之源。喝水是生命体通过口腔摄入水以补充自身细胞内水份，是生命体新陈代谢的重要一环，也是补充微量元素的方式之一"
		} else if hour == 9 {
			return "「" + nickname + "」小主，喝杯温水补补水哦💦，也可以通过其他的方式给自己补水，例如：加湿器、保湿喷雾等。"
		} else if hour == 12 {
			return "「" + nickname + "」小主，中午吃过午饭之后，来杯淡盐水既补水又消炎，小小的眯一下，可以让下午精神百倍，精神百倍。"
		} else if hour == 15 {
			return "「" + nickname + "」小主，下午喝杯茶提神养颜哦，如果有失眠的状态那就改喝温水哦〜"
		} else if hour == 17 {
			return "「" + nickname + "」小主，喝杯清水助消化，好准备吃晚饭啦。"
		} else if hour == 20 {
			return "「" + nickname + "」小主，清水增进血液循环，这个时间可以给自己的肌肤补补水，放松放松。"
		}
	}

	return ""
}

func GetWaterOption(userId int) int {
	water := GetWaterByUserId(userId)
	if water == (model.Water{}) {
		return 0
	}
	return water.Option
}

func GetWaterGetUpText(userId int, nickname string) string {
	option := GetWaterOption(userId)
	if option == 0 {
		return ""
	}

	return "「" + nickname + "」小主，早起第一杯水，见意喝蜂蜜|牛奶水排毒养颜哦〜"
}

func GetWaterSleepText(userId int, nickname string) string {
	option := GetWaterOption(userId)
	if option == 0 {
		return ""
	}

	if option == 4 {
		return "「" + nickname + "」小主，清水增进血液循环，这个时间可以给自己的肌肤补补水，放松放松。"
	} else if option == 6 {
		return "「" + nickname + "」小主，牛奶补水安神，晚间来一杯牛奶可以舒缓紧张的心情，让自己变得平静下来、进入梦乡。"
	} else if option == 8 {
		return "「" + nickname + "」小主，牛奶补水安神，晚间来一杯牛奶可以舒缓紧张的心情，让自己变得平静下来、进入梦乡。"
	}

	return ""
}
