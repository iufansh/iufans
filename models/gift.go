package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Gift struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64     // 创建人Id
	Modifior   int64     // 更新人Id
	Version    int       // 版本
	OrgId      int64     // 关联ID(公司id等)
	AppNo      string    // App编号
	Code       string    // 礼包码
	Price      int       // 价值
	Status     int       // 状态0-未使用；1-已使用
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(Gift))
}

func (model *Gift) Paginate(page int, limit int, orgId int64, param1 string, status int) (list []Gift, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(Gift))
	cond := orm.NewCondition()
	cond = cond.And("OrgId", orgId)
	if param1 != "" {
		cond = cond.And("Code__contains", param1)
	}
	if status != -1 {
		cond = cond.And("Status", status)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Id")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
