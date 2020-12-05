package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type SiteConfig struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64                                         // 创建人Id
	Modifior   int64                                         // 更新人Id
	Version    int                                           // 版本
	Code       string    `orm:"unique;size(63)"`                      // 代码
	Value      string    `orm:"size(1024)"`                  // 值
	IsSystem   int8                                          // 是否内置(内置不可删除)
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(SiteConfig))
}

func GetSiteConfigValue(code string) string {
	var model SiteConfig
	o := orm.NewOrm()
	if err := o.QueryTable(new(SiteConfig)).Filter("Code", code).Limit(1).One(&model, "Value"); err != nil {
		beego.Error("SiteConfig get value error", err)
		return ""
	}
	return model.Value
}

func GetSiteConfigMap(code ...interface{}) map[string]string {
	m := make(map[string]string)
	if len(code) == 0 {
		return m
	}
	var list []SiteConfig
	o := orm.NewOrm()
	qs := o.QueryTable(new(SiteConfig))
	qs = qs.Filter("Code__in", code...)
	_, err := qs.All(&list, "Code", "Value")
	if err != nil {
		beego.Error("GetSiteConfigMap error", err)
		return m
	}
	for _, v := range list {
		m[v.Code] = v.Value
	}
	return m
}
