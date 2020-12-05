package models

import (
	"github.com/astaxie/beego/logs"
	"time"

	"github.com/astaxie/beego/orm"
)

type SmsLog struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64     // 创建人Id
	Modifior   int64     // 更新人Id
	Version    int       // 版本
	OrgId      int64     // 关联ID(公司id等)
	AppInfo    string    // App信息
	Ip         string    // 发送ip
	Receiver   string    // 接收人
	Info       string    // 短信内容
	Status     int       // 0-草稿；1-发送中；2-发送成功；3-发送失败
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(SmsLog))
}

func conditionSet(orgId int64, param1 string, status int, timeStart string, timeEnd string) orm.QuerySeter {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SmsLog))
	cond := orm.NewCondition()
	cond = cond.And("OrgId", orgId)
	if param1 != "" {
		cond = cond.AndCond(cond.And("Ip__contains", param1).Or("Receiver__contains", param1).Or("Info__contains", param1))
	}
	if status != -1 {
		cond = cond.And("Status", status)
	}
	if timeStart != "" {
		cond = cond.And("CreateDate__gte", timeStart)
	}
	if timeEnd != "" {
		cond = cond.And("CreateDate__lte", timeEnd)
	}
	return qs.SetCond(cond)
}

func (model *SmsLog) Paginate(page int, limit int, orgId int64, param1 string, status int, timeStart string, timeEnd string) (list []SmsLog, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	qs := conditionSet(orgId, param1, status, timeStart, timeEnd)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Id")
	qs.All(&list)
	total, _ = qs.Count()
	return
}

// 条件必须和Paginate一样
func (model *SmsLog) Del(orgId int64, param1 string, status int, timeStart string, timeEnd string) int64 {
	qs := conditionSet(orgId, param1, status, timeStart, timeEnd)
	if num, err := qs.Delete(); err != nil {
		logs.Error("SmsLog Del err:", err)
		return 0
	} else {
		return num
	}
}

func (model *SmsLog) InsertLog() error {
	o := orm.NewOrm()
	if _, err := o.Insert(model); err != nil {
		return err
	}
	return nil
}
