package sysapi

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
	"github.com/iufansh/iutils"
	"time"
)

type loginQqParam struct {
	AccessToken string `json:"accessToken"`
	Nickname    string `json:"nickname"`
	OpenId      string `json:"openId"`
	Avatar      string `json:"avatar"`
}

type qqUnionIdResp struct {
	ClientId         string `json:"client_id"`
	Openid           string `json:"openid"`
	Unionid          string `json:"unionid"`
	ErrorCode        int    `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type LoginQqApiController struct {
	Base2ApiController
}

/*
api QQ登录
param:
body:{"accessToken":"access_token", "",}
return:{"code":1,"msg":"成功","data":{"id":1,"token":"11111111111111111111","nickname":"微信用户1","autoLogin":true,"accessToken":"ddfesfsf"}}
*/
func (c *LoginQqApiController) Post() {
	defer c.RetJSON()
	var p loginQqParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	if p.AccessToken == "" {
		c.Msg = "参数错误"
		return
	}
	// 获取unionId
	reqUrl := fmt.Sprintf("https://graph.qq.com/oauth2.0/me?access_token=%s&unionid=1&fmt=json", p.AccessToken)
	var unionInfo qqUnionIdResp
	if err := httplib.Get(reqUrl).ToJSON(&unionInfo); err != nil {
		logs.Error("LoginQqApiController get unionId err:", err)
		c.Msg = "验证失败S1"
		return
	}
	if unionInfo.ErrorCode != 0 || unionInfo.Unionid == "" {
		c.Msg = "验证失败S2"
		return
	}
	if unionInfo.Openid != p.OpenId {
		c.Msg = "验证失败S3"
		return
	}
	o := orm.NewOrm()
	var member models.Member
	if err := o.QueryTable(new(models.Member)).Filter("ThirdAuthId", unionInfo.Unionid).Limit(1).One(&member); err != nil && err != orm.ErrNoRows {
		logs.Error("LoginWxApiController QueryTable Member err:", err)
		c.Msg = "用户查询异常"
		return
	} else if err == orm.ErrNoRows {
		if member, err = CreateMemberReg(4, c.AppNo, c.AppChannel, c.AppVersionCode, 0, unionInfo.Unionid, unionInfo.Unionid, p.Nickname, unionInfo.Unionid, p.Avatar); err != nil {
			c.Msg = "登录失败，请重试"
			return
		}
	}

	// 自动登录
	member.LoginIp = c.Ctx.Input.IP()
	// 以下两个是用于统计登录次数
	member.AppNo = c.AppNo
	member.AppChannel = c.AppChannel
	member.AppVersion = c.AppVersionCode
	_, _, token := UpdateMemberLoginStatus(member)

	c.Code = utils.CODE_OK
	c.Msg = "登录成功"
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
		"vipExpire":  iutils.FormatDate(member.VipExpire),
	}
}
