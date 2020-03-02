package models

import (
	"github.com/astaxie/beego/orm"
)

type RolePermission struct {
	Id           int64
	PermissionId int64
	RoleId       int64
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(RolePermission))
}

func (model *RolePermission) TableUnique() [][]string {
	return [][]string{
		{"PermissionId", "RoleId"},
	}
}
