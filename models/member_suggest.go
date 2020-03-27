package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type MemberSuggest struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64     // 创建人Id
	Modifior   int64     // 更新人Id
	Version    int       // 版本
	OrgId      int64     // 组织ID
	MemberId   int64     // 会员ID
	Mobile     string
	Name       string
	Suggest    string
	Status     int    // 0：未处理；1：已回复接受未读；2：已回复拒绝未读；3：接受建议已读；4：拒绝建议已读
	Feedback   string // 回复
	AppInfo    string // App信息
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(MemberSuggest))
}

func (model *MemberSuggest) Paginate(page int, limit int, param1 string, status int) (list []MemberSuggest, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(MemberSuggest))
	cond := orm.NewCondition()
	if param1 != "" {
		cond = cond.AndCond(cond.And("Name__contains", param1).Or("Mobile__contains", param1).Or("Suggest__contains", param1))
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
