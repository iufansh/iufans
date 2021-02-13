package appbanner

import (
	"github.com/astaxie/beego/logs"
	"github.com/iufansh/iufans/controllers/sysmanage"
	"github.com/astaxie/beego"
	. "github.com/iufansh/iufans/models"
	"html/template"
	"github.com/astaxie/beego/orm"
)

type AppBannerIndexController struct {
	sysmanage.BaseController
}

func (c *AppBannerIndexController) NestPrepare()  {
	c.EnableRender = false
}

func (c *AppBannerIndexController) Get() {
	param1 := c.GetString("param1")
	status, _ := c.GetInt("status", -1)
	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(AppBanner).Paginate(page, limit, c.LoginAdminOrgId, param1, status)
	c.SetPaginator(limit, total)
	c.Data["dataList"] = &list

	c.Data["condArr"] = map[string]interface{}{"param1": param1, "status": status}

	c.Data["urlAppBannerIndexGet"] = c.URLFor("AppBannerIndexController.Get")
	c.Data["urlAppBannerIndexDelone"] = c.URLFor("AppBannerIndexController.Delone")
	c.Data["urlAppBannerAddGet"] = c.URLFor("AppBannerAddController.Get")
	c.Data["urlAppBannerEditGet"] = c.URLFor("AppBannerEditController.Get")

	if t, err := template.New("tplAppBannerIndex.tpl").Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *AppBannerIndexController) Delone() {
	var code int
	var msg string
	url := beego.URLFor("AppBannerIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := AppBanner{}
	model.Id = id
	err := o.Read(&model)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	}
	_, err1 := o.Delete(&model, "Id")
	if err1 != nil {
		logs.Error("Delete AppBanner eror", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type AppBannerAddController struct {
	sysmanage.BaseController
}

func (c *AppBannerAddController) NestPrepare()  {
	c.EnableRender = false
}

func (c *AppBannerAddController) Get() {
	c.Data["urlSyscommonUpload"] = c.URLFor("SyscommonController.Upload")
	c.Data["urlAppBannerIndexGet"] = c.URLFor("AppBannerIndexController.Get")
	c.Data["urlAppBannerAddPost"] = c.URLFor("AppBannerAddController.Post")

	if t, err := template.New("tplAddAppBanner.tpl").Parse(tplAdd); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *AppBannerAddController) Post() {
	var code int
	var msg string
	var url string
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	AppBanner := AppBanner{}
	if err := c.ParseForm(&AppBanner); err != nil {
		msg = "参数异常"
		return
	}
	AppBanner.Creator = c.LoginAdminId
	AppBanner.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Insert(&AppBanner); err != nil {
		msg = "添加失败"
		logs.Error("添加失败", err)
	} else {
		code = 1
		msg = "添加成功"
	}
}

type AppBannerEditController struct {
	sysmanage.BaseController
}

func (c *AppBannerEditController) NestPrepare()  {
	c.EnableRender = false
}

func (c *AppBannerEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := AppBanner{}
	model.Id = id
	err := o.Read(&model)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(beego.URLFor("AppBannerIndexController.Get"), 302)
	}
	c.Data["data"] = &model

	c.Data["urlSyscommonUpload"] = c.URLFor("SyscommonController.Upload")
	c.Data["urlAppBannerIndexGet"] = c.URLFor("AppBannerIndexController.Get")
	c.Data["urlAppBannerEditPost"] = c.URLFor("AppBannerEditController.Post")

	if t, err := template.New("tplEditAppBanner.tpl").Parse(tplEdit); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *AppBannerEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("AppBannerIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	AppBanner := AppBanner{}
	if err := c.ParseForm(&AppBanner); err != nil {
		msg = "参数异常"
		return
	}
	cols := []string{"AppNo", "Seq", "Title", "Banner", "JumpUrl", "Status", "Modifior", "ModifyDate"}
	AppBanner.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Update(&AppBanner, cols...); err != nil {
		msg = "更新失败"
		logs.Error("更新失败", err)
	} else {
		code = 1
		msg = "更新成功"
	}
}
