package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	pb "xiaozhu/protos"

	"github.com/labstack/echo/v4"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/material"
	"github.com/silenceper/wechat/message"
	"github.com/silenceper/wechat/qr"
	"xiaozhu/mark.xiaozhu/conf"
)

const (
	morningStartTime int = 4
	morningEndTime   int = 12
	eveningStartTime int = 18
	eveningEndTime   int = 4
	//
	morningPunchType int = 1
	eveningPunchType int = 2
)

var (
	WechatService wechatService
)

type wechatService struct {
	Wc *wechat.Wechat
}

func init() {
	WechatService.Wc = conf.XzWechat
}

func (this *wechatService) WechatHandle(c echo.Context) error {
	// 传入request和responseWriter
	server := this.Wc.GetServer(c.Request(), c.Response())

	//设置接收消息的处理方法
	fmt.Println("消息预定义开始！")
	server.SetMessageHandler(func(reqInfo message.MixMessage) *message.Reply {
		replyMsgType := message.MsgTypeText
		var replyMsgData interface{}

		switch reqInfo.MsgType {

		case message.MsgTypeText: //文本消息
			replyMsgType, replyMsgData = this.getTextMsg(reqInfo)

		//case message.MsgTypeImage: //图片消息
		//case message.MsgTypeVoice: //语音消息
		//case message.MsgTypeVideo: //视频消息
		//case message.MsgTypeShortVideo: //小视频消息
		//case message.MsgTypeLocation: //地理位置消息
		case message.MsgTypeLink: //链接消息
			replyMsgType, replyMsgData = this.getTextMsg(reqInfo)
		case message.MsgTypeEvent: //事件推送消息
			fmt.Println(reqInfo.Event)
			switch reqInfo.Event {
			case message.EventSubscribe: //关注消息
				replyMsgType, replyMsgData = this.setSubscribe(reqInfo)
			case message.EventUnsubscribe: //事件推送消息
				replyMsgType, replyMsgData = this.unSetSubscribe(reqInfo)
			case message.EventClick:
				replyMsgType, replyMsgData = this.getEventClickInfo(reqInfo)
			//case message.MsgTypeImage: //图片消息
			//case message.MsgTypeVoice: //语音消息
			//case message.MsgTypeVideo: //视频消息
			//case message.MsgTypeShortVideo: //小视频消息
			//case message.MsgTypeLocation: //地理位置消息
			//case message.MsgTypeLink: //链接消息
			default:
				replyMsgType, replyMsgData = this.getDefaultText(reqInfo.FromUserName)
			}
		default:
			replyMsgType = message.MsgTypeText
			replyMsgData = message.NewText("")
		}
		return &message.Reply{MsgType: replyMsgType, MsgData: replyMsgData}
	})

	fmt.Println("消息预定义结束！")
	//处理消息接收以及回复
	//server.SetDebug(true)
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return err
	}
	//发送回复的消息
	err = server.Send()
	fmt.Println("消息接收结束！")
	return err
}

// 根据文本返回相应的数据
func (this *wechatService) getTextMsg(reqInfo message.MixMessage) (message.MsgType, interface{}) {
	fmt.Println("msg_receive:", reqInfo.Content)
	replyType := message.MsgTypeText
	var replyData interface{}
	switch reqInfo.Content {
	case "去打卡":
		replyType, replyData = this.getPunchCardText(reqInfo.FromUserName)
	default:
		replyType, replyData = this.getDefaultText(reqInfo.FromUserName)
	}

	return replyType, replyData
}

//默认文本文案
func (this *wechatService) getDefaultText(openId string) (message.MsgType, interface{}) {
	userInfo, err := conf.GrpcServer.UserClient.GetUserByOpenid(context.Background(), &pb.GetUserByOpenidRequest{
		OpenId: openId,
	})
	if err != nil {
		fmt.Println("getDefaultText:", err)
	}

	replyMsg := "「" + userInfo.GetNickname() + "」" + "小主您好，给小主您请安。\n\n"
	replyMsg += "醒了记得打卡开启美好的一天哦~\n\n<a href='weixin://bizmsgmenu?msgmenucontent=去打卡&msgmenuid=101'>去打卡</a>\n\n"

	replyData := message.NewText(replyMsg)
	return message.MsgTypeText, replyData
}

