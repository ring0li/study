
├── Makefile
├── README.en.md
├── README.md
├── build_mysql_grom.sh   #把数据表生成 sturct 放到当前 model 下，要手动复制到user.xiaozhu/model中，不要手动修改，如果有特定的在 datamodels 中定义
├── build_proto.sh # 生成 user.xiaozhu grpc 的proto; 
├── cmd
│   └── gen
├── go.mod
├── go.sum
├── mark.xiaozhu
│   ├── app.yaml
│   ├── conf
│   │   ├── config.go
│   │   └── mysql.go
│   ├── controllers
│   │   └── work.go
│   ├── datamodels
│   ├── main.go
│   ├── middleware
│   │   └── auth.go
│   ├── model
│   │   └── xz_user.go
│   └── services
├── protos
│   ├── user.pb.go
│   ├── user.pb.gw.go
│   └── user.proto
├── user.xiaozhu
│   ├── app
│   │   ├── demo.go
│   │   ├── grpc_client.go
│   │   └── grpc_server.go
│   ├── app.yaml
│   ├── conf
│   │   ├── config.go
│   │   └── mysql.go
│   ├── controllers
│   │   └── user.go
│   ├── datamodels
│   ├── middleware
│   │   └── auth.go
│   ├── model
│   │   └── xz_user.go
│   └── services
└── utils  # 抽像出来放公共方法，类，特殊的放在项目目录下
    ├── common
    │   ├── Errors.go
    │   ├── format.go
    │   ├── output.go
    │   ├── tools.go
    │   └── validate.go
    └── libs
        ├── conf.go
        ├── logger.go
        ├── mysql.go
        ├── path.go
        └── redis.go

