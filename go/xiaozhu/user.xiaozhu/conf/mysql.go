package conf

import (
	"fmt"

	"xiaozhu/utils/libs"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func InitDb(child MysqlChild) *gorm.DB {

	//if ins == nil {
	//    ins = &singleton{}
	//}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("defer panic:", r)
		}
	}()

	//
	_conn, err := libs.NewMysqlConn(DataHandle.Conf.DbType, child.User, child.Pwd, child.Host, child.Dbname)
	//
	if err != nil {
		fmt.Println(err)
		return nil
	}

	//设置连接池
	_conn.LogMode(DataHandle.Conf.Debug)
	_conn.DB().SetMaxIdleConns(child.MinConns) //闲置的连接数
	_conn.DB().SetMaxOpenConns(child.MaxConns) //最大打开的连接数
	//
	_conn.SingularTable(true) // 全局禁用表名复数,使用TableName修改表名

	//全局重置表名规则
	//gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
	//    return "prefix_" + defaultTableName;
	//}

	return _conn
}