// 获取打卡信息
func (this *wechatService) getPunchCardText(openId string) (message.MsgType, interface{}) {
	userInfo, err := conf.GrpcServer.UserClient.GetUserByOpenid(context.Background(), &pb.GetUserByOpenidRequest{
		OpenId: openId,
	})
	if err != nil {
		fmt.Println("conf.GrpcServer.UserClient.getUserInfoByOpenId err:", err)
	}
	replyMsg := ""
	curHour := time.Now().Hour()
	if curHour >= morningStartTime && curHour <= morningEndTime {
		replyMsg = this.getMorningPunchCardText(userInfo)
	} else if curHour >= eveningStartTime || curHour < eveningEndTime {
		replyMsg = this.getEveningPunchCardText(userInfo)
	} else {
		replyMsg = this.getDefaultPunchCardText(userInfo)
	}

	replyData := message.NewText(replyMsg)
	return message.MsgTypeText, replyData
}

// 晚打卡
func (this *wechatService) getEveningPunchCardText(userInfo *pb.GetUserByOpenidResponse) string {
	replyMsg := "「" + userInfo.GetNickname() + "」"
	redisOpenIdPunchStatusKey := userInfo.OpenId + "_e_punchstatus"
	redisOpenIdMediaIDValue := this.Wc.Context.Cache.Get(redisOpenIdPunchStatusKey)
	// 下午卡
	if redisOpenIdMediaIDValue != nil {
		replyMsg += "小主您吉祥，您已经打过晚安卡了。\n\n"
		replyMsg += getMorningPunchCardTitle() + fmt.Sprintf("\n\n<a href=\"%s?openid=%s\">提醒管理</a>", conf.WEB_TIPSWITCH, userInfo.OpenId)
	} else {
		punchCardNum := this.getPunchCardNum(userInfo.OpenId, int32(eveningPunchType))
		replyMsg += "小主晚安，您已经坚持打卡" + strconv.Itoa(int(punchCardNum)) + "天了，加油了。\n\n"
		replyMsg += getMorningPunchCardTitle() + fmt.Sprintf("\n\n<a href=\"%s?openid=%s\">提醒管理</a>", conf.WEB_TIPSWITCH, userInfo.OpenId)
	}
	//当前时间小于24点，用 加上第二日4点前的时间戳，减去当前时间戳
	sleepTimeOut := 0
	if time.Now().Hour() < 24 && time.Now().Hour() >= eveningStartTime {
		sleepTimeOut = 86400 + (eveningEndTime+1)*3600 - time.Now().Hour()*3600 - time.Now().Minute()*60 - time.Now().Second()
	} else {
		sleepTimeOut = (eveningEndTime+1)*3600 - time.Now().Hour()*3600 - time.Now().Minute()*60 - time.Now().Second()
	}
	mediaIdTimeout := time.Duration(sleepTimeOut) * time.Second
	fmt.Println("getEveningPunchCardText punchstatus:", sleepTimeOut)
	if sleepTimeOut > 0 {
		this.setRedisPunchCardStatus(redisOpenIdPunchStatusKey, userInfo, strconv.Itoa(eveningPunchType), mediaIdTimeout)
	}
	return replyMsg
}

// 早打卡
func (this *wechatService) getMorningPunchCardText(userInfo *pb.GetUserByOpenidResponse) string {
	replyMsg := "「" + userInfo.GetNickname() + "」"
	redisOpenIdPunchStatusKey := userInfo.OpenId + "_m_punchstatus"
	redisOpenIdPunchStatus := this.Wc.Context.Cache.Get(redisOpenIdPunchStatusKey)
	// 上午卡
	if redisOpenIdPunchStatus != nil {
		replyMsg += "小主您吉祥，您已经打过早安卡了。\n\n"
	} else {
		replyMsg += getMorningPunchCardTitle() + fmt.Sprintf("\n\n<a href=\"%s?openid=%s\">提醒管理</a>", conf.WEB_TIPSWITCH, userInfo.OpenId)
	}
	//12点前的时间戳，减去当前时间戳
	getUpTimeOut := (morningEndTime+1)*3600 - time.Now().Hour()*3600 - time.Now().Minute()*60 - time.Now().Second()
	mediaIdTimeout := time.Duration(getUpTimeOut) * time.Second
	fmt.Println("getMorningPunchCardText punchstatus:", getUpTimeOut)

	if getUpTimeOut > 0 {
		this.setRedisPunchCardStatus(redisOpenIdPunchStatusKey, userInfo, strconv.Itoa(morningPunchType), mediaIdTimeout)
	}
	return replyMsg
}

