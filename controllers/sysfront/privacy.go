package sysfront

import (
	fm "github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
)

type PrivacyFrontController struct {
	Base2FrontController
}

func (c *PrivacyFrontController) Get() {
	m := fm.GetSiteConfigMap(utils.Scname, utils.Sccompanyaddress, utils.Sccompanyconcattel, utils.Sccompanyconcatqq, utils.Sccompanyname)
	c.Data["siteName"] = m[utils.Scname]
	c.Data["companyName"] = m[utils.Sccompanyname]
	c.Data["companyAddress"] = m[utils.Sccompanyaddress]
	c.Data["companyConcatTel"] = m[utils.Sccompanyconcattel]
	c.Data["companyConcatQQ"] = m[utils.Sccompanyconcatqq]
	c.TplName = "front/privacy.html"
}
