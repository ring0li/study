set -x

cd ../../
git pull
cd user.xiaozhu/app/
go build grpc_server.go
mv grpc_server /home/work/wwwroot/app/
killall grpc_server
#cd /home/work/wwwroot/app/
#nohup ./grpc_server >> grpc_server.log 2>&1 &
