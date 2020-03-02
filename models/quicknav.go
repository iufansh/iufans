package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type QuickNav struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64                                         // 创建人Id
	Modifior   int64                                         // 更新人Id
	Version    int                                           // 版本
	Name       string                                        // 网址名称
	WebSite    string                                        // 网址
	Icon       string    `orm:"null"`                        // 图标
	Seq        int                                           // 排序(升序显示)
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(QuickNav))
}

func (model *QuickNav) Paginate(page int, limit int) (list []QuickNav, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(QuickNav))
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("Seq", "Id")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
