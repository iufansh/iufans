package login

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/controllers/sysmanage"
	. "github.com/iufansh/iufans/models"
	. "github.com/iufansh/iufans/utils"
	"github.com/iufansh/iuplugins/googleauth"
	. "github.com/iufansh/iutils"
	"html/template"
	"net/http"
	"time"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	logs.Warn("Login Get from ip:", c.Ctx.Input.IP())
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

	// 登录页面模板
	var tpl string
	loginTpl := beego.AppConfig.DefaultInt("adminlogintpl", -1)
	switch loginTpl {
	case 0:
		tpl = tplLogin
		break
	default:
		tpl = tplLoginV1
	}

	if t, err := template.New("tplLogin.tpl").Funcs(map[string]interface{}{
		"create_captcha": GetCpt().CreateCaptchaHTML,
	}).Parse(tpl); err != nil {
		logs.Error("template Parse err", err)
	} else {
		if err := t.Execute(c.Ctx.ResponseWriter, c.Data); err != nil {
			logs.Error("LoginController.Get execute response err:", err)
		}
	}
}

func (c *LoginController) Post() {
	ret := make(map[string]interface{})
	username := c.GetString("username")
	pwd := c.GetString("password")
	logs.Info("Login username=", username, "password=", pwd)
	defer func() {
		logs.Warn("LoginRequest from Ip:", c.Ctx.Input.IP(), "账号:", username, "结果:", ret["msg"])
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
		logs.Error("Login error", err)
		ret["msg"] = "用户名或密码错误"
		return
	}
	if admin.Enabled == 0 {
		ret["msg"] = "账号不可用"
		return
	}
	if admin.Locked == 1 {
		ret["msg"] = "账号不可用"
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
		if _, err := o.Update(&admin, cols...); err != nil {
			logs.Error("Login update admin err:", err)
		}

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
		ret["msg"] = "请输入安全码"
		return
	}

	token := GetGuid()
	lifeTime, err := beego.AppConfig.Int("sessiongcmaxlifetime")
	if err != nil {
		lifeTime = 3600
	}
	if err := SetCache(fmt.Sprintf("loginAdminId%d", admin.Id), token, lifeTime); err != nil {
		ret["msg"] = "登录失败C1"
		return
	}
	c.SetSession("token", token)
	c.SetSession("loginAdminId", admin.Id)
	c.SetSession("loginAdminOrgId", admin.OrgId)
	c.SetSession("loginAdminName", admin.Name)
	c.SetSession("loginAdminUsername", admin.Username)

	admin.LoginFailureCount = 0
	admin.LoginIp = c.Ctx.Input.IP()
	admin.LoginDate = time.Now()
	if num, err := o.Update(&admin, "LoginFailureCount", "LoginIp", "LoginDate"); err != nil {
		ret["msg"] = "登录失败U1"
		return
	} else if num != 1 {
		ret["msg"] = "登录失败N1"
		return
	}

	ret["code"] = 1
	ret["msg"] = "登录成功"
	ret["url"] = c.URLFor("BaseIndexController.Get")
}

func (c *LoginController) LoginVerify() {
	var code int
	var msg string
	var reUrl string
	username := c.GetString("username")
	defer func() {
		sysmanage.Retjson(c.Ctx, &msg, &code, &reUrl)
		logs.Warn("LoginVerify from Ip:", c.Ctx.Input.IP(), "账号:", username, "结果:", msg)
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
	if err := SetCache(fmt.Sprintf("loginAdminId%d", admin.Id), token, lifeTime); err != nil {
		msg = "登录失败C2"
		return
	}
	c.SetSession("token", token)
	c.SetSession("loginAdminId", admin.Id)
	c.SetSession("loginAdminOrgId", admin.OrgId)
	c.SetSession("loginAdminName", admin.Name)
	c.SetSession("loginAdminUsername", admin.Username)

	admin.LoginFailureCount = 0
	admin.LoginIp = c.Ctx.Input.IP()
	admin.LoginDate = time.Now()
	if num, err := o.Update(&admin, "LoginFailureCount", "LoginIp", "LoginDate"); err != nil {
		msg= "登录失败U2"
		return
	} else if num != 1 {
		msg = "登录失败N2"
		return
	}

	code = 1
	msg = "验证成功"
	reUrl = c.URLFor("BaseIndexController.Get")
}

func (c *LoginController) Logout() {
	if err := DelCache(fmt.Sprintf("loginAdminId%v", c.GetSession("loginAdminId"))); err != nil {
		logs.Error("LoginController.Logout DelCache err:", err)
	}
	c.DelSession("token")
	c.DelSession("loginAdminId")
	c.DelSession("loginAdminOrgId")
	c.DelSession("loginAdminName")
	c.DelSession("loginAdminUsername")
	c.Redirect(c.URLFor("LoginController.Get"), 302)
}
