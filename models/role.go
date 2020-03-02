package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Role struct {
	Id          int64     `auto`                              // 自增主键
	CreateDate  time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate  time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator     int64                                         // 创建人Id
	Modifior    int64                                         // 更新人Id
	Version     int                                           // 版本
	Enabled     int8                                          // 是否启用
	Description string    `orm:"null"`                        // 描述
	IsSystem    int8                                          // 是否内置(内置不可选择)
	Name        string                                        // 名称
	IsOrg       int8                                          // 是否适用组织
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(Role))
}

func GetRoleList(isOrg bool) (roles []Role) {
	var roleList []Role
	o := orm.NewOrm()
	qs := o.QueryTable(new(Role)).Filter("Enabled", 1).Filter("IsSystem", 0)
	if isOrg {
		qs = qs.Filter("IsOrg", 1)
	}
	qs.All(&roleList)
	return roleList
}
