package initial

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/models"
	. "github.com/iufansh/iufans/utils"
	"html/template"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	license "github.com/iufansh/iuplugins/iu-authorize"
	"net/http"
	"time"
)

func InitFilter() {
	beego.InsertFilter("*", beego.BeforeRouter, filterMethod)
	var adminRouter = beego.AppConfig.String("adminrouter")
	//beego.InsertFilter(adminRouter+"/login", beego.BeforeRouter, filterLicense)
	//beego.InsertFilter(adminRouter+"/sys/index", beego.BeforeRouter, filterLicense)
	//beego.InsertFilter(adminRouter+"/admin/*", beego.BeforeRouter, filterLicense)
	//beego.InsertFilter(adminRouter+"/site/*", beego.BeforeRouter, filterLicense)
	beego.InsertFilter(adminRouter+"/*", beego.BeforeRouter, filterAuth)
	beego.InsertFilter(adminRouter+"/*", beego.BeforeExec, filterBeforeExec)
	beego.InsertFilter(adminRouter+"/*", beego.BeforeExec, filterHttpGet)
}

var staticUrl = beego.AppConfig.String("staticurl")

var htmlHead = `
    <meta charset="UTF-8">
	<title></title>
	<meta name="renderer" content="webkit|ie-comp|ie-stand">
	<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
	<meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, user-scalable=0">
	<meta name="_xsrf" content="{{.xsrf_token}}" />
	<link rel="icon" href="data:image/ico;base64,aWNv">
    <link rel="shortcut icon" href="data:image/x-icon;," type="image/x-icon">
	<link rel="stylesheet" href="{{.static_url}}/static/layui/css/layui.css" media="all">
	<link rel="stylesheet" href="{{.static_url}}/static/back/css/common.css" media="all">
`
var scripts = `
<script src="{{.static_url}}/static/layui/layui.js"></script>
<script src="{{.static_url}}/static/back/js/admin.js?v=1.0"></script>
`

var filterMethod = func(ctx *context.Context) {
	if ctx.Input.Query("_method") != "" && ctx.Input.IsPost() {
		ctx.Request.Method = ctx.Input.Query("_method")
	}
}

var filterHttpGet = func(ctx *context.Context) {
	if ctx.Input.Method() == http.MethodGet {
		// 设置公共参数
		expire := int64(beego.BConfig.WebConfig.XSRFExpire)
		xsrfKey := ctx.XSRFToken(beego.BConfig.WebConfig.XSRFKey, expire)
		ctx.Input.SetData("xsrf_token", xsrfKey)
		ctx.Input.SetData("static_url", staticUrl)

		if t, err := template.New("HtmlHead.tpl").Parse(htmlHead); err != nil {
			logs.Error("filterAfterExec err1", err)
		} else {
			var buf bytes.Buffer
			if err := t.Execute(&buf, map[string]string{
				"xsrf_token": xsrfKey,
				"static_url": staticUrl,
			}); err != nil {

			}
			ctx.Input.SetData("HtmlHead", template.HTML(buf.String()))
		}
		if t, err := template.New("Scripts.tpl").Parse(scripts); err != nil {
			logs.Error("filterAfterExec err2", err)
		} else {
			var buf bytes.Buffer
			if err := t.Execute(&buf, map[string]string{
				"static_url": staticUrl,
			}); err != nil {

			}
			ctx.Input.SetData("Scripts", template.HTML(buf.String()))
		}
	}
}

