package sysapi

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	. "github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
	. "github.com/iufansh/iutils"
	"strconv"
	"strings"
)

type MemberApiController struct {
	BaseApiController
}

type bindParam struct {
	AuthCode string `json:"authCode"`
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
api已通过其他方式登录的（如微信），通过此接口绑定手机号
param:
body:{"authCode":2356,"username":"13111111111","password":"32md5小写"}
return:{"code":1,"msg":"成功","data":"ok"}
desc:
*/
func (c *MemberApiController) BindPhone() {
	defer c.RetJSON()
	var p bindParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	p.Password = strings.ToLower(p.Password)

	model := Member{
		Username: p.Username,
		Password: p.Password,
	}
	if p.AuthCode == "" {
		c.Msg = "短信验证码必填"
		return
	}
	if ok := utils.VerifySmsVerifyCode(model.Username, p.AuthCode); !ok {
		c.Msg = "短信验证码错误"
		return
	}
	valid := validation.Validation{}
	valid.Required(model.Username, "errmsg").Message("手机号必填")
	valid.MaxSize(model.Username, 11, "errmsg").Message("手机号最长11位")
	valid.Required(model.Password, "errmsg").Message("密码必填")
	valid.Length(model.Password, 32, "errmsg").Message("密码格式错误")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.Msg = err.Message
			return
		}
	}
	// 验证用户名是否存在
	o := orm.NewOrm()
	if isExist := o.QueryTable(new(Member)).Filter("Username", model.Username).Exist(); isExist {
		c.Msg = "当前手机号已存在账号，请换其他手机号"
		return
	}

	salt := GetGuid()
	model.Password = Md5(model.Password, Pubsalt, salt)
	model.Salt = salt
	model.Mobile = model.Username

	model.Id = c.LoginMemberId

	if num, err := o.Update(&model, "Username", "Mobile", "Password", "Salt"); err != nil || num != 1 {
		c.Msg = "绑定失败，请重试"
		return
	}

	c.Msg = "绑定成功"
	c.Code = utils.CODE_OK
	c.Dta = "绑定成功"
}

/*
api注销账号
param:
body:
return:{"code":1,"msg":"成功","data":"ok"}
desc:
*/
func (c *MemberApiController) CancelAccount() {
	defer c.RetJSON()
	member := Member{}
	member.Id = c.LoginMemberId
	member.Username = "mCancelled_" + strconv.FormatInt(member.Id, 10)
	member.Password = "mCancelled"
	member.ThirdAuthId = ""
	member.Token = ""
	member.Cancelled = 1
	o := orm.NewOrm()
	if num, err := o.Update(&member, "Username", "Password", "ThirdAuthId", "Token", "Cancelled"); err != nil || num != 1 {
		c.Msg = "注销失败，请重试"
		logs.Error("MemberApiController Update member error:", err)
		return
	}
	if err := utils.DelCache(fmt.Sprintf("loginMemberId%d", member.Id)); err != nil {
		logs.Error("MemberApiController DelCache error:", err)
	}

	c.Msg = "注销成功"
	c.Code = utils.CODE_OK
	c.Dta = "ok"
}
