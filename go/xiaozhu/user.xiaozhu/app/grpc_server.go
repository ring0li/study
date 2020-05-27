package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	servers "xiaozhu/user.xiaozhu/controllers"

	"golang.org/x/net/context"
	pb "xiaozhu/protos"
)

func main() {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	// Register interceptor
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(c context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("请求:", info, req)
		resp, err = handler(c, req)
		fmt.Println("返回:", resp)
		return resp, err
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	s := grpc.NewServer(opts...)

	pb.RegisterUserServiceServer(s, &servers.UserServer{})
	pb.RegisterClockServiceServer(s, &servers.ClockServer{})
	pb.RegisterBlockServiceServer(s, &servers.BlockServer{})
	pb.RegisterWaterServiceServer(s, &servers.WaterServer{})
	pb.RegisterLightServiceServer(s, &servers.LightServer{})

	s.Serve(lis)
}
