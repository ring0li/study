package libs

import (
	"fmt"
	"path"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func NewMysqlConn(dbtype, user, pwd, host, dbname string) (*gorm.DB, error) {
	//
	var _conn *gorm.DB
	var err error
	//
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("defer panic:", r)
		}
	}()
	//
	if dbtype == "mysql" {
		_conn, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=%s", user, pwd, host, dbname, "Asia%2FShanghai"))
	} else {
		_conn, err = gorm.Open("sqlite3", path.Join(RootPath, dbname+".db"))
	}
	//
	if err != nil {
		fmt.Println("failed to connect database." + err.Error())
		//
		return nil, err
	}
	//设置连接池
	//_conn.LogMode(true)
	//_conn.DB().SetMaxIdleConns(5)  //闲置的连接数
	//_conn.DB().SetMaxOpenConns(100) //最大打开的连接数
	//
	_conn.SingularTable(true) // 全局禁用表名复数,使用TableName修改表名

	//全局重置表名规则
	//gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
	//    return "prefix_" + defaultTableName;
	//}

	//
	return _conn, err
}
