package memberlogincount

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/controllers/sysmanage"
	models2 "github.com/iufansh/iufans/models"
	"html/template"
	"strings"
	"time"
)

type MemberLoginCountIndexController struct {
	sysmanage.BaseController
}

type memberResult struct {
	AppNo         string
	AppChannel    string
	AppVersion    int
	CountDate    time.Time
	Count      int
}

// 用户分析
func (c *MemberLoginCountIndexController) Get() {
	appNo := c.GetString("appNo")
	appChannel := c.GetString("appChannel")
	appVersion, _ := c.GetInt("appVersion", 0)
	timeStart := strings.TrimSpace(c.GetString("timeStart"))
	timeEnd := strings.TrimSpace(c.GetString("timeEnd"))

	group_appNo := c.GetString("group_appNo")
	group_appChannel := c.GetString("group_appChannel")
	group_appVersion := c.GetString("group_appVersion")
	group_countDate := c.GetString("group_countDate")

	if timeStart == "" {
		timeStart = time.Now().Format("2006-01-02") + " 00:00:00"
	}

	qb, _ := orm.NewQueryBuilder(beego.AppConfig.String("dbdriver"))

	fields := []string{
		"sum(m.Count) as count",
	}
	groupItems := make([]string, 0)
	if group_appNo == "on" {
		fields = append(fields, "m.app_no")
		groupItems = append(groupItems, "m.app_no")
	}
	if group_appChannel == "on" {
		fields = append(fields, "m.app_channel")
		groupItems = append(groupItems, "m.app_channel")
	}
	if group_appVersion == "on" {
		fields = append(fields, "m.app_version")
		groupItems = append(groupItems, "m.app_version")
	}
	if group_countDate == "on" {
		fields = append(fields, "cast(m.count_date AS date) AS count_date")
		groupItems = append(groupItems, "cast(m.count_date AS date)")
	}
	args := []interface{}{timeStart}
	// 构建查询对象
	qb.Select(fields...).
		From(models2.SysDbPrefix + "member_login_count m").
		Where("m.count_date >= ?")
	if timeEnd != "" {
		qb.And("m.count_date <= ?")
		args = append(args, timeEnd)
	}
	if appNo != "" {
		qb.And("m.app_no = ?")
		args = append(args, appNo)
	}
	if appChannel != "" {
		qb.And("m.app_channel = ?")
		args = append(args, appChannel)
	}
	if appVersion != 0 {
		qb.And("m.app_version = ?")
		args = append(args, appVersion)
	}
	if len(groupItems) > 0 {
		qb.GroupBy(groupItems...)
	}
	qb.OrderBy("m.app_no", "m.app_channel", "m.app_version", "m.count_date")

	// 导出 SQL 语句
	sql := qb.String()

	// 执行 SQL 语句
	o := orm.NewOrm()
	var list []memberResult
	if _, err := o.Raw(sql, args...).QueryRows(&list); err != nil {
		logs.Error("MemberLoginCountIndexController Raw sql err,", err)
		c.Msg = "查询异常"
	}

	c.Data["condArr"] = map[string]interface{}{
		"appNo":            appNo,
		"appChannel":       appChannel,
		"appVersion":       appVersion,
		"timeStart":        timeStart,
		"timeEnd":          timeEnd,
		"group_appNo":      group_appNo,
		"group_appChannel": group_appChannel,
		"group_appVersion": group_appVersion,
		"group_countDate": group_countDate}
	c.Data["dataList"] = &list

	c.Data["urlMemberLoginCountIndexGet"] = c.URLFor("MemberLoginCountIndexController.Get")

	if t, err := template.New("tplIndexMemberLoginCount.tpl").Funcs(map[string]interface{}{
		"date": beego.Date,
	}).Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}
