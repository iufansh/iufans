package backtask

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/iufansh/iufans/controllers/sysmanage"
	"github.com/iufansh/iufans/taskback"
	"html/template"
)

type BackTaskIndexController struct {
	sysmanage.BaseController
}

func (c *BackTaskIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *BackTaskIndexController) Get() {

	list, _ := taskback.GetAllTaskBack()
	c.Data["dataList"] = &list

	c.Data["urlTaskBackIndexGet"] = c.URLFor("BackTaskIndexController.Get")

	if t, err := template.New("tplBackTaskIndex.tpl").Funcs(map[string]interface{}{
		"date": beego.Date,
	}).Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}
