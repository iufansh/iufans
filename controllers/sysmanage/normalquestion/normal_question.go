package normalquestion

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/controllers/sysmanage"
	. "github.com/iufansh/iufans/models"
	"html/template"
	"strings"
)

type NormalQuestionIndexController struct {
	sysmanage.BaseController
}

func (c *NormalQuestionIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *NormalQuestionIndexController) Get() {
	param1 := strings.TrimSpace(c.GetString("param1"))
	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(NormalQuestion).Paginate(page, limit, param1)
	c.SetPaginator(limit, total)
	c.Data["dataList"] = &list
	c.Data["condArr"] = map[string]interface{}{"param1": param1}
	c.Data["urlNormalQuestionIndexDelone"] = c.URLFor("NormalQuestionIndexController.Delone")
	c.Data["urlNormalQuestionAddGet"] = c.URLFor("NormalQuestionAddController.Get")
	c.Data["urlNormalQuestionEditGet"] = c.URLFor("NormalQuestionEditController.Get")

	if t, err := template.New("tplNormalQuestionIndex.tpl").Funcs(map[string]interface{}{
		"date": beego.Date,
	}).Parse(tplIndex); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *NormalQuestionIndexController) Delone() {
	var code int
	var msg string
	url := beego.URLFor("NormalQuestionIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := NormalQuestion{}
	model.Id = id
	err := o.Read(&model)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	}
	_, err1 := o.Delete(&model, "Id")
	if err1 != nil {
		beego.Error("Delete NormalQuestion eror", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type NormalQuestionAddController struct {
	sysmanage.BaseController
}

func (c *NormalQuestionAddController) NestPrepare() {
	c.EnableRender = false
}

func (c *NormalQuestionAddController) Get() {
	c.Data["urlNormalQuestionIndexGet"] = c.URLFor("NormalQuestionIndexController.Get")
	c.Data["urlNormalQuestionAddPost"] = c.URLFor("NormalQuestionAddController.Post")

	if t, err := template.New("tplAddNormalQuestion.tpl").Parse(tplAdd); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *NormalQuestionAddController) Post() {
	var code int
	var msg string
	var url = beego.URLFor("NormalQuestionIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	model := NormalQuestion{}
	if err := c.ParseForm(&model); err != nil {
		msg = "参数异常"
		return
	}
	model.OrgId = c.LoginAdminOrgId
	model.Creator = c.LoginAdminId
	model.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Insert(&model); err != nil {
		msg = "添加失败"
		beego.Error("添加失败", err)
	} else {
		code = 1
		msg = "添加成功"
	}
}

type NormalQuestionEditController struct {
	sysmanage.BaseController
}

func (c *NormalQuestionEditController) NestPrepare() {
	c.EnableRender = false
}

func (c *NormalQuestionEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := NormalQuestion{}
	model.Id = id
	err := o.Read(&model)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(beego.URLFor("NormalQuestionIndexController.Get"), 302)
	}
	c.Data["data"] = &model

	c.Data["urlNormalQuestionIndexGet"] = c.URLFor("NormalQuestionIndexController.Get")
	c.Data["urlNormalQuestionEditPost"] = c.URLFor("NormalQuestionEditController.Post")

	if t, err := template.New("tplEditNormalQuestion.tpl").Funcs(map[string]interface{}{
		"date": beego.Date,
	}).Parse(tplEdit); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *NormalQuestionEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("NormalQuestionIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	model := NormalQuestion{}
	if err := c.ParseForm(&model); err != nil {
		msg = "参数异常"
		return
	}
	cols := []string{"Seq", "Question", "Answer"}
	model.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Update(&model, cols...); err != nil {
		msg = "更新失败"
		beego.Error("更新失败", err)
	} else {
		code = 1
		msg = "更新成功"
	}
}
