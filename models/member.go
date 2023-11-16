package models

import (
	"github.com/astaxie/beego"
	"strings"
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
	Levels            string    // 层级关系
	LevelsDeep        int       // 层级深度
	AppNo             string
	AppChannel        string
	AppVersion        int
	Username          string `orm:"unique;size(127)"`
	ThirdAuthId       string // 三方登录的ID, 比如微信的unionid，华为的AuthHuaweiId
	RegType           int    // 注册类型 1-手机号；2-微信；3-支付宝；4-QQ；5-本机号码一键登录；6-Apple登录；7-游客模式
	Name              string
	Mobile            string
	Password          string
	Salt              string
	Vip               int
	VipTime           time.Time `orm:"null"` // 最近VIP获得时间
	VipExpire         time.Time `orm:"null"` // VIP过期时间
	Avatar            string
	Locked            int8
	LockedDate        time.Time `orm:"null"`
	LoginDate         time.Time `orm:"null"`
	LoginFailureCount int
	LoginIp           string
	Enabled           int8
	Token             string
	TokenExpTime      time.Time `orm:"null"`
	Cancelled         int8      // 是否注销 0-正常使用；1-已注销
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(Member))
}

/**
 * 获取格式化的手机号，格式：131*****234
 */
func (model *Member) GetFmtMobile() string {
	if len(model.Mobile) != 11 {
		return ""
	}
	return beego.Substr(model.Mobile, 0, 3) + "*****" + beego.Substr(model.Username, 8, 3)
}

/*
 * 获取完整的头像地址
 */
func (model *Member) GetFullAvatar(domain string) string {
	if model.Avatar == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(model.Avatar), "http") {
		return model.Avatar
	}
	if strings.HasPrefix(model.Avatar, "/") {
		return domain + model.Avatar
	}
	return domain + "/" + model.Avatar
}

func (model *Member) Paginate(page int, limit int, orderBy int, id int64, param1 string, regType int) (list []Member, total int64) {
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
	if regType != 0 {
		cond = cond.And("RegType", regType)
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
