package quicknav

import (
	"github.com/iufansh/iufans/controllers/sysmanage"
	"github.com/astaxie/beego"
	. "github.com/iufansh/iufans/models"
	"html/template"
	"github.com/astaxie/beego/orm"
)

type QuickNavIndexController struct {
	sysmanage.BaseController
}

func (c *QuickNavIndexController) NestPrepare()  {
	c.EnableRender = false
}

func (c *QuickNavIndexController) Get() {
	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(QuickNav).Paginate(page, limit)
	c.SetPaginator(limit, total)
	c.Data["dataList"] = &list

	c.Data["urlQuickNavIndexDelone"] = c.URLFor("QuickNavIndexController.Delone")
	c.Data["urlQuickNavAddGet"] = c.URLFor("QuickNavAddController.Get")
	c.Data["urlQuickNavEditGet"] = c.URLFor("QuickNavEditController.Get")

	if t, err := template.New("tplQuickNavIndex.tpl").Parse(tplIndex); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *QuickNavIndexController) Delone() {
	var code int
	var msg string
	url := beego.URLFor("QuickNavIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := QuickNav{}
	model.Id = id
	err := o.Read(&model)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	}
	_, err1 := o.Delete(&model, "Id")
	if err1 != nil {
		beego.Error("Delete QuickNav eror", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type QuickNavAddController struct {
	sysmanage.BaseController
}

func (c *QuickNavAddController) NestPrepare()  {
	c.EnableRender = false
}

func (c *QuickNavAddController) Get() {
	c.Data["urlQuickNavIndexGet"] = c.URLFor("QuickNavIndexController.Get")
	c.Data["urlQuickNavAddPost"] = c.URLFor("QuickNavAddController.Post")

	if t, err := template.New("tplAddQuickNav.tpl").Parse(tplAdd); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *QuickNavAddController) Post() {
	var code int
	var msg string
	var url = beego.URLFor("QuickNavIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	quicknav := QuickNav{}
	if err := c.ParseForm(&quicknav); err != nil {
		msg = "参数异常"
		return
	}
	quicknav.Creator = c.LoginAdminId
	quicknav.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Insert(&quicknav); err != nil {
		msg = "添加失败"
		beego.Error("添加失败", err)
	} else {
		code = 1
		msg = "添加成功"
	}
}

type QuickNavEditController struct {
	sysmanage.BaseController
}

func (c *QuickNavEditController) NestPrepare()  {
	c.EnableRender = false
}

func (c *QuickNavEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := QuickNav{}
	model.Id = id
	err := o.Read(&model)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(beego.URLFor("QuickNavIndexController.Get"), 302)
	}
	c.Data["data"] = &model

	c.Data["urlQuickNavIndexGet"] = c.URLFor("QuickNavIndexController.Get")
	c.Data["urlQuickNavEditPost"] = c.URLFor("QuickNavEditController.Post")

	if t, err := template.New("tplEditQuickNav.tpl").Parse(tplEdit); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *QuickNavEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("QuickNavIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	quicknav := QuickNav{}
	if err := c.ParseForm(&quicknav); err != nil {
		msg = "参数异常"
		return
	}
	cols := []string{"Name", "WebSite", "Icon", "Seq", "ModifyDate"}
	quicknav.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Update(&quicknav, cols...); err != nil {
		msg = "更新失败"
		beego.Error("更新失败", err)
	} else {
		code = 1
		msg = "更新成功"
	}
}
