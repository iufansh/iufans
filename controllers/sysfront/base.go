package sysfront

import (
	"github.com/astaxie/beego"
	"net/http"
	"fmt"
	"github.com/iufansh/iufans/utils"
	"net/url"
	"github.com/iufansh/iufans/controllers"
)

type FrontPreparer interface {
	FrontPreparer()
}

// 验证登录
type BaseFrontController struct {
	controllers.BaseMainController
	LoginMemberId   int64
	LoginMemberName string
	LoginMemberVip  int
}

func (c *BaseFrontController) Prepare() {
	beego.Info("\r\n----------request---------",
		"\r\nUri:", c.Ctx.Input.URI(),
		"\r\nMethod:", c.Ctx.Input.Method(),
		"\r\nFrom ip:", c.Ctx.Input.IP(),
		"\r\nUserAgent:", c.Ctx.Input.UserAgent(),
		"\r\nBody:", string(c.Ctx.Input.RequestBody),
		"\r\n--------------------------")
	var isLogon bool
	if memberId, ok := c.GetSession("loginMemberId").(int64); ok {
		var cacToken string
		if err := utils.GetCache(fmt.Sprintf("loginMemberId%d", memberId), &cacToken); err != nil || cacToken == "" {
			beego.Warn("BaseValidate1", err)
			c.Redirect(c.URLFor("LoginFrontController.Get", "redirect", url.QueryEscape(c.Ctx.Input.URI())), http.StatusFound)
			return
		}
		sesToken, ok := c.GetSession("loginMemberToken").(string)
		if !ok || cacToken != sesToken {
			beego.Warn("BaseValidate2 cacToken=", cacToken, "sesToken=", sesToken)
			c.Redirect(c.URLFor("LoginFrontController.Get", "redirect", url.QueryEscape(c.Ctx.Input.URI())), http.StatusFound)
			return
		}
		c.LoginMemberId = memberId
		isLogon = true
	} else {
		beego.Warn("BaseValidate3")
		c.Redirect(c.URLFor("LoginFrontController.Get", "redirect", url.QueryEscape(c.Ctx.Input.URI())), http.StatusFound)
		return
	}
	if name, ok := c.GetSession("loginMemberName").(string); ok {
		c.LoginMemberName = name
	} else {
		beego.Warn("BaseValidate4")
		c.Redirect(c.URLFor("LoginFrontController.Get", "redirect", url.QueryEscape(c.Ctx.Input.URI())), http.StatusFound)
		return
	}
	c.LoginMemberVip, _ = c.GetSession("loginMemberVip").(int)

	c.Data["isLogon"] = isLogon
	c.Data["loginMemberId"] = c.LoginMemberId
	c.Data["loginMemberName"] = c.LoginMemberName
	c.Data["loginMemberVip"] = c.LoginMemberVip
	if app, ok := c.AppController.(FrontPreparer); ok {
		app.FrontPreparer()
	}
	if c.Ctx.Input.Method() == http.MethodGet {
		c.Data["xsrf_token"] = c.XSRFToken()
	}
}

// 不验证登录
type Base2FrontController struct {
	controllers.BaseMainController
	IsLogon       bool
	LoginMemberId int64
}

func (c *Base2FrontController) Prepare() {
	beego.Info("\r\n----------request---------",
		"\r\nUri:", c.Ctx.Input.URI(),
		"\r\nMethod:", c.Ctx.Input.Method(),
		"\r\nFrom ip:", c.Ctx.Input.IP(),
		"\r\nUserAgent:", c.Ctx.Input.UserAgent(),
		"\r\nBody:", string(c.Ctx.Input.RequestBody),
		"\r\n--------------------------")
	if memberId, ok := c.GetSession("loginMemberId").(int64); ok {
		var cacToken string
		if err := utils.GetCache(fmt.Sprintf("loginMemberId%d", memberId), &cacToken); err != nil || cacToken == "" {
			c.IsLogon = false
		} else {
			sesToken, ok := c.GetSession("loginMemberToken").(string)
			if !ok || cacToken != sesToken {
				c.IsLogon = false
			} else {
				c.IsLogon = true
				c.LoginMemberId = memberId
			}
		}
	} else {
		c.IsLogon = false
	}
	if app, ok := c.AppController.(FrontPreparer); ok {
		app.FrontPreparer()
	}
	if c.Ctx.Input.Method() == http.MethodGet {
		c.Data["xsrf_token"] = c.XSRFToken()
	}
}

/*
func Retjson(ctx *context.Context, msg *string, code *int, data ...interface{}) {
	ret := make(map[string]interface{})
	ret["code"] = code
	ret["msg"] = msg
	if len(data) > 0 {
		d := data[0]
		switch d.(type) {
		case string:
			ret["url"] = d
			break
		case *string:
			ret["url"] = d
			break
		}
		ret["data"] = d
	}
	ctx.Output.JSON(ret, false, false)
}
*/