func getMorningPunchCardTitle() string {
	w := int(time.Now().Weekday())
	title := ""
	switch w {
	case 1:
		title = "今天是星期一要努力的挣开双眼哦〜"
	case 2:
		title = "星期二，要打起精神、把活干完"
	case 3:
		title = "周三已经到了，加油哦〜"
	case 4:
		title = "周四了，要准备好迎接周末了。"
	case 5:
		title = "周五了，抵档不住内心的“哈哈哈”"
	case 6:
		title = "今天是周六，来呀〜 造作啊〜"
	case 0:
		title = "周日了，啊啊啊〜，抓紧抓紧〜"
	}
	return title
}

// 打卡接口
func (this *wechatService) punchCard(openId string, punchCardType int) (string, string) {
	qrcodeUrl := this.getQrUrl(openId)
	//早起
	date := time.Now().Format("2006-01-02 15:04:05")
	var imgUrl string
	var waterText string
	if punchCardType == morningPunchType {
		getUpResult, err := conf.GrpcServer.ClockClient.GetUp(context.Background(),
			&pb.GetUpRequest{
				OpenId:    openId,
				ClockTime: date,
				QrcodeUrl: qrcodeUrl,
			})
		if err == nil {
			fmt.Println(getUpResult)
		} else {
			fmt.Println(err.Error())
		}
		imgUrl = getUpResult.GetImgUrl()
		waterText = getUpResult.GetWaterText()
	} else {
		sleepResult, err := conf.GrpcServer.ClockClient.Sleep(context.Background(),
			&pb.SleepRequest{
				OpenId:    openId,
				ClockTime: date,
				QrcodeUrl: qrcodeUrl,
			})
		if err == nil {
			fmt.Println(sleepResult)
		} else {
			fmt.Println(err.Error())
		}
		imgUrl = sleepResult.GetImgUrl()
		waterText = sleepResult.GetWaterText()
	}
	fmt.Println("imgUrl:", imgUrl)

	mt := this.Wc.GetMaterial()
	uploadInfo, err := mt.MediaUpload(material.MediaTypeImage, imgUrl)
	if err != nil {
		fmt.Println("MediaUpload:", err)
	} else {
		fmt.Println(uploadInfo)
	}
	return uploadInfo.MediaID, waterText
}

// 获取用户二维码
func (this *wechatService) getQrUrl(openId string) string {
	defUrl := "https://gss1.bdstatic.com/-vo3dSag_xI4khGkpoWK1HF6hhy/baike/c0%3Dbaike92%2C5%2C5%2C92%2C30/sign=82e604c0252dd42a4b0409f9625230d0/314e251f95cad1c8f24aad8e7f3e6709c83d51a1.jpg"
	redisClient := conf.DataHandle.RedisClient
	redisQrUrl := redisClient.Get("QRURL_" + openId).Val()

	if len(redisQrUrl) > 0 {
		fmt.Println("getRedisQrUrl succ:", redisQrUrl, openId)
		return redisQrUrl
	}
	wxQr := this.Wc.GetQR()
	// 三十天的时间戳，秒
	exp := 30 * 24 * time.Hour
	newTq := qr.NewTmpQrRequest(exp, openId)
	tk, err := wxQr.GetQRTicket(newTq)
	if err != nil {
		fmt.Println("GetQRTicket err:", err)
		return defUrl
	}
	qrUrl := qr.ShowQRCode(tk)
	fmt.Println("getTicketQrUrl succ:", qrUrl, openId)
	// redis的过期时间是 25天，residExp
	residExp := 25 * 24 * time.Hour
	redisClient.Set("QRURL_"+openId, qrUrl, residExp).Result()
	return qrUrl
}

// 默认打卡文案
func (this *wechatService) getDefaultPunchCardText(userInfo *pb.GetUserByOpenidResponse) string {
	replyMsg := "「" + userInfo.GetNickname() + "」"

	// 未到打卡时间
	replyMsg += "小主您吉祥，让我们来记录您的起居，让小主您养成良好的习惯。\n\n"
	replyMsg += "早打卡时间：" + strconv.Itoa(morningStartTime) + ":00 ~ " + strconv.Itoa(morningEndTime) + ":00\n\n"
	replyMsg += "晚打卡时间：" + strconv.Itoa(eveningStartTime) + ":00 ~ " + strconv.Itoa(eveningEndTime) + ":00\n\n"

	return replyMsg
}

