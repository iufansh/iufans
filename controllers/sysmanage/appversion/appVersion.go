package appversion

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/controllers/sysmanage"
	. "github.com/iufansh/iufans/models"
	"html/template"
	"strings"
)

type AppVersionIndexController struct {
	sysmanage.BaseController
}

func (c *AppVersionIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *AppVersionIndexController) Get() {

	param1 := c.GetString("param1")
	osType := c.GetString("osType")
	versionNo, _ := c.GetInt("versionNo", 0)
	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(AppVersion).Paginate(page, limit, c.LoginAdminOrgId, param1, osType, versionNo)
	c.SetPaginator(limit, total)
	c.Data["dataList"] = &list

	c.Data["condArr"] = map[string]interface{}{"param1": param1, "osType": osType, "versionNo": versionNo}

	c.Data["urlAppVersionIndexDelone"] = c.URLFor("AppVersionIndexController.Delone")
	c.Data["urlAppVersionAddGet"] = c.URLFor("AppVersionAddController.Get")
	c.Data["urlAppVersionEditGet"] = c.URLFor("AppVersionEditController.Get")

	if t, err := template.New("tplAppVersionIndex.tpl").Funcs(map[string]interface{}{
		"date": beego.Date,
	}).Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *AppVersionIndexController) Delone() {
	var code int
	var msg string
	url := beego.URLFor("AppVersionIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := AppVersion{}
	model.Id = id
	err := o.Read(&model)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	}
	_, err1 := o.Delete(&model, "Id")
	if err1 != nil {
		logs.Error("Delete AppVersion eror", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type AppVersionAddController struct {
	sysmanage.BaseController
}

func (c *AppVersionAddController) NestPrepare() {
	c.EnableRender = false
}

func (c *AppVersionAddController) Get() {
	c.Data["urlAppVersionIndexGet"] = c.URLFor("AppVersionIndexController.Get")
	c.Data["urlAppVersionAddPost"] = c.URLFor("AppVersionAddController.Post")
	c.Data["appChannel"] = ""
	if t, err := template.New("tplAddAppVersion.tpl").Funcs(map[string]interface{}{
		"urlfor": beego.URLFor,
	}).Parse(tplAdd); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *AppVersionAddController) Post() {
	var code int
	var msg string
	var url = beego.URLFor("AppVersionIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	model := AppVersion{}
	if err := c.ParseForm(&model); err != nil {
		msg = "参数异常"
		return
	}
	appChannel := c.GetString("appChannel")
	if appChannel != "" {
		model.AppNo = model.AppNo + "-" + appChannel
	}
	model.OrgId = c.LoginAdminOrgId
	model.Creator = c.LoginAdminId
	model.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Insert(&model); err != nil {
		msg = "添加失败"
		logs.Error("添加失败", err)
	} else {
		code = 1
		msg = "添加成功"
	}
}

type AppVersionEditController struct {
	sysmanage.BaseController
}

func (c *AppVersionEditController) NestPrepare() {
	c.EnableRender = false
}

func (c *AppVersionEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := AppVersion{}
	model.Id = id
	err := o.Read(&model)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(beego.URLFor("AppVersionIndexController.Get"), 302)
	}
	if strings.Contains(model.AppNo, "-") {
		c.Data["appChannel"] = strings.Split(model.AppNo, "-")[1]
		model.AppNo = strings.Split(model.AppNo, "-")[0]
	} else {
		c.Data["appChannel"] = ""
	}
	c.Data["data"] = &model

	c.Data["urlAppVersionIndexGet"] = c.URLFor("AppVersionIndexController.Get")
	c.Data["urlAppVersionEditPost"] = c.URLFor("AppVersionEditController.Post")

	if t, err := template.New("tplEditAppVersion.tpl").Funcs(map[string]interface{}{
		"date":   beego.Date,
		"urlfor": beego.URLFor,
	}).Parse(tplEdit); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *AppVersionEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("AppVersionIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	model := AppVersion{}
	if err := c.ParseForm(&model); err != nil {
		msg = "参数异常"
		return
	}
	appChannel := c.GetString("appChannel")
	if appChannel != "" {
		model.AppNo = model.AppNo + "-" + appChannel
	}
	cols := []string{"AppNo", "OsType", "VersionNo", "VersionName", "VersionDesc", "DownloadUrl", "AppSize", "EncryptValue", "PublishTime", "ForceUpdate", "Ignorable"}
	model.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Update(&model, cols...); err != nil {
		msg = "更新失败"
		logs.Error("更新失败", err)
	} else {
		code = 1
		msg = "更新成功"
	}
}
