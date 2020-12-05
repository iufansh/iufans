package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Permission struct {
	Id          int64     `auto`                              // 自增主键
	CreateDate  time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate  time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator     int64                                         // 创建人Id
	Modifior    int64                                         // 更新人Id
	Version     int                                           // 版本
	Pid         int64                                         // 父节点Id
	Enabled     int8                                          // 是否启用
	Display     int8                                          // 是否在菜单显示
	Description string    `orm:"null"`                        // 描述
	Url         string                                        // 链接地址
	Name        string    `orm:"unique;size(127)"`                      // 名称
	Icon        string    `orm:"null"`                        // 图标
	Sort        int                                           // 排序
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(Permission))
}

func GetPermissionList() (permissions []Permission) {
	var permissionList []Permission
	o := orm.NewOrm()
	qs := o.QueryTable(new(Permission))
	cond := orm.NewCondition()
	qs = qs.SetCond(cond.And("Enabled", 1))
	qs = qs.OrderBy("Pid", "Sort", "Id")
	qs.All(&permissionList)
	return permissionList
}
