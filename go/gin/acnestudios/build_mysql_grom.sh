go get -u github.com/jinzhu/gorm #https://github.com/smallnest/gen
set -x

./cmd/gen -c "root:benchi@tcp(127.0.0.1:3306)/acnestudios?charset=utf8&parseTime=True&loc=Local" --json --gorm --guregu --rest


rm -rf api