package models

import (
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type AppBanner struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64     // 创建人Id
	Modifior   int64     // 更新人Id
	Version    int       // 版本
	OrgId      int64     // 关联ID(公司id等)
	AppNo      string    // App编号
	Seq        int       // 排序(升序显示)
	Title      string    // 标题
	Banner     string    // banner
	JumpUrl    string    // 跳转地址，网址或者App内部页
	Status     int       // 状态0-禁用；1-启用
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(AppBanner))
}

/*
 * 获取完整的banner地址
 */
func (model *AppBanner) GetFullBanner(domain string) string {
	if model.Banner == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(model.Banner), "http") {
		return model.Banner
	}
	if strings.HasPrefix(model.Banner, "/") {
		return domain + model.Banner
	}
	return domain + "/" + model.Banner
}

func (model *AppBanner) Paginate(page int, limit int, orgId int64, param1 string, status int) (list []AppBanner, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(AppBanner))
	cond := orm.NewCondition()
	cond = cond.And("OrgId", orgId)
	if param1 != "" {
		cond = cond.And("AppNo", param1)
	}
	if status != -1 {
		cond = cond.And("Status", status)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("Seq", "-Id")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
