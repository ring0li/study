package conf

import (
	"github.com/golang/glog"
	"google.golang.org/grpc"
	pb "xiaozhu/protos"
)

var GrpcServer grpcServer

type grpcServer struct {
	UserClient  pb.UserServiceClient
	ClockClient pb.ClockServiceClient
	BlockClient pb.BlockServiceClient
	WaterClient pb.WaterServiceClient
	LightClient pb.LightServiceClient
}

func init() {
	address := DataHandle.Conf.GRPCUserServer.Host
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		glog.Exitf("run after grpc server : err = %+v", err)
	}
	//
	GrpcServer.UserClient = pb.NewUserServiceClient(conn)
	GrpcServer.ClockClient = pb.NewClockServiceClient(conn)
	GrpcServer.BlockClient = pb.NewBlockServiceClient(conn)
	GrpcServer.WaterClient = pb.NewWaterServiceClient(conn)
	GrpcServer.LightClient = pb.NewLightServiceClient(conn)
}
