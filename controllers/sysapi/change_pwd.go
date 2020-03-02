package sysapi

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/models"
	utils2 "github.com/iufansh/iufans/utils"
	utils "github.com/iufansh/iutils"
	"strings"
)

type changePwdParam struct {
	NewPassword string `json:"newPwd"` // 必填
	OldPassword string `json:"oldPwd"` // 可选，不填不验证，填写后验证
}

type ChangePwdApiController struct {
	BaseApiController
}

/*
api修改密码
param:
body: {"newPwd":"32md5小写","oldPwd":"32md5小写"}
return:{"code":1,"msg":"成功"}
desc: 密码修改成功后，不需要重新登录
*/
func (c *ChangePwdApiController) Post() {
	defer c.RetJSON()
	var p changePwdParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils2.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	if p.NewPassword == "" {
		c.Msg = "新密码不能为空"
		return
	}
	p.NewPassword = strings.ToLower(p.NewPassword)
	o := orm.NewOrm()

	member := models.Member{Id: c.LoginMemberId}
	if p.OldPassword != "" {
		if err := o.Read(&member); err != nil {
			c.Msg = "异常，请重试"
			return
		}
		if utils.Md5(p.OldPassword, utils.Pubsalt, member.Salt) != member.Password {
			c.Msg = "旧密码错误"
			return
		}
	}
	salt := utils.GetGuid()
	pa := utils.Md5(p.NewPassword, utils.Pubsalt, salt)
	member.Password = pa
	member.Salt = salt

	if _, err2 := o.Update(&member, "Password", "Salt", "ModifyDate"); err2 != nil {
		c.Msg = "修改失败，请重试"
		beego.Error("PwdFront Change password error", err2)
		return
	}
	c.Code = utils2.CODE_OK
	c.Msg = "修改成功"
}
