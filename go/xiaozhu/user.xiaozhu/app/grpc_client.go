package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
	pb "xiaozhu/protos"
)

const (
	address = "127.0.0.1:8082"
)

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//查询用户信息
	//userClient := pb.NewUserServiceClient(conn)
	//msg := &pb.GetUserByOpenidRequest{OpenId: "1"}
	//result, err := userClient.GetUserByOpenid(context.Background(), msg)

	//msg := &pb.GetUsersByOpenidsRequest{OpenIds: "45644fea-ef15-11e9-8a6d-0242c0a87003,506d20ea-ef15-11e9-8a6d-0242c0a87003"}
	//result, err := userClient.GetUsersByOpenids(context.Background(), msg)

	//保存用户信息
	//date := time.Now().Format("2006-01-02 15:04:05")
	//msg := &pb.SaveUserRequest{OpenId: "test111", UnionId: date, Nickname: date, Sex: 1, City: "北京", Province: "北京", Country: "中国",
	//	HeadImgUrl: "http://thirdwx.qlogo.cn/mmopen/v10Xv25cWiaBibBUEiak34JFPEWGXqn1icSfic6udlgP540RlbJoLuqHXiaQShDQ6lykdDBpMd6VaZbpzVA9sFm8gIX3CPUuQXsXdU/132", FriendOpenId: "45644fea-ef15-11e9-8a6d-0242c0a87003"}
	//result, err := userClient.SaveUser(context.Background(), msg)

	//取消订阅
	//userClient := pb.NewUserServiceClient(conn)
	//msg := &pb.UnsubscribeRequest{OpenId: "test"}
	//result, err := userClient.Unsubscribe(context.Background(), msg)

	//早起
	clockClient := pb.NewClockServiceClient(conn)
	date := time.Now().Format("2006-01-02 15:04:05")
	msg := &pb.GetUpRequest{OpenId: "test111", ClockTime: date, QrcodeUrl: "http://thirdwx.qlogo.cn/mmopen/v10Xv25cWiaBibBUEiak34JFPEWGXqn1icSfic6udlgP540RlbJoLuqHXiaQShDQ6lykdDBpMd6VaZbpzVA9sFm8gIX3CPUuQXsXdU/132"}
	result, err := clockClient.GetUp(context.Background(), msg)

	//获得打卡天数
	//msg := &pb.GetTotalClockDaysRequest{OpenId: "1", ClockType: 1}
	//result, err := clockClient.GetTotalClockDays(context.Background(), msg)

	//早睡
	//date := time.Now().Format("2006-01-02 15:04:05")
	//msg := &pb.SleepRequest{OpenId: "test111", ClockTime: date, QrcodeUrl: "http://thirdwx.qlogo.cn/mmopen/v10Xv25cWiaBibBUEiak34JFPEWGXqn1icSfic6udlgP540RlbJoLuqHXiaQShDQ6lykdDBpMd6VaZbpzVA9sFm8gIX3CPUuQXsXdU/132"}
	//result, err := clockClient.Sleep(context.Background(), msg)

	// --------------- 主动推送  --------------

	//未打卡的用户列表
	//clockClient := pb.NewClockServiceClient(conn)
	//msg := &pb.NoClockUsersRequest{ClockType: 1, StartTime: "00:00:00", EndTime: "23:05:00"}
	//result, err := clockClient.NoClockUsers(context.Background(), msg)

	//未打卡的好友列表
	//clockClient := pb.NewClockServiceClient(conn)
	//msg := &pb.GetNoClockFriendsRequest{OpenId: "1", ClockType: 1}
	//result, err := clockClient.GetNoClockFriends(context.Background(), msg)

	//崩溃测试
	//date := time.Now().Format("2006-01-02 15:04:05")
	//msg := &pb.GetUpRequest{OpenId: "test111", ClockTime: date, QrcodeUrl: "http://thirdwx.qlogo.cn/mmopen/v10Xv25cWiaBibBUEiak34JFPEWGXqn1icSfic6udlgP540RlbJoLuqHXiaQShDQ6lykdDBpMd6VaZbpzVA9sFm8gIX3CPUuQXsXdU/132"}
	//result, err := clockClient.Test(context.Background(), msg)

	// --------------- 屏蔽消息  --------------
	//屏蔽所有消息
	//blockClient := pb.NewBlockServiceClient(conn)
	//msg := &pb.AllRequest{OpenId: "1", Switchval: 1} //1是屏蔽，0是取消屏蔽
	//result, err := blockClient.All(context.Background(), msg)

	//屏蔽好友消息
	//blockClient := pb.NewBlockServiceClient(conn)
	//msg := &pb.FriendRequest{OpenId: "1", FriendOpenId: "2", Switchval: 0}
	//result, err := blockClient.Friend(context.Background(), msg)

	//屏蔽好友列表
	//blockClient := pb.NewBlockServiceClient(conn)
	//msg := &pb.GetMyFriendsRequest{OpenId: "1"}
	//result, err := blockClient.GetMyFriends(context.Background(), msg)

	// --------------- 8杯水  --------------
	//get喝水消息，为空表示没有或者屏蔽了
	//waterClient := pb.NewWaterServiceClient(conn)
	//msg := &pb.GetWaterTextRequest{LastUserId: 0} //userId从0开始，直到没有任何返回值，表示结束。
	//result, err := waterClient.GetWaterText(context.Background(), msg)

	//set喝水值
	//waterClient := pb.NewWaterServiceClient(conn)
	//msg := &pb.SetRequest{OpenId: "1", Waterval: 4} //次数选项。0、4、6、8。0表示屏蔽
	//result, err := waterClient.Set(context.Background(), msg)

	//get喝水值
	//waterClient := pb.NewWaterServiceClient(conn)
	//msg := &pb.GetRequest{OpenId: "1"} //次数选项。0、4、6、8。0表示屏蔽
	//result, err := waterClient.Get(context.Background(), msg)

	if err == nil {
		fmt.Println(result)
	} else {
		fmt.Println(err.Error())
	}
}
