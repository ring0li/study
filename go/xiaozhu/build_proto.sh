curr=$(pwd)
dir="protos"

#protoc -I${curr}/${dir} ${curr}/${dir}/user.proto --go_out=plugins=grpc:${dir} \
#   -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
#   --grpc-gateway_out=logtostderr=true:${curr}/${dir}

protoc -I${curr}/${dir} ${curr}/${dir}/*.proto --go_out=plugins=grpc:${dir}
#protoc -I${curr}/${dir} ${curr}/${dir}/user.proto --go_out=plugins=grpc:${dir}
#protoc -I${curr}/${dir} ${curr}/${dir}/clock.proto --go_out=plugins=grpc:${dir}
#protoc -I${curr}/${dir} ${curr}/${dir}/block.proto --go_out=plugins=grpc:${dir}
#protoc -I${curr}/${dir} ${curr}/${dir}/water.proto --go_out=plugins=grpc:${dir}
#protoc -I${curr}/${dir} ${curr}/${dir}/light.proto --go_out=plugins=grpc:${dir}

#删掉protoc的omitempty，但是没有效果呢。。
#ls ${curr}/${dir}/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
