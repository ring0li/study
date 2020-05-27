set -x

cd ../../
git pull
cd mark.xiaozhu/app/
go build main.go
mv main /home/work/wwwroot/app/mark
killall mark
#cd /home/work/wwwroot/app/
#nohup ./mark >> mark.log 2>&1 &
