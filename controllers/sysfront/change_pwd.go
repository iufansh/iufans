package sysfront

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/models"
	utils "github.com/iufansh/iutils"
	"strings"
)

type ChangePwdFrontController struct {
	BaseFrontController
}

func (c *ChangePwdFrontController) Get() {
	c.TplName = "front/changePwd.html"
}

func (c *ChangePwdFrontController) Post() {
	defer c.RetJSON()
	id := c.LoginMemberId
	if id <= 0 {
		c.Code = 1
		c.Msg = "请重新登录"
		c.Dta = c.URLFor("LoginFrontController.Get")
		return
	}
	oldPwd := strings.TrimSpace(c.GetString("oldPwd"))
	newPwd := strings.TrimSpace(c.GetString("newPwd"))
	newPwd2 := strings.TrimSpace(c.GetString("newPwd2"))
	if oldPwd == "" {
		c.Msg = "旧密码不能为空"
		return
	}
	if newPwd == "" {
		c.Msg = "新密码不能为空"
		return
	}
	if newPwd != newPwd2 {
		c.Msg = "确认密码不一致"
		return
	}
	o := orm.NewOrm()
	member := models.Member{Id: id}
	if err := o.Read(&member); err != nil {
		c.Msg = "用户信息错误，请重新登录"
		c.Dta = c.URLFor("LoginFrontController.Get")
	} else if utils.Md5(utils.Md5(oldPwd), utils.Pubsalt, member.Salt) != member.Password {
		c.Msg = "旧密码错误"
	} else {
		salt := utils.GetGuid()
		pa := utils.Md5(utils.Md5(newPwd), utils.Pubsalt, salt)
		member.Password = pa
		member.Salt = salt

		if _, err2 := o.Update(&member, "Password", "Salt", "ModifyDate"); err2 != nil {
			c.Msg = "修改失败，请重试"
			beego.Error("PwdFront Change password error", err2)
		} else {
			c.Code = 1
			c.Msg = "修改成功，请重新登录"
			c.Dta = c.URLFor("LoginFrontController.Get")
		}
	}
}
