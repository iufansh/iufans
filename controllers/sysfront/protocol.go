package sysfront

import (
	fm "github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
)

type ProtocolFrontController struct {
	Base2FrontController
}

func (c *ProtocolFrontController) Get() {
	c.Data["siteName"] = fm.GetSiteConfigValue(utils.Scname)
	c.TplName = "front/protocol.html"
}
