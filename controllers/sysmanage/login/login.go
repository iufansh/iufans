package login

import (
	"fmt"
	"html/template"
	. "github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/controllers/sysmanage"
	. "github.com/iufansh/iutils"
	"time"
	. "github.com/iufansh/iufans/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iuplugins/googleauth"
	"net/http"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	beego.Warn("Login Get from ip:", c.Ctx.Input.IP())
	if beego.BConfig.RunMode == "dev" {
		c.Data["username"] = "admin"
		c.Data["pass"] = "111111"
		c.Data["captchaValue"] = "1"
	} else {
		c.Data["username"] = ""
		c.Data["pass"] = ""
		c.Data["captchaValue"] = ""
	}
	c.Data["year"] = time.Now().Year()
	c.Data["siteName"] = GetSiteConfigValue(Scname)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())

	c.Data["urlLoginPost"] = c.URLFor("LoginController.Post")
	c.Data["urlLoginVerify"] = c.URLFor("LoginController.LoginVerify")

	if t, err := template.New("tplLogin.tpl").Funcs(map[string]interface{}{
		"create_captcha": GetCpt().CreateCaptchaHTML,
	}).Parse(tplLogin); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *LoginController) Post() {
	ret := make(map[string]interface{})
	username := c.GetString("username")
	pwd := c.GetString("password")
	beego.Info("Login username=", username, "password=", pwd)
	defer func() {
		beego.Warn("LoginRequest from Ip:", c.Ctx.Input.IP(), "账号:", username, "结果:", ret["msg"])
		c.Data["json"] = &ret
		c.ServeJSON()
	}()
	if username == "" {
		ret["msg"] = "用户名不能为空"
		return
	}
	if pwd == "" {
		c.Redirect(c.URLFor("LoginController.Get"), http.StatusFound)
		return
	}
	if beego.BConfig.RunMode == "prod" && !GetCpt().VerifyReq(c.Ctx.Request) {
		ret["msg"] = "验证码错误"
		return
	}
	o := orm.NewOrm()
	admin := Admin{Username: username}
	if err := o.Read(&admin, "Username"); err != nil {
		beego.Error("Login error", err)
		ret["msg"] = "用户名或密码错误"
		return
	}
	if admin.Enabled == 0 {
		ret["msg"] = "账号已禁用"
		return
	}
	if admin.Locked == 1 {
		ret["msg"] = "账号已锁定"
		return
	}
	if admin.Password != Md5(pwd, Pubsalt, admin.Salt) {
		cols := make([]string, 0)
		admin.LoginFailureCount += 1
		if admin.LoginFailureCount >= 5 {
			admin.Locked = 1
			cols = append(cols, "Locked")
		}
		cols = append(cols, "LoginFailureCount")
		o.Update(&admin, cols...)

		ret["msg"] = "用户名或密码错误"
		return
	}
	// ip黑白名单验证
	ipValid := false
	curIp := c.Ctx.Input.IP()
	var lists []IpList
	// 白名单验证
	if _, err := o.QueryTable(new(IpList)).Filter("OrgId", admin.OrgId).Filter("Black", 0).All(&lists, "Ip"); err != nil {
		ret["msg"] = "登录异常，请重试"
		return
	}
	if len(lists) > 0 {
		for _, v := range lists {
			if v.Ip == curIp {
				ipValid = true
				break
			}
		}
		if !ipValid {
			ret["msg"] = "IP不可用，请联系管理员"
			return
		}
	} else { // 白名单未配置，则验证黑名单
		if exists := o.QueryTable(new(IpList)).Filter("OrgId", admin.OrgId).Filter("Black", 1).Filter("Ip", curIp).Exist(); exists {
			ret["msg"] = "非法IP，请联系管理员"
			return
		}
	}

	if admin.LoginVerify == 2 { // 需要谷歌安全码验证
		ret["code"] = 3
		ret["msg"] = "请输入谷歌安全码"
		return
	}

	token := GetGuid()
	lifeTime, err := beego.AppConfig.Int("sessiongcmaxlifetime")
	if err != nil {
		lifeTime = 3600
	}
	SetCache(fmt.Sprintf("loginAdminId%d", admin.Id), token, lifeTime)
	c.SetSession("token", token)
	c.SetSession("loginAdminId", admin.Id)
	c.SetSession("loginAdminOrgId", admin.OrgId)
	c.SetSession("loginAdminName", admin.Name)
	c.SetSession("loginAdminUsername", admin.Username)

	admin.LoginFailureCount = 0
	admin.LoginIp = c.Ctx.Input.IP()
	admin.LoginDate = time.Now()
	o.Update(&admin, "LoginFailureCount", "LoginIp", "LoginDate")

	ret["code"] = 1
	ret["msg"] = "登录成功"
	ret["url"] = c.URLFor("BaseIndexController.Get")
}

func (c *LoginController) LoginVerify() {
	var code int
	var msg string
	var reurl string
	username := c.GetString("username")
	defer func() {
		sysmanage.Retjson(c.Ctx, &msg, &code, &reurl)
		beego.Warn("LoginVerify from Ip:", c.Ctx.Input.IP(), "账号:", username, "结果:", msg)
	}()
	verifyCode := c.GetString("code")
	verifyType, _ := c.GetInt("verify", 2)
	if verifyCode == "" {
		msg = "验证码不能为空"
		return
	}

	o := orm.NewOrm()
	admin := Admin{Username: username}
	if err := o.Read(&admin, "Username"); err != nil {
		msg = "验证失败，请重试"
		return
	}
	var isVerify bool
	if verifyType == 3 {
		isVerify, _ = googleauth.VerifyGAuth(admin.GaSecret, verifyCode)
	}
	if !isVerify {
		msg = "验证失败"
		return
	}

	token := GetGuid()
	lifeTime, err := beego.AppConfig.Int("sessiongcmaxlifetime")
	if err != nil {
		lifeTime = 3600
	}
	SetCache(fmt.Sprintf("loginAdminId%d", admin.Id), token, lifeTime)
	c.SetSession("token", token)
	c.SetSession("loginAdminId", admin.Id)
	c.SetSession("loginAdminOrgId", admin.OrgId)
	c.SetSession("loginAdminName", admin.Name)
	c.SetSession("loginAdminUsername", admin.Username)

	admin.LoginFailureCount = 0
	admin.LoginIp = c.Ctx.Input.IP()
	admin.LoginDate = time.Now()
	o.Update(&admin, "LoginFailureCount", "LoginIp", "LoginDate")

	code = 1
	msg = "验证成功"
	reurl = c.URLFor("BaseIndexController.Get")
}

func (c *LoginController) Logout() {
	DelCache(fmt.Sprintf("loginAdminId%v", c.GetSession("loginAdminId")))
	c.DelSession("token")
	c.DelSession("loginAdminId")
	c.DelSession("loginAdminOrgId")
	c.DelSession("loginAdminName")
	c.DelSession("loginAdminUsername")
	c.Redirect(c.URLFor("LoginController.Get"), 302)
}
