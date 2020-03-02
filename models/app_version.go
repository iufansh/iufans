package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type AppVersion struct {
	Id           int64     `auto`                              // 自增主键
	CreateDate   time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate   time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator      int64     // 创建人Id
	Modifior     int64     // 更新人Id
	Version      int       // 版本
	OrgId        int64     // 关联ID(公司id等)
	AppNo        string    // App编号，每个App编号唯一
	OsType       string    // android,ios
	VersionNo    int       // 版本号
	VersionName  string    // 版本名称
	VersionDesc  string    // 版本描述
	DownloadUrl  string    // 下载地址
	PublishTime  time.Time // 发布时间
	ForceUpdate  int8      // 是否强制升级
	Ignorable    int8      // 是否可忽略
	AppSize      int       // app大小(单位：KB)
	EncryptValue string    // 加密值(默认MD5)
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(AppVersion))
}

func (model *AppVersion) Paginate(page int, limit int, orgId int64, param1 string, osType string, versionNo int) (list []AppVersion, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(AppVersion))
	cond := orm.NewCondition()
	if param1 != "" {
		cond = cond.AndCond(cond.And("VersionName__contains", param1).Or("VersionDesc__contains", param1).Or("AppNo__contains", param1))
	}
	cond = cond.And("OrgId", orgId)
	if osType != "" {
		cond = cond.And("OsType", osType)
	}
	if versionNo != 0 {
		cond = cond.And("VersionNo", versionNo)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Id")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
