package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type NormalQuestion struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64     // 创建人Id
	Modifior   int64     // 更新人Id
	Version    int       // 版本
	OrgId      int64     // 组织ID
	Seq        int       // 排序
	Question   string    // 问题
	Answer     string    // 答复
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(NormalQuestion))
}

func (model *NormalQuestion) Paginate(page int, limit int, param1 string) (list []NormalQuestion, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(NormalQuestion))
	cond := orm.NewCondition()
	if param1 != "" {
		cond = cond.AndCond(cond.And("Question__contains", param1).Or("Answer__contains", param1))
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("Seq")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
