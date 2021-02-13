package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type MemberLoginCount struct {
	Id         int64     `auto`           // 自增主键
	AppNo      string    `orm:"size(63)"` // App编号
	AppChannel string    `orm:"size(63)"` // App渠道
	AppVersion int       // App版本
	CountDate  time.Time // 统计日期
	Count      int       // 统计
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(MemberLoginCount))
}

func (model *MemberLoginCount) TableUnique() [][]string {
	return [][]string{
		{"AppNo", "AppChannel", "AppVersion", "CountDate"},
	}
}
