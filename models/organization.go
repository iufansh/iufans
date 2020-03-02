package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Organization struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64                                         // 创建人Id
	Modifior   int64                                         // 更新人Id
	Version    int                                           // 版本
	OrgId      int64                                         // 上级ID
	Levels     string                                        // 层级关系
	Vip        int                                           // VIP等级
	Name       string                                        // 名称
	OrgType    int                                           // 组织类型，用于业务自定义
	BindDomain string                                        // 绑定域名
	ExpireTime time.Time                                     // 过期时间
	EncryptKey string                                        // 秘钥(二级密码)
	Enabled    int8                                          // 状态 0：禁用；1：启用
	Remark     string                                        // 备注
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(Organization))
}

func (model *Organization) Paginate(page int, limit int, orgId int64, param1 string, orgTypes []int) (list []Organization, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(Organization))
	cond := orm.NewCondition()
	if param1 != "" {
		cond = cond.AndCond(cond.And("Name__contains", param1).Or("Remark__contains", param1))
	}
	if orgId >= 0 {
		cond = cond.And("OrgId", orgId)
	}
	if len(orgTypes) > 0 {
		cond = cond.And("OrgType__in", orgTypes)
	}
	qs = qs.SetCond(cond)
	qs = qs.OrderBy("-Id")
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs.All(&list)
	total, _ = qs.Count()
	return
}
