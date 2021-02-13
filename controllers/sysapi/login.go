package sysapi

import (
	"fmt"
	"github.com/astaxie/beego/logs"
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
		logs.Error("Login error", err)
		c.Msg = "用户名或密码错误"
		return
	}
	if member.Cancelled == 1 {
		c.Msg = "用户名或密码错误" // 已注销
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
		if member.LoginFailureCount > 5 {
			member.Locked = 1
			cols = append(cols, "Locked")
		}
		cols = append(cols, "LoginFailureCount")
		if _, err := o.Update(&member, cols...); err != nil {
			logs.Error("MemberLogin update login fail err:", err)
		}
		if member.LoginFailureCount >= 3 {
			c.Msg = fmt.Sprintf("再错误%d次，将锁定账号", 6-member.LoginFailureCount)
		}
		c.Msg = "用户名或密码错误"
		return
	}

	member.LoginIp = c.Ctx.Input.IP()
	// 以下两个是用于统计登录次数
	member.AppNo = c.AppNo
	member.AppChannel = c.AppChannel
	member.AppVersion = c.AppVersionCode
	code, msg, token := UpdateMemberLoginStatus(member)
	c.Code = code
	c.Msg = msg
	var vipEffect int
	if member.Vip > 0 && !member.VipExpire.IsZero() && member.VipExpire.After(time.Now().AddDate(0, 0, -1)) {
		vipEffect = 1
	}
	c.Dta = map[string]interface{}{
		"id":         member.Id,
		"token":      token,
		"phone":      member.GetFmtMobile(),
		"nickname":   member.Name,
		"autoLogin":  true,
		"avatar":     member.GetFullAvatar(c.Ctx.Input.Site()),
		"inviteCode": utils.GenInviteCode(member.Id),
		"vipEffect":  vipEffect,
		"vip":        member.Vip,
		"vipExpire":  FormatDate(member.VipExpire),
	}
}

func UpdateMemberLoginStatus(member Member) (code int, msg, token string) {
	if TimeSub(time.Now(), member.LoginDate) >= 1 {
		go UpdateMemberLoginCount(member.AppNo, member.AppChannel, member.AppVersion)
	}
	o := orm.NewOrm()
	lifeTime := beego.AppConfig.String("apitokenlifetime")
	expTime, _ := time.ParseDuration(lifeTime)

	member.TokenExpTime = time.Now().Add(expTime)
	token = Md5(member.Username, Pubsalt, strconv.FormatInt(time.Now().Unix(), 10))
	member.Token = token
	member.Locked = 0
	member.LoginFailureCount = 0
	member.LoginDate = time.Now()
	if num, err := o.Update(&member, "LoginFailureCount", "Locked", "LoginIp", "LoginDate", "TokenExpTime", "Token"); err != nil || num != 1 {
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
		logs.Error("Member logout err:", err)
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
	var member Member
	if err := o.QueryTable(new(Member)).Filter("Id", c.LoginMemberId).One(&member, "LoginDate"); err != nil {
		logs.Error("RefreshLoginApiController.Post query member err:", err)
	} else if TimeSub(time.Now(), member.LoginDate) >= 1 {
		go UpdateMemberLoginCount(c.AppNo, c.AppChannel, c.AppVersionCode)
	}
	if _, err := o.QueryTable(new(Member)).Filter("Id", c.LoginMemberId).Update(orm.Params{
		"TokenExpTime": time.Now().Add(expTime),
		"LoginDate":    time.Now(),
	}); err != nil {
		logs.Error("RefreshLoginApiController.Post Update member err:", err)
		c.Msg = "登录失败"
		return
	}
	c.Code = utils.CODE_OK
	c.Msg = "登录成功"
	c.Dta = "ok"
}

func UpdateMemberLoginCount(appNo, appChannel string, appVersion int) {
	o := orm.NewOrm()
	if exist := o.QueryTable(new(MemberLoginCount)).Filter("AppNo", appNo).Filter("AppChannel", appChannel).Filter("AppVersion", appVersion).Filter("CountDate", FormatDate(time.Now())).Exist(); exist {
		if _, err := o.QueryTable(new(MemberLoginCount)).Filter("AppNo", appNo).Filter("AppChannel", appChannel).Filter("AppVersion", appVersion).Filter("CountDate", FormatDate(time.Now())).Update(orm.Params{
			"Count": orm.ColValue(orm.ColAdd, 1),
		}); err != nil {
			logs.Error("UpdateMemberLoginCount update err:", err)
		}
	} else {
		model := MemberLoginCount{
			AppNo:      appNo,
			AppChannel: appChannel,
			AppVersion: appVersion,
			CountDate:  ParseDate(time.Now()),
			Count:      1,
		}
		if _, err := o.Insert(&model); err != nil {
			logs.Error("UpdateMemberLoginCount Insert err:", err)
		}
	}
}