var filterBeforeExec = func(ctx *context.Context) {
	// 不需要鉴权的url
	var currentUrl = ctx.Request.URL.Path
	switch currentUrl {
	case beego.URLFor("LoginController.Get"):
		return
	case beego.URLFor("LoginController.Logout"):
		return
	case beego.URLFor("LoginController.LoginVerify"):
		return
	case beego.URLFor("SysIndexController.GetAuth"):
		return
	case beego.URLFor("SysIndexController.PostAuth"):
		return
	case beego.URLFor("BaseIndexController.Get"):
		return
	case beego.URLFor("ChangePwdController.Get"):
		return
	case beego.URLFor("SysIndexController.Get"):
		return
	}
	o := orm.NewOrm()
	lid, _ := ctx.Input.Session("loginAdminId").(int64)
	// 谷歌验证
	forceGaAuth := beego.AppConfig.DefaultInt("forcegaauth", 0)
	if forceGaAuth == 1 {
		var admin models.Admin
		if err := o.QueryTable(new(models.Admin)).Filter("Id", lid).One(&admin, "LoginVerify"); err != nil {
			logs.Error("filterBeforeExec query admin err:", err)
			ctx.Abort(401, "内部错误，请刷新重试")
			return
		}
		if admin.LoginVerify == 0 {
			ctx.Redirect(http.StatusFound, beego.URLFor("SysIndexController.GetAuth"))
			return
		}
	}
	// 鉴权
	var arList orm.ParamsList
	_, err := o.QueryTable(new(models.AdminRole)).Filter("AdminId", lid).ValuesFlat(&arList, "RoleId")
	if err != nil {
		logs.Error("FilterAuth Query AdminRole error", err)
		ctx.Abort(500, "内部错误, 请联系管理员")
		return
	}
	if len(arList) == 0 {
		logs.Error("FilterAuth user no AdminRole, user id", lid)
		ctx.Abort(401, "没有权限1")
		return
	}
	_, err = o.QueryTable(new(models.Role)).Filter("Id__in", arList).Filter("Enabled", 1).ValuesFlat(&arList, "Id")
	if err != nil {
		logs.Error("FilterAuth Query AdminRole error", err)
		ctx.Abort(500, "内部错误, 请联系管理员")
		return
	}
	if len(arList) == 0 {
		logs.Error("FilterAuth user no Role, user id", lid)
		ctx.Abort(401, "没有权限2")
		return
	}
	var rpList orm.ParamsList
	_, err = o.QueryTable(new(models.RolePermission)).Filter("RoleId__in", arList).Distinct().ValuesFlat(&rpList, "PermissionId")
	if err != nil {
		logs.Error("FilterAuth Query RolePermission error", err)
		ctx.Abort(500, "内部错误, 请联系管理员")
		return
	}
	if len(rpList) == 0 {
		logs.Error("FilterAuth user no RolePermission, user id", lid)
		ctx.Abort(401, "没有权限3")
		return
	}
	var permList orm.ParamsList
	_, err = o.QueryTable(new(models.Permission)).Filter("Id__in", rpList).Filter("Enabled", 1).ValuesFlat(&permList, "Url")
	if err != nil {
		logs.Error("FilterAuth Query Permission error", err)
		ctx.Abort(500, "内部错误, 请联系管理员")
		return
	}
	ps := ctx.Input.Params()
	var urlArgs = make([]interface{}, 0)
	for k, v := range ps {
		if k != ":splat" && strings.HasPrefix(k, ":") {
			urlArgs = append(urlArgs, k, v)
		}
	}
	for _, perm := range permList {
		if perm != nil && perm.(string) != "" && beego.URLFor(perm.(string), urlArgs...) == ctx.Request.URL.Path {
			return
		}
	}
	// 没有权限
	if _, err := ctx.ResponseWriter.Write([]byte("没有权限或页面不存在")); err != nil {

	}
	ctx.Abort(401, "没有权限或页面不存在")
}

/**
 * 登录验证
 */
var filterAuth = func(ctx *context.Context) {
	// 不需要鉴权的url
	var currentUrl = ctx.Request.URL.Path
	switch currentUrl {
	case beego.URLFor("LoginController.Get"):
		return
	case beego.URLFor("LoginController.Logout"):
		return
	case beego.URLFor("LoginController.LoginVerify"):
		return
	}
	// 登录验证
	lid, ok := ctx.Input.Session("loginAdminId").(int64)
	if !ok {
		ctx.Redirect(302, beego.URLFor("LoginController.Get"))
	}
	// token验证
	sesToken, ok := ctx.Input.Session("token").(string)
	var cacToken string
	if err := GetCache(fmt.Sprintf("loginAdminId%d", lid), &cacToken); err != nil {
		if _, err := ctx.ResponseWriter.Write([]byte("Token获取失败")); err != nil {

		}
		ctx.Abort(401, "Token获取失败")
	}
	if !ok || sesToken == "" || sesToken != cacToken {
		if _, err := ctx.ResponseWriter.Write([]byte("登录过期，请重新登录")); err != nil {

		}
		ctx.Abort(401, "登录过期，请重新登录")
	}
}

var filterLicense = func(ctx *context.Context) {
	// 不需要登录的url
	switch ctx.Request.RequestURI {
	case beego.URLFor("SysIndexController.Systeminfo"):
		return
	}
	lic := beego.AppConfig.String("serverlicense")
	if lic == "" {
		if _, err := ctx.ResponseWriter.Write([]byte("当前系统为试用版，请购买正版")); err != nil {

		}
		ctx.Abort(500, "当前系统为试用版，请购买正版")
		logs.Error("License not found, please config!")
		return
	}

	payTime, err := time.ParseInLocation("20060102150405", lic[:14], time.Local)
	if err != nil {
		if _, err := ctx.ResponseWriter.Write([]byte("注册码日期异常")); err != nil {

		}
		ctx.Abort(500, "注册码日期异常")
		logs.Error("License exp time err, please check, format is 20060102150405!")
		return
	}
	ok, msg := license.CheckLicense(lic[14:], payTime, false, "")
	if ok {
		return
	}
	logs.Error(msg)

	if _, err := ctx.ResponseWriter.Write([]byte("当前系统为试用版，请购买正版")); err != nil {

	}
	ctx.Abort(500, "当前系统为试用版，请购买正版")
	return
}
