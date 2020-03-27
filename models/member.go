package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Member struct {
	Id                int64     `auto`                              // 自增主键
	CreateDate        time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate        time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator           int64     // 创建人Id
	Modifior          int64     // 更新人Id
	Version           int       // 版本
	OrgId             int64     // 组织ID
	RefId             int64     // 推荐人ID
	AppNo             string
	AppChannel        string
	AppVersion        int
	Username          string
	ThirdAuthId       string // 三方登录的ID, 比如微信的unionid，华为的AuthHuaweiId
	Name              string
	Mobile            string
	Password          string
	Salt              string
	Vip               int
	Avatar            string
	Locked            int8
	LockedDate        time.Time `orm:"null"`
	LoginDate         time.Time `orm:"null"`
	LoginFailureCount int
	LoginIp           string
	Enabled           int8
	Token             string
	TokenExpTime      time.Time `orm:"null"`
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(Member))
}

func (model *Member) Paginate(page int, limit int, orderBy int, id int64, param1 string) (list []Member, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(Member))
	cond := orm.NewCondition()
	if param1 != "" {
		cond = cond.AndCond(cond.And("Name__contains", param1).Or("Username__contains", param1).Or("Mobile__contains", param1))
	}
	if id != -1 {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("RefId", id).Or("Id", id)
		cond = cond.AndCond(cond1)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	switch orderBy {
	case 1:
		qs = qs.OrderBy("Id")
		break
	case 2:
		qs = qs.OrderBy("-LoginDate")
		break
	case 3:
		qs = qs.OrderBy("LoginDate")
		break
	default:
		qs = qs.OrderBy("-Id")
		break
	}
	qs.All(&list)
	total, _ = qs.Count()
	return
}
