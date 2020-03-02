package sysapi

import (
	. "github.com/iufansh/iufans/models"
	. "github.com/iufansh/iutils"
	"time"

	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/utils"
	"strconv"
	"strings"
)

type LoginApiController struct {
	Base2ApiController
}

type loginParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
api登录
param:
body:{"username":"aaaaa","password":"32md5小写"}
return:{"code":1,"msg":"成功","data":{"id":1,"token":"11111111111111111111","phone":"13111111111","nickname":"昵称","autoLogin":true}}
*/
func (c *LoginApiController) Post() {
	defer c.RetJSON()
	var p loginParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	if p.Username == "" {
		c.Msg = "用户名不能为空"
		return
	}
	if p.Password == "" {
		c.Msg = "密码不能为空"
		return
	}
	p.Password = strings.ToLower(p.Password)
	o := orm.NewOrm()
	member := Member{Username: p.Username}
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
		c.Msg = "账号已锁定，请通过忘记密码找回"
		return
	}
	if member.Password != Md5(p.Password, Pubsalt, member.Salt) {
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

	member.LoginIp = c.Ctx.Input.IP()
	code, msg, token := UpdateMemberLoginStatus(member)
	c.Code = code
	c.Msg = msg
	c.Dta = map[string]interface{}{
		"id":        member.Id,
		"token":     token,
		"phone":     member.Username,
		"nickname":  member.Name,
		"autoLogin": true,
	}
}

func UpdateMemberLoginStatus(member Member) (code int, msg, token string) {
	o := orm.NewOrm()
	lifeTime := beego.AppConfig.String("apitokenlifetime")
	expTime, _ := time.ParseDuration(lifeTime)

	member.TokenExpTime = time.Now().Add(expTime)
	token = Md5(member.Username, Pubsalt, strconv.FormatInt(time.Now().Unix(), 10))
	member.Token = token
	member.Locked = 0
	member.LoginFailureCount = 0
	member.LoginDate = time.Now()
	if num, err := o.Update(&member, "LoginFailureCount", "LoginIp", "LoginDate", "TokenExpTime", "Token"); err != nil || num != 1 {
		return utils.CODE_ERROR, "异常，请重试", ""
	}
	return utils.CODE_OK, "登录成功", token
}

func (c *LoginApiController) Logout() {
	defer c.RetJSON()
	if c.LoginMemberId <= 0 {
		c.Msg = "退出成功"
		c.Code = utils.CODE_OK
	}
	o := orm.NewOrm()
	if num, err := o.Update(&Member{Id: c.LoginMemberId, Token: ""}, "Token"); err != nil || num != 1 {
		beego.Error("Member logout err:", err)
		c.Msg = "退出失败"
		return
	}
	c.Code = utils.CODE_OK
	c.Msg = "退出成功"
}

type RefreshLoginApiController struct {
	BaseApiController
}

func (c *RefreshLoginApiController) Post() {
	defer c.RetJSON()

	if c.LoginMemberId <= 0 {
		c.Msg = "未登录"
		return
	}

	lifeTime := beego.AppConfig.String("apitokenlifetime")
	expTime, _ := time.ParseDuration(lifeTime)

	o := orm.NewOrm()
	o.QueryTable(new(Member)).Filter("Id", c.LoginMemberId).Update(orm.Params{
		"TokenExpTime": time.Now().Add(expTime),
		"LoginDate":    time.Now(),
	})
	c.Code = utils.CODE_OK
	c.Msg = "登录成功"
	c.Dta = "ok"
}
