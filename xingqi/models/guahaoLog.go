package models

import "time"

//CREATE TABLE `guahao_log` (
//`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
//`create_time` datetime DEFAULT NULL,
//`org_id` int(10) unsigned DEFAULT NULL,
//`dept_id` int(10) unsigned DEFAULT NULL,
//`employee_id` int(10) unsigned DEFAULT NULL,
//`employee_name` varchar(32) DEFAULT NULL,
//`consult_count` int(10) unsigned DEFAULT NULL,
//PRIMARY KEY (`id`),
//UNIQUE KEY `create_time_employee_id` (`create_time`,`employee_id`) USING BTREE
//) ENGINE=InnoDB AUTO_INCREMENT=22549 DEFAULT CHARSET=utf8

type GuahaoLog struct {
	Id           int       `orm:"auto"`
	CreateTime   time.Time `orm:"type(datetime)"`
	OrgId        int
	DeptId       int
	EmployeeId   int
	EmployeeName string `orm:"size(32)"`
	ConsultCount int
}

// 多字段索引
//func (u *User) TableIndex() [][]string {
//	return [][]string{
//		[]string{"Id", "Name"},
//	}
//}

//多字段唯一键
func (u *GuahaoLog) TableUnique() [][]string {
	return [][]string{
		[]string{"CreateTime", "EmployeeId"},
	}
}
