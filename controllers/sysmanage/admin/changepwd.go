package admin

import (
	"html/template"
	"github.com/iufansh/iufans/controllers/sysmanage"
	. "github.com/iufansh/iufans/models"
	. "github.com/iufansh/iutils"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ChangePwdController struct {
	sysmanage.BaseController
}

func (c *ChangePwdController) NestPrepare()  {
	c.EnableRender = false
}

func (c *ChangePwdController) Get() {
	c.Data["urlAdminChangePwd"] = c.URLFor("ChangePwdController.Post")
	if t, err := template.New("tplChangePwd.tpl").Parse(tplChangePwd); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *ChangePwdController) Post() {
	var code int
	var msg string
	var url string
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	id := c.LoginAdminId
	if id == 0 {
		code = 1
		msg = "请重新登录"
		url = c.URLFor("LoginController.Get")
		return
	}
	oldPwd := c.GetString("oldPassword")
	newPwd := c.GetString("newPassword")
	reNewPwd := c.GetString("reNewPassword")
	if oldPwd == "" || strings.TrimSpace(oldPwd) == "" {
		msg = "旧密码不能为空"
		return
	} else if newPwd == "" || strings.TrimSpace(newPwd) == "" {
		msg = "新密码不能为空"
		return
	} else if strings.TrimSpace(newPwd) != strings.TrimSpace(reNewPwd) {
		msg = "两次输入的新密码不一致"
		return
	}
	o := orm.NewOrm()
	admin := Admin{Id: id}
	if err := o.Read(&admin); err != nil {
		msg = "用户信息错误，请重试"
	} else if Md5(oldPwd, Pubsalt, admin.Salt) != admin.Password {
		msg = "旧密码错误"
	} else {
		salt := GetGuid()
		pa := Md5(newPwd, Pubsalt, salt)
		admin.Password = pa
		admin.Salt = salt

		if _, err2 := o.Update(&admin, "Password", "Salt", "ModifyDate"); err2 != nil {
			msg = "更新失败"
			beego.Error("Change password error", err2)
		} else {
			code = 1
			msg = "更新成功，请重新登录"
			url = c.URLFor("LoginController.Logout")
		}
	}
}
