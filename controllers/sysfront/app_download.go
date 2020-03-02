package sysfront

import (
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/models"
	"net/http"
)

type AppDownloadFrontController struct {
	Base2FrontController
}

func (c *AppDownloadFrontController) DownloadRedirect() {
	appNo := c.Ctx.Input.Param(":appNo")
	o := orm.NewOrm()
	var appVersion models.AppVersion
	if err := o.QueryTable(new(models.AppVersion)).Filter("AppNo", appNo).OrderBy("-VersionNo").One(&appVersion, "DownloadUrl"); err != nil {
		c.Abort("404")
		return
	}
	c.Redirect(appVersion.DownloadUrl, http.StatusFound)
}
