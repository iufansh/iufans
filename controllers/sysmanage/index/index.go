package index

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/controllers/sysmanage"
	. "github.com/iufansh/iufans/models"
	"github.com/iufansh/iuplugins/googleauth"
	"html/template"
	"net/http"
	"time"
)

type SysIndexController struct {
	sysmanage.BaseController
}

func (c *SysIndexController) NestPrepare()  {
	c.EnableRender = false
}

func (c *SysIndexController) Get() {
	o := orm.NewOrm()
	var admin = Admin{Id: c.LoginAdminId}
	if err := o.Read(&admin); err != nil {
		logs.Error("SysIndexController.Get read admin err:", err)
	}

	forceGaAuth := beego.AppConfig.DefaultInt("forcegaauth", 0)
	if forceGaAuth == 1 && admin.LoginVerify == 0 {
		c.Redirect(c.URLFor("SysIndexController.GetAuth"), http.StatusFound)
		return
	}

	c.Data["loginVerify"] = admin.LoginVerify

	c.Data["urlIndexGetAuth"] = c.URLFor("SysIndexController.GetAuth")
	c.Data["urlBackIndexGet"] = c.URLFor("BackIndexController.Get")

	if t, err := template.New("tplSysIndex.tpl").Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *SysIndexController) GetAuth() {
	user := "LOGIN-" + c.LoginAdminUsername
	ok, secret, qrCode := googleauth.GetGAuthQr(user)
	o := orm.NewOrm()
	if num, err := o.QueryTable(new(Admin)).Filter("Id", c.LoginAdminId).Update(orm.Params{
		"GaSecret":   secret,
		"Modifior":   c.LoginAdminId,
		"ModifyDate": time.Now(),
	}); err != nil || num != 1 {
		logs.Error("SysIndexController GetAuth", err, num)
		ok = false
	}
	c.Data["ok"] = ok
	c.Data["qrCode"] = qrCode
	c.Data["urlSysIndexPostAuth"] = c.URLFor("SysIndexController.PostAuth")

	if t, err := template.New("tplGaAuth.tpl").Parse(sysmanage.TplGaAuth); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *SysIndexController) PostAuth() {
	var code int
	var msg string
	var reUrl string
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &reUrl)
	authCode := c.GetString("auth_code")
	o := orm.NewOrm()
	var admin Admin
	if err := o.QueryTable(new(Admin)).Filter("Id", c.LoginAdminId).One(&admin, "GaSecret"); err != nil {
		msg = "绑定失败，请重试"
		return
	}

	if ok, err := googleauth.VerifyGAuth(admin.GaSecret, authCode); err != nil || !ok {
		logs.Error("SysIndexController PostAuth", err, ok)
		msg = "安全码验证失败，请确认"
		return
	}
	if num, err := o.QueryTable(new(Admin)).Filter("Id", c.LoginAdminId).Update(orm.Params{
		"LoginVerify": 2,
		"Modifior":    c.LoginAdminId,
		"ModifyDate":  time.Now(),
	}); err != nil || num != 1 {
		logs.Error("SysIndexController PostAuth", err, num)
		msg = "绑定失败，请重试"
		return
	}
	code = 1
	msg = "绑定成功"
	reUrl = c.URLFor("SysIndexController.Get")
}
