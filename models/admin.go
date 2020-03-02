package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const SysDbPrefix = "sys_"

type Admin struct {
	Id                int64     `auto`                              // 自增主键
	CreateDate        time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate        time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator           int64                                         // 创建人Id
	Modifior          int64                                         // 更新人Id
	Version           int                                           // 版本
	OrgId             int64                                         // 关联ID(公司id等)
	Enabled           int8                                          // 是否启用
	Locked            int8                                          // 是否锁定
	IsSystem          int8                                          // 是否系统内置
	LockedDate        time.Time `orm:"null"`                        // 锁定时间
	LoginDate         time.Time `orm:"null"`                        // 登录时间
	LoginFailureCount int                                           // 登录失败次数
	LoginIp           string    `orm:"null"`                        // 登录ip
	Salt              string                                        // 盐
	Name              string                                        // 名称
	Password          string                                        // 密码
	Username          string    `orm:"unique"`                      // 用户名
	Mobile            string    `orm:"null"`                        // 手机
	LoginVerify       int8                                          // 验证状态(0:无验证;1:短信验证；2：谷歌安全码验证;)
	GaSecret          string    `orm:"null"`                        // 谷歌验证秘钥
}

func init() {
	beego.Info("Init model admin")
	orm.RegisterModelWithPrefix(SysDbPrefix, new(Admin))
}

func (model *Admin) Paginate(page int, limit int, adminOrgId int64, param1 string) (list []Admin, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(Admin))
	cond := orm.NewCondition()
	if param1 != "" {
		cond = cond.AndCond(cond.And("Name__contains", param1).Or("Username__contains", param1))
	}
	cond = cond.And("IsSystem", 0)
	if adminOrgId != 0 {
		cond = cond.And("OrgId", adminOrgId)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("OrgId", "-Id")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