// SetSubscribe 关注事件
func (this *wechatService) setSubscribe(reqInfo message.MixMessage) (message.MsgType, interface{}) {
	fmt.Println("FromUserName:", reqInfo.FromUserName)
	userInfo, err := this.Wc.GetUser().GetUserInfo(reqInfo.FromUserName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(userInfo)

	msg := &pb.SaveUserRequest{
		OpenId:       userInfo.OpenID,
		UnionId:      userInfo.UnionID,
		Nickname:     userInfo.Nickname,
		HeadImgUrl:   userInfo.Headimgurl,
		Sex:          userInfo.Sex,
		City:         userInfo.City,
		Province:     userInfo.Province,
		Country:      userInfo.Country,
		FriendOpenId: userInfo.QrSceneStr,
	}
	result, err := conf.GrpcServer.UserClient.SaveUser(context.Background(), msg)

	if err == nil {
		fmt.Println(result)
	} else {
		fmt.Println(err.Error())
	}
	replyType := message.MsgTypeText

	replyMsg := "「" + userInfo.Nickname + "」，小主您吉祥，请让我们来记录您的起居。让小主您变的越来越好。\n\n"
	replyMsg += "<a href='weixin://bizmsgmenu?msgmenucontent=去打卡&msgmenuid=101'>去打卡</a>\n\n"

	replyData := message.NewText(replyMsg)
	return replyType, replyData
}

// UnSetSubscribe 取消关注事件
func (this *wechatService) unSetSubscribe(reqInfo message.MixMessage) (message.MsgType, interface{}) {

	msg := &pb.UnsubscribeRequest{
		OpenId: reqInfo.FromUserName,
	}
	result, err := conf.GrpcServer.UserClient.Unsubscribe(context.Background(), msg)

	if err == nil {
		fmt.Println(result)
	} else {
		fmt.Println(err.Error())
	}
	replyType := message.MsgTypeText
	replyData := message.NewText("欢迎光临")
	return replyType, replyData
}

// GetEventClickInfo 处理菜单事件
func (this *wechatService) getEventClickInfo(reqInfo message.MixMessage) (message.MsgType, interface{}) {
	openId := reqInfo.FromUserName
	userInfo, err := conf.GrpcServer.UserClient.GetUserByOpenid(context.Background(), &pb.GetUserByOpenidRequest{
		OpenId: openId,
	})
	if err != nil {
		fmt.Println("getPunchCardText:", err)
	}
	replyMsg := ""
	msgType := message.MsgTypeText
	curHour := time.Now().Hour()
	switch reqInfo.EventKey {
	// case "morning":
	// case "evening":
	case "checkin":
		if curHour >= morningStartTime && curHour <= morningEndTime {
			//morning
			replyMsg = this.getMorningPunchCardText(userInfo)
		} else if curHour >= eveningStartTime || curHour < eveningEndTime {
			//evening
			replyMsg = this.getEveningPunchCardText(userInfo)
		} else if curHour > morningEndTime {
			replyMsg = this.getDefaultMorningPunchCardText(userInfo)
		} else {
			replyMsg = this.getDefaultEveningPunchCardText(userInfo)
		}
	case "zhongguojiayou":
		fallthrough
	case "wuhanjiayou":

		replyMsg = "小主您吉祥，头像生成中..."
		//
		go func() {
			//
			_addTagRequest := pb.AddTagRequest{
				OpenId:   openId,
				TypeName: reqInfo.EventKey,
				Qrcode:   this.getQrUrl(openId),
			}
			imgResult, err := conf.GrpcServer.LightClient.AddTag(context.Background(), &_addTagRequest)
			if err == nil {
				cms := message.NewMessageManager(this.Wc.Context)
				//
				mt := this.Wc.GetMaterial()
				uploadInfo, err := mt.MediaUpload(material.MediaTypeImage, imgResult.Imgurl)
				if err == nil {
					msgData := message.NewCustomerImgMessage(openId, uploadInfo.MediaID)
					cms.Send(msgData)
				}
				if imgResult.Shareimgurl != "" {
					uploadInfo, err := mt.MediaUpload(material.MediaTypeImage, imgResult.Shareimgurl)
					if err == nil {
						msgData := message.NewCustomerImgMessage(openId, uploadInfo.MediaID)
						cms.Send(msgData)
					}
				}
			}
		}()

	default:
		replyMsg = "小主您吉祥!"
	}

	replyData := message.NewText(replyMsg)
	return msgType, replyData
}

// 推送每日提醒
func (this wechatService) PushPunchCardMsg() {
	redisClient := conf.DataHandle.RedisClient
	for {
		pushKey := "pushPunchCardList"
		popLen := redisClient.LLen(pushKey).Val()
		if popLen <= 0 {
			//time.Sleep(time.Duration(10) * time.Second)
			time.Sleep(time.Duration(200) * time.Millisecond)
			continue
		}
		pushValue := make(map[string]string)
		baseResult, err := redisClient.RPop(pushKey).Bytes()
		if err != nil {
			fmt.Println("get pushPunchCardList error:", err)
			continue
		}
		err = json.Unmarshal(baseResult, &pushValue)
		if err != nil {
			fmt.Println("json error:", err)
		}
		fmt.Println(pushValue)

		openId := pushValue["openId"]
		punchCardType, _ := strconv.Atoi(pushValue["punchCardType"])
		// userIdInt, _ := strconv.Atoi(pushValue["userId"])
		// userId := int32(userIdInt)
		nickname := pushValue["nickname"]
		redisOpenIdMediaIDKey := openId + "_mediakey_" + pushValue["punchCardType"]
		redisOpenIdMediaIDValue := this.Wc.Context.Cache.Get(redisOpenIdMediaIDKey)
		mediaId := ""
		var waterText string //打卡后提醒8杯水
		if redisOpenIdMediaIDValue != nil {
			mediaId = redisOpenIdMediaIDValue.(string)
		} else {
			mediaId, waterText = this.punchCard(openId, punchCardType)
			this.Wc.Context.Cache.Set(redisOpenIdMediaIDKey, mediaId, 12*time.Hour)
			this.pushMsgToFriend(openId, nickname)
		}

		cms := message.NewMessageManager(this.Wc.Context)
		msgData := message.NewCustomerImgMessage(openId, mediaId)
		err = cms.Send(msgData)
		if err != nil {
			fmt.Println("PushPunchCardMsg:", err)
		} else {
			if len(waterText) > 0 {
				waterData := message.NewCustomerTextMessage(openId, waterText)
				err = cms.Send(waterData)
			}
		}
	}
}

// 打卡提醒
func (this wechatService) CronPunchCardRemind(stepTime int) {
	if (time.Now().Hour() > morningEndTime && time.Now().Hour() < eveningStartTime) || (time.Now().Hour() < morningStartTime && time.Now().Hour() > eveningEndTime) {
		//return
	}

	var lastUserId int32 = 0
	var clockType int32 = 1
	var startTime = time.Now().Format("15:04:00")
	var durStepTime = time.Duration(stepTime)
	var endTime = time.Now().Add(durStepTime * time.Minute).Format("15:04:00")
	var sendMsg string = "小主早安，早起的鸟儿先得食，去打卡开启美好的一天。\n\n<a href='weixin://bizmsgmenu?msgmenucontent=去打卡&msgmenuid=101'>去打卡</a>"
	if time.Now().Hour() >= eveningStartTime || time.Now().Hour() < eveningEndTime {
		clockType = 2
		sendMsg = "小主该就寝了，保证充足睡眠才可以变的美丽健康哦〜\n\n<a href='weixin://bizmsgmenu?msgmenucontent=去打卡&msgmenuid=102'>去打卡</a>"
	}
	fmt.Println("CronPunchCardRemind Start Get User")

	for {
		clockUser := &pb.NoClockUsersRequest{ClockType: clockType, LastUserId: lastUserId, StartTime: startTime, EndTime: endTime}
		result, err := conf.GrpcServer.ClockClient.NoClockUsers(context.Background(), clockUser)

		if err != nil {
			fmt.Println(err.Error())
			break
		}

		if len(result.List) <= 0 {
			fmt.Println("推送打卡提示已完成：", clockType)
			break
		}

		for _, userInfo := range result.List {
			cms := message.NewMessageManager(this.Wc.Context)
			sendMsgText := "「" + userInfo.GetNickname() + "」" + sendMsg
			msgData := message.NewCustomerTextMessage(userInfo.OpenId, sendMsgText)
			err = cms.Send(msgData)
			if err != nil {
				fmt.Println("PunchCardRemind error:", userInfo.UserId, err)
			} else {
				fmt.Println("PunchCardRemind succ:", userInfo.UserId)
			}
			lastUserId = userInfo.UserId
		}
		if lastUserId > 0 {
			continue
		} else {
			fmt.Println("PunchCardRemind Complete")
			break
		}
	}
	return
}

// 打卡后准备通知好友
func (this *wechatService) pushMsgToFriend(openId, nickname string) {
	redisClient := conf.DataHandle.RedisClient
	pushKey := "pushMsgToFriendList"
	pushValue := map[string]string{
		// "userid":   strconv.Itoa(int(userId)),
		"openid":   openId,
		"nickname": nickname,
	}
	jsonPushValue, _ := json.Marshal(pushValue)
	redisClient.LPush(pushKey, jsonPushValue)
}

// 本人打开后通知好友
func (this *wechatService) CronPushMsgToFriend() {
	type pushFriend struct {
		myName      string
		friendNames string
		num         int32
	}
	curHour := time.Now().Hour()
	var (
		clockType        int32 = 1
		startMsg         string
		endMsg           string
		pushUserDataList = make(map[string]*pushFriend, 10)
	)
	if curHour >= morningStartTime && curHour < morningEndTime {
		clockType = 1
		startMsg = "小主早安，您的朋友:"
		endMsg = `已经起床打卡了，小主不要忘记打卡哦〜\n\n<a href="weixin://bizmsgmenu?msgmenucontent=去打卡&msgmenuid=101">去打卡</a>`
	} else if curHour >= eveningStartTime || curHour < eveningEndTime {
		startMsg = "小主您好，你的好友:"
		endMsg = `已经打卡休息了，忙碌了一整天，睡前别忘打卡哦~\n\n<a href="weixin://bizmsgmenu?msgmenucontent=去打卡&msgmenuid=101">去打卡</a>`
		clockType = 2
	} else {
		return
	}
	redisClient := conf.DataHandle.RedisClient

	pushKey := "pushMsgToFriendList"
	for i := 0; i < 1000; i++ {
		popLen := redisClient.LLen(pushKey).Val()
		if popLen <= 0 {
			break
		}
		baseResult, err := redisClient.RPop(pushKey).Bytes()
		if err != nil {
			fmt.Println("get pushPunchCardList error:", err)
			continue
		}
		pushValue := map[string]string{}
		err = json.Unmarshal(baseResult, &pushValue)
		if len(pushValue) <= 0 {
			fmt.Println("get pushValue isempty:")
			continue
		}
		fromOpenId := pushValue["openid"]
		fromNickname := pushValue["nickname"]

		if fromOpenId != "" {
			fmt.Println("get pushValue err:", err, fromOpenId)
			continue
		}
		fmt.Println("get pushValue:", fromOpenId, fromNickname)

		reqData := &pb.GetNoClockFriendsRequest{OpenId: fromOpenId, ClockType: clockType}
		friendList, err := conf.GrpcServer.ClockClient.GetNoClockFriends(context.Background(), reqData)

		if err == nil {
			fmt.Println(fromOpenId, " friendList:", len(friendList.List))
		} else {
			fmt.Println(err.Error())
		}
		for _, friend := range friendList.List {
			_, ok := pushUserDataList[friend.OpenId]
			if ok {
				if pushUserDataList[friend.OpenId].num < 5 {
					pushUserDataList[friend.OpenId].friendNames += "、" + fromNickname
				}
				pushUserDataList[friend.OpenId].num++
			} else {
				newPushFriend := pushFriend{
					myName:      friend.Nickname,
					friendNames: fromNickname,
					num:         int32(1),
				}
				fmt.Println(newPushFriend)
				pushUserDataList[friend.OpenId] = &newPushFriend
			}
		}
	}
	cms := message.NewMessageManager(this.Wc.Context)

	for friendOpenId, toVal := range pushUserDataList {

		sendMsgText := "「" + toVal.myName + "」" + startMsg + toVal.friendNames
		//if toVal.num > 1 {
		//	sendMsgText += " 等" + strconv.Itoa(int(toVal.num)) + "位好友"
		//} else {
		//	sendMsgText += " "
		//}
		sendMsgText += endMsg
		sendMsgText += fmt.Sprintf("\n\n<a href=\"%s?openid=%s\">好友管理</a>", conf.WEB_FRIENDSWITCH, friendOpenId)
		msgData := message.NewCustomerTextMessage(friendOpenId, sendMsgText)
		err := cms.Send(msgData)
		if err != nil {
			fmt.Println("CronPushMsgToFriend error:", friendOpenId, err)
		} else {
			fmt.Println("CronPushMsgToFriend succ:", friendOpenId)
		}
	}
}

// 8杯水提醒
func (this *wechatService) CronPushWaterRemind() {

	var lastUserId int32 = 0

	var openId string
	var remindText string
	cms := message.NewMessageManager(this.Wc.Context)
	for {
		msg := &pb.GetWaterTextRequest{LastUserId: lastUserId} //userId从0开始，直到没有任何返回值，表示结束。
		result, err := conf.GrpcServer.WaterClient.GetWaterText(context.Background(), msg)
		if err != nil {
			fmt.Println("获取8杯水提醒错误")
			break
		}
		if len(result.List) <= 0 {
			fmt.Println("该时段没有提醒：", time.Now().Format("2006-01-02 15:04:05"))
			break
		}

		for _, water := range result.List {
			openId = water.GetOpenId()
			remindText = water.GetWaterText()
			remindText += fmt.Sprintf("\n\n<a href=\"%s?openid=%s\">8杯水管理</a>", conf.WEB_WATER8SWITCH, openId)
			lastUserId = water.GetUserId()
			msgData := message.NewCustomerTextMessage(openId, remindText)
			err := cms.Send(msgData)
			if err != nil {
				fmt.Println("CronPushMsgToFriend error:", openId, err)
			} else {
				fmt.Println("CronPushMsgToFriend succ:", openId, remindText)
			}
		}
		if lastUserId <= 0 {
			break
		}
		time.Sleep(500 * time.Microsecond)
	}
}

func (this *wechatService) getDefaultMorningPunchCardText(userInfo *pb.GetUserByOpenidResponse) string {
	replyMsg := "「" + userInfo.GetNickname() + "」"
	// 未到打卡时间
	replyMsg += "小主您吉祥，" + fmt.Sprintf("%02d", morningStartTime) + ":00-" + fmt.Sprintf("%02d", morningEndTime) + ":00才可以打早安卡哦~~"

	return replyMsg
}

func (this *wechatService) getDefaultEveningPunchCardText(userInfo *pb.GetUserByOpenidResponse) string {
	replyMsg := "「" + userInfo.GetNickname() + "」"
	// 未到打卡时间
	replyMsg += "小主您吉祥，" + fmt.Sprintf("%02d", eveningStartTime) + ":00-" + fmt.Sprintf("%02d", eveningEndTime) + ":00才可以打晚安卡哦~~"

	return replyMsg
}

func (this *wechatService) setRedisPunchCardStatus(redisOpenIdPunchStatusKey string, userInfo *pb.GetUserByOpenidResponse, punchCardType string, mediaIdTimeout time.Duration) {
	fmt.Println("punchCard:", userInfo.OpenId, punchCardType)
	err := this.Wc.Context.Cache.Set(redisOpenIdPunchStatusKey, punchCardType, mediaIdTimeout)
	if err != nil {
		fmt.Println(err)
		return
	}
	redisClient := conf.DataHandle.RedisClient
	if err != nil {
		fmt.Println("redis conn :", err)
	} else {
		pushKey := "pushPunchCardList"
		pushValue := map[string]string{
			"openId":        userInfo.OpenId,
			"punchCardType": punchCardType,
			"userId":        strconv.Itoa(int(userInfo.UserId)),
			"nickname":      userInfo.Nickname,
		}
		jsonPushValue, _ := json.Marshal(pushValue)
		redisClient.LPush(pushKey, jsonPushValue).Result()
	}
}

// 获取打卡天数
func (this *wechatService) getPunchCardNum(openId string, punchCardType int32) int32 {
	msg := &pb.GetTotalClockDaysRequest{OpenId: openId, ClockType: punchCardType}
	result, err := conf.GrpcServer.ClockClient.GetTotalClockDays(context.Background(), msg)
	if err == nil {
		fmt.Println(result)
	} else {
		fmt.Println(err.Error())
	}
	return result.GetTotalClockDays()
}
