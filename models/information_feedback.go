package models

import (
	"github.com/astaxie/beego/orm"
)

type InformationFeedback struct {
	Id            int64 `auto` // 自增主键
	Version       int   // 超过3次就等于已读
	InformationId int64
	MemberId      int64
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(InformationFeedback))
}

func (model *InformationFeedback) TableUnique() [][]string {
	return [][]string{
		{"InformationId", "MemberId"},
	}
}
