package sysmanage

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	. "github.com/iufansh/iufans/models"
	"github.com/astaxie/beego/orm"
	pagination2 "github.com/astaxie/beego/utils/pagination"
	"github.com/iufansh/iufans/utils"
	"html/template"
	"bytes"
	"github.com/iufansh/iufans/controllers"
	"net/http"
)

type NestPreparer interface {
	NestPrepare()
}

type BaseController struct {
	controllers.BaseMainController
	LoginAdminId       int64
	LoginAdminUsername string
	LoginAdminName     string
	LoginAdminOrgId    int64
}

func (c *BaseController) Prepare() {

	beego.Info("\r\n----------request---------",
		"\r\nUri:", c.Ctx.Input.URI(),
		"\r\nMethod:", c.Ctx.Input.Method(),
		"\r\nFrom ip:", c.Ctx.Input.IP(),
		"\r\nUserAgent:", c.Ctx.Input.UserAgent(),
		"\r\nBody:", string(c.Ctx.Input.RequestBody),
		"\r\n--------------------------")
	adminId, ok := c.GetSession("loginAdminId").(int64)
	if !ok {
		if c.Ctx.Input.Method() == http.MethodGet {
			c.Abort("请先登录")
		} else {
			c.Data["json"] = map[string]interface{}{"msg": "请先登录"}
			c.ServeJSON()
		}
		return
	}
	c.LoginAdminId = adminId
	c.LoginAdminUsername = c.GetSession("loginAdminUsername").(string)
	c.LoginAdminName = c.GetSession("loginAdminName").(string)
	c.LoginAdminOrgId, _ = c.GetSession("loginAdminOrgId").(int64)
	if app, ok := c.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

// Deprecated TODO 改方法废弃
func Retjson(ctx *context.Context, msg *string, code *int, data ...interface{}) {
	ret := make(map[string]interface{})
	ret["code"] = code
	ret["msg"] = msg
	if len(data) > 0 {
		d := data[0]
		switch d.(type) {
		case string:
			ret["url"] = d
			break
		case *string:
			ret["url"] = d
			break
		}
		ret["data"] = d
	}
	ctx.Output.JSON(ret, false, false)
}

var pagination = `
<div>
	<ul class="pagination">
		<li class="disabled"><a>总记录数：{{.paginator.Nums}} 条</a></li>
	{{if .paginator.HasPrev}}
	    <li><a href="{{.paginator.PageLinkFirst}}">第一页</a></li>
	    <li><a href="{{.paginator.PageLinkPrev}}">&laquo;</a></li>
	{{else}}
	    <li class="disabled"><a>第一页</a></li>
	    <li class="disabled"><a>&laquo;</a></li>
	{{end}}
	{{range $index, $page := .paginator.Pages}}
	    <li{{if $.paginator.IsActive .}} class="active"{{end}}>
	        <a href="{{$.paginator.PageLink $page}}">{{$page}}</a>
	    </li>
	{{end}}
	{{if .paginator.HasNext}}
	    <li><a href="{{.paginator.PageLinkNext}}">&raquo;</a></li>
	    <li><a href="{{.paginator.PageLinkLast}}">最后一页</a></li>
	{{else}}
	    <li class="disabled"><a>&raquo;</a></li>
	    <li class="disabled"><a>最后一页</a></li>
	{{end}}
	</ul>
</div>
`

func (c *BaseController) SetPaginator(per int, nums int64) {
	paginator := pagination2.NewPaginator(c.Ctx.Request, per, nums)
	if t, err := template.New("Pagination.tpl").Parse(pagination); err != nil {
		beego.Error("filterAfterExec err3", err)
	} else {
		var buf bytes.Buffer
		t.Execute(&buf, map[string]interface{}{
			"paginator": paginator,
		})
		c.Data["Pagination"] = template.HTML(buf.String())
	}
}

type BaseIndexController struct {
	BaseController
}

func (c *BaseIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *BaseIndexController) Get() {
	// 获取左侧菜单
	o := orm.NewOrm()
	sql := "select * from "+SysDbPrefix+"permission a where a.enabled = 1 and display = 1 and exists(select b.id from "+SysDbPrefix+"role_permission b, "+SysDbPrefix+"admin_role c where b.role_id = c.role_id and b.permission_id = a.id and c.admin_id = ?) order by a.pid, a.sort, a.id"
	var permissions []Permission
	_, err := o.Raw(sql, c.LoginAdminId).QueryRows(&permissions)
	if err != nil {
		beego.Error("Query admin permission error", err)
		c.Abort("内部错误，请重试")
	} else {
		var mainMenuList []Permission
		secdMenuMap := make(map[int64][]Permission)
		for _, pe := range permissions {
			// 构建菜单
			if pe.Pid == 0 {
				mainMenuList = append(mainMenuList, pe)
			} else {
				if val, ok := secdMenuMap[pe.Pid]; ok {
					val = append(val, pe)
					secdMenuMap[pe.Pid] = val
				} else {
					var menuList []Permission
					menuList = append(menuList, pe)
					secdMenuMap[pe.Pid] = menuList
				}
			}
		}
		c.Data["loginAdminName"] = c.GetSession("loginAdminName")
		c.Data["mainMenuList"] = mainMenuList
		c.Data["secdMenuMap"] = secdMenuMap
	}
	c.Data["static_url"] = beego.AppConfig.String("staticurl")
	c.Data["siteName"] = GetSiteConfigValue(utils.Scname)

	if t, err := template.New("tplBaseIndex.tpl").Funcs(map[string]interface{}{
		"date":   beego.Date,
		"urlfor": beego.URLFor,
		"substr": beego.Substr,
	}).Parse(tplBase); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}
