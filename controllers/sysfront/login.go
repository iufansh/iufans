package sysfront

import (
	"fmt"
	"net/http"
	. "github.com/iufansh/iufans/models"
	fm "github.com/iufansh/iufans/models"
	. "github.com/iufansh/iufans/utils"
	. "github.com/iufansh/iutils"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"net/url"
)

type LoginFrontController struct {
	Base2FrontController
}

func (c *LoginFrontController) Get() {
	c.Data["msg"] = c.GetString("msg", "")
	c.Data["siteName"] = fm.GetSiteConfigValue(Scname)
	c.Data["redirect"] = c.GetString("redirect")
	c.TplName = "front/login.html"
}

func (c *LoginFrontController) Post() {
	defer c.RetJSON()
	redirect := c.GetString("redirect")
	username := c.GetString("username")
	pwd := c.GetString("password")
	defer func() {
		beego.Warn("LoginFrontRequest from Ip:", c.Ctx.Input.IP(), "账号:", username, "结果:", c.Msg)
	}()
	if username == "" {
		c.Msg = "用户名不能为空"
		return
	}
	if pwd == "" {
		c.Redirect(c.URLFor("LoginFrontController.Get"), http.StatusFound)
		return
	}
	if beego.BConfig.RunMode == "prod" && !GetCpt().VerifyReq(c.Ctx.Request) {
		c.Msg = "验证码错误"
		return
	}
	o := orm.NewOrm()
	member := Member{Username: username}
	if err := o.Read(&member, "Username"); err != nil {
		beego.Error("Login error", err)
		c.Msg = "用户名或密码错误"
		return
	}
	if member.Enabled == 0 {
		c.Msg = "账号已禁用"
		return
	}
	if member.Locked == 1 {
		c.Msg = "账号已锁定"
		return
	}
	if member.Password != Md5(pwd, Pubsalt, member.Salt) {
		cols := make([]string, 0)
		member.LoginFailureCount += 1
		if member.LoginFailureCount >= 5 {
			member.Locked = 1
			cols = append(cols, "Locked")
		}
		cols = append(cols, "LoginFailureCount")
		o.Update(&member, cols...)

		c.Msg = "用户名或密码错误"
		return
	}

	token := GetGuid()
	lifeTime, err := beego.AppConfig.Int("sessiongcmaxlifetimefront")
	if err != nil {
		lifeTime = 604800 // 7天
	}
	SetCache(fmt.Sprintf("loginMemberId%d", member.Id), token, lifeTime)
	c.SetSession("loginMemberToken", token)
	c.SetSession("loginMemberId", member.Id)
	c.SetSession("loginMemberName", member.Name)
	c.SetSession("loginMemberUsername", member.Username)
	c.SetSession("loginMemberVip", member.Vip)

	member.LoginFailureCount = 0
	member.LoginIp = c.Ctx.Input.IP()
	member.LoginDate = time.Now()
	o.Update(&member, "LoginFailureCount", "LoginIp", "LoginDate")

	c.Code = 1
	c.Msg = "登录成功"
	if redirect != "" {
		if c.Dta, err = url.QueryUnescape(redirect); err != nil {
			c.Dta = c.URLFor("HomeFrontController.Get")
		}
	} else {
		c.Dta = c.URLFor("HomeFrontController.Get")
	}
}

func (c *LoginFrontController) Logout() {
	DelCache(fmt.Sprintf("loginMemberId%v", c.GetSession("loginMemberId")))
	c.DelSession("loginMemberToken")
	c.DelSession("loginMemberId")
	c.DelSession("loginMemberName")
	c.DelSession("loginMemberUsername")
	c.DelSession("loginMemberVip")
	c.Redirect(c.URLFor("LoginFrontController.Get"), http.StatusFound)
}
