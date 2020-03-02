package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error404() {
	//c.Data["content"] = "page not found"
	c.TplName = "404.tpl"
}

func (c *ErrorController) Error401() {
	//c.Data["content"] = "server error"
	c.TplName = "401.tpl"
}
