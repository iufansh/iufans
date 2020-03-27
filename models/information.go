package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Information struct {
	Id           int64     `auto`                              // 自增主键
	CreateDate   time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate   time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator      int64     // 创建人Id
	Modifior     int64     // 更新人Id
	Version      int       // 版本
	OrgId        int64     // 关联ID(公司id等)
	Title        string    // 标题
	Info         string    // 消息
	EffectTime   time.Time // 生效时间
	ExpireTime   time.Time // 失效时间
	NeedFeedback int8      // 0-不需要反馈；1-需要反馈
	Receiver     string    // 接收者 规则：App编号:渠道:会员ID；不填表示全部；如：a:oppo:
	ReadNum      int       // 查看人数
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(Information))
}

func (model *Information) Paginate(page int, limit int, orgId int64, param1 string) (list []Information, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(Information))
	cond := orm.NewCondition()
	cond = cond.And("OrgId", orgId)
	if param1 != "" {
		cond = cond.AndCond(cond.And("Title__contains", param1).Or("Info__contains", param1).Or("Receiver__contains", param1))
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Id")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
