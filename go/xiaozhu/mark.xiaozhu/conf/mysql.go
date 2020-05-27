package conf

import (
	"xiaozhu/utils/libs"

	// "github.com/jinzhu/gorm"
	"github.com/golang/glog"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {

	child := DataHandle.Conf.Mysql
	Db, err := libs.NewMysqlConn(DataHandle.Conf.DbType, child.User, child.Pwd, child.Host, child.Dbname)
	if err != nil {
		glog.Exitf("conn mysql err = %+v", err)
	}

	//设置连接池
	Db.LogMode(DataHandle.Conf.Debug)
	Db.DB().SetMaxIdleConns(child.MinConns) //闲置的连接数
	Db.DB().SetMaxOpenConns(child.MaxConns) //最大打开的连接数
	//
	Db.SingularTable(true) // 全局禁用表名复数,使用TableName修改表名

	//全局重置表名规则
	//gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
	//    return "prefix_" + defaultTableName;
	//}

	DataHandle.MainDb = Db
}
