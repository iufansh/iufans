package initial

import (
	"fmt"
	. "github.com/iufansh/iufans/models"
	. "github.com/iufansh/iufans/utils"

	"bytes"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	license "github.com/iufansh/iuplugins/iu-authorize"
	"html/template"
	"net/http"
	"time"
)

func InitFilter() {
	var adminRouter = beego.AppConfig.String("adminrouter")
	//beego.InsertFilter(adminRouter+"/login", beego.BeforeRouter, filterLicense)
	//beego.InsertFilter(adminRouter+"/sys/index", beego.BeforeRouter, filterLicense)
	//beego.InsertFilter(adminRouter+"/admin/*", beego.BeforeRouter, filterLicense)
	//beego.InsertFilter(adminRouter+"/site/*", beego.BeforeRouter, filterLicense)
	beego.InsertFilter(adminRouter+"/*", beego.BeforeRouter, filterAuth)
	beego.InsertFilter(adminRouter+"/*", beego.BeforeExec, filterBeforeExec)
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

var filterBeforeExec = func(ctx *context.Context) {
	expire := int64(beego.BConfig.WebConfig.XSRFExpire)
	xsrfKey := ctx.XSRFToken(beego.BConfig.WebConfig.XSRFKey, expire)
	ctx.Input.SetData("xsrf_token", xsrfKey)
	ctx.Input.SetData("static_url", staticUrl)

	if ctx.Input.Method() == http.MethodGet {
		if t, err := template.New("HtmlHead.tpl").Parse(htmlHead); err != nil {
			beego.Error("filterAfterExec err1", err)
		} else {
			var buf bytes.Buffer
			t.Execute(&buf, map[string]string{
				"xsrf_token": xsrfKey,
				"static_url": staticUrl,
			})
			ctx.Input.SetData("HtmlHead", template.HTML(buf.String()))
		}
		if t, err := template.New("Scripts.tpl").Parse(scripts); err != nil {
			beego.Error("filterAfterExec err2", err)
		} else {
			var buf bytes.Buffer
			t.Execute(&buf, map[string]string{
				"static_url": staticUrl,
			})
			ctx.Input.SetData("Scripts", template.HTML(buf.String()))
		}
	}
}

/**
 * 登录验证、鉴权
 */
var filterAuth = func(ctx *context.Context) {
	// 不需要鉴权的url
	switch ctx.Request.RequestURI {
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
	sestoken, ok := ctx.Input.Session("token").(string)
	var cactoken string
	GetCache(fmt.Sprintf("loginAdminId%d", lid), &cactoken)
	if !ok || sestoken == "" || sestoken != cactoken {
		ctx.ResponseWriter.Write([]byte("登录过期，请重新登录"))
		ctx.Abort(401, "登录过期，请重新登录")
	}

	// 鉴权
	o := orm.NewOrm()
	var arList orm.ParamsList
	_, err := o.QueryTable(new(AdminRole)).Filter("AdminId", lid).ValuesFlat(&arList, "RoleId")
	if err != nil {
		beego.Error("FilterAuth Query AdminRole error", err)
		ctx.Abort(500, "内部错误, 请联系管理员")
		return
	}
	if len(arList) == 0 {
		beego.Error("FilterAuth user no AdminRole, user id", lid)
		ctx.Abort(401, "没有权限1")
		return
	}
	_, err = o.QueryTable(new(Role)).Filter("Id__in", arList).Filter("Enabled", 1).ValuesFlat(&arList, "Id")
	if err != nil {
		beego.Error("FilterAuth Query AdminRole error", err)
		ctx.Abort(500, "内部错误, 请联系管理员")
		return
	}
	if len(arList) == 0 {
		beego.Error("FilterAuth user no Role, user id", lid)
		ctx.Abort(401, "没有权限2")
		return
	}
	var rpList orm.ParamsList
	_, err = o.QueryTable(new(RolePermission)).Filter("RoleId__in", arList).Distinct().ValuesFlat(&rpList, "PermissionId")
	if err != nil {
		beego.Error("FilterAuth Query RolePermission error", err)
		ctx.Abort(500, "内部错误, 请联系管理员")
		return
	}
	if len(rpList) == 0 {
		beego.Error("FilterAuth user no RolePermission, user id", lid)
		ctx.Abort(401, "没有权限3")
		return
	}
	var permList orm.ParamsList
	_, err = o.QueryTable(new(Permission)).Filter("Id__in", rpList).Filter("Enabled", 1).ValuesFlat(&permList, "Url")
	if err != nil {
		beego.Error("FilterAuth Query Permission error", err)
		ctx.Abort(500, "内部错误, 请联系管理员")
		return
	}
	var currentUrl = ctx.Request.URL.EscapedPath()
	var isAuth = false
	for _, perm := range permList {
		if perm != nil && perm.(string) != "" && beego.URLFor(perm.(string)) == currentUrl {
			isAuth = true
		}
	}
	// 没有权限
	if !isAuth {
		ctx.ResponseWriter.Write([]byte("没有权限或页面不存在"))
		ctx.Abort(401, "没有权限或页面不存在")
	}
}

var filterLicense = func(ctx *context.Context) {
	beego.Info("filter license")
	// 不需要登录的url
	switch ctx.Request.RequestURI {
	case beego.URLFor("SysIndexController.Systeminfo"):
		return
	}
	lic := beego.AppConfig.String("serverlicense")
	if lic == "" {
		ctx.ResponseWriter.Write([]byte("当前系统为试用版，请购买正版"))
		ctx.Abort(500, "当前系统为试用版，请购买正版")
		beego.Error("License not found, please config!")
		return
	}

	payTime, err := time.ParseInLocation("20060102150405", lic[:14], time.Local)
	if err != nil {
		ctx.ResponseWriter.Write([]byte("注册码日期异常"))
		ctx.Abort(500, "注册码日期异常")
		beego.Error("License exp time err, please check, format is 20060102150405!")
		return
	}
	ok, msg := license.CheckLicense(lic[14:], payTime, false, "")
	if ok {
		return
	}
	beego.Error(msg)

	ctx.ResponseWriter.Write([]byte("当前系统为试用版，请购买正版"))
	ctx.Abort(500, "当前系统为试用版，请购买正版")
	return
}
