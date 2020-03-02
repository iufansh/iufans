package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type IpList struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64                                         // 创建人Id
	Modifior   int64                                         // 更新人Id
	Version    int                                           // 版本
	OrgId      int64                                         // 组织ID
	Ip         string                                        // Ip
	Black      int8                                          // 是否黑名单 1：黑名单，0：白名单
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(IpList))
}

func (model *IpList) Paginate(page int, limit int, orgId int64, param1 string, black int8) (list []IpList, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(IpList))
	cond := orm.NewCondition()
	if param1 != "" {
		cond = cond.And("Ip__contains", param1)
	}
	cond = cond.And("OrgId", orgId)
	if black != -1 {
		cond = cond.And("Black", black)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("Ip", "-Id")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
