package models

import (
	"time"
)

//`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
//`create_date` date DEFAULT NULL,
//`org_id` int(10) unsigned DEFAULT NULL,
//`dept_id` int(10) unsigned DEFAULT NULL,
//`employee_id` int(10) unsigned DEFAULT NULL,
//`employee_name` varchar(32) DEFAULT NULL,
//`consult_count` int(10) unsigned DEFAULT NULL,
//PRIMARY KEY (`id`) USING BTREE,
//UNIQUE KEY `create_time_employee_id` (`create_date`,`employee_id`) USING BTREE

type Guahao struct {
	Id           int       `orm:"auto"`
	CreateDate   time.Time `orm:"type(date)"`
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
func (u *Guahao) TableUnique() [][]string {
	return [][]string{
		[]string{"CreateDate", "EmployeeId"},
	}
}

func (m *Guahao) week() {
	//o := orm.NewOrm()
}
