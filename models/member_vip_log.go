package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type MemberVipLog struct {
	Id           int64     `auto`                              // 自增主键
	CreateDate   time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate   time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator      int64     // 创建人Id
	Modifior     int64     // 更新人Id
	Version      int       // 版本
	OrderNo      string    `orm:"unique"` // 订单号
	MemberId     int64     // 会员ID
	VipDays      int       // 购买时长
	PreVip       int       // 上回VIP等级
	PreVipTime   time.Time `orm:"null"` // 上回VIP获得时间
	PreVipExpire time.Time `orm:"null"` // 上回VIP过期时间
	CurVip       int       // 当前VIP等级
	CurVipTime   time.Time `orm:"null"` // 当前VIP获得时间
	CurVipExpire time.Time `orm:"null"` // 当前VIP过期时间
	Remark       string    // 备注
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(MemberVipLog))
}

func (model *MemberVipLog) Paginate(page int, limit int, memberId int64) (list []MemberVipLog, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(MemberVipLog))
	cond := orm.NewCondition()
	if memberId > 0 {
		cond = cond.And("MemberId", memberId)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Id")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
