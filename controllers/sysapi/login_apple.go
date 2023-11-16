package sysapi

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
	"github.com/iufansh/iutils"
	"time"
)

/*
uni.login({provider: 'apple'}) 获得结果
{
	"authResult": {
		"access_token": "eyJraWQiOiI4NkQ4OEtmIiwiYWxnIjoiUlMyNTYifQ.eyJpc3MiOiJodHRwczovL2FwcGxlaWQuYXBwbGUuY29tIiwiYXVkIjoiY29tLnFpZXh1ZS5zZWFob3JzZXJlYWQiLCJleHAiOjE2MjIxMDQzMDEsImlhdCI6MTYyMjAxNzkwMSwic3ViIjoiMDAwMzY3LmJhZmJlMjRjMjQ1MjQ2NTk4MTRiYTE5M2U1ZDg3ZjMxLjA4MzEiLCJjX2hhc2giOiJrdUtRUWlNTDdYUFdzaE9QZV9YWDRRIiwiZW1haWwiOiI3d2Y5bTV0amc0QHByaXZhdGVyZWxheS5hcHBsZWlkLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjoidHJ1ZSIsImlzX3ByaXZhdGVfZW1haWwiOiJ0cnVlIiwiYXV0aF90aW1lIjoxNjIyMDE3OTAxLCJub25jZV9zdXBwb3J0ZWQiOnRydWV9.D6SCH01B0dAASGCXih0hkdMeA8sIczpblQrfJtREZeXwM9MG3CyRUtjwWsPHfI6L5BRZJzTaifwW3uwmBit-jDFdfcFjbqYUrMgQ2fpqFa5a6BLmWKRaKRHhEFwWEjDvwr9w1sJLA4kasRZZPje_afkJiDhij-QqAsTxQJ_wq3fRfENREP_TRQMytLRfIXfytW0MqqVZlQEwn_dYIs42lP3rf3QFYrtzdDgGHBUe3I3jimE9vHGj_lwrzhv9StIoqF7FgcHQrVT9ggDrPjWArFS_5WFMyYxfgTO9KC1aXS8Rw0QAMM2uAqlSGRsC9EjHEPUvZQHw9LlqXfzOkRHnMg",
		"openid": "012367.bafbe24c24384939814ba193e5d87f31.0821"
	},
	"errMsg": "login:ok"
}

uni.getUserInfo({ provider: 'apple' }) 获得结果
{
	"errMsg": "getUserInfo:ok",
	"userInfo": {
		"openId": "012367.bafbe24c24384939814ba193e5d87f31.0821",
		"fullName": {
			"familyName": "周",
			"giveName": "XX",
			"givenName": "XX"
		},
		"email": "7wf9sf3tjg4@privaterelay.appleid.com",
		"authorizationCode": "c07bb8cesdf3402293a206458a48ac3c.0.rtwx.6QhSsnxTxEhe8wKgJ_RyVg",
		"identityToken": "eyJraWQiOiI4NkQ4OEsfg4sWxnIjoiUlMyNTYifQ.eyJpc3MiOiJodHRwczovL2FwcGxlaWQuYXBwbGUuY29tIiwiYXVkIjoiY29tLnFpZXh1ZS5zZWFob3JzZXJlYWQiLCJleHAiOjE2MjIxMDQzMDEsImlhdCI6MTYyMjAxNzkwMSwic3ViIjoiMDAwMzY3LmJhZmJlMjRjMjQ1MjQ2NTk4MTRiYTE5M2U1ZDg3ZjMxLjA4MzEiLCJjX2hhc2giOiJrdUtRUWlNTDdYUFdzaE9QZV9YWDRRIiwiZW1haWwiOiI3d2Y5bTV0amc0QHByaXZhdGVyZWxheS5hcHBsZWlkLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjoidHJ1ZSIsImlzX3ByaXZhdGVfZW1haWwiOiJ0cnVlIiwiYXV0aF90aW1lIjoxNjIyMDE3OTAxLCJub25jZV9zdXBwb3J0ZWQiOnRydWV9.D6SCH01B0dAASGCXih0hkdMeA8sIczpblQrfJtREZeXwM9MG3CyRUtjwWsPHfI6L5BRZJzTaifwW3uwmBit-jDFdfcFjbqYUrMgQ2fpqFa5a6BLmWKRaKRHhEFwWEjDvwr9w1sJLA4kasRZZPje_afkJiDhij-QqAsTxQJ_wq3fRfENREP_TRQMytLRfIXfytW0MqqVZlQEwn_dYIs42lP3rf3QFYrtzdDgGHBUe3I3jimE9vHGj_lwrzhv9StIoqF7FgcHQrVT9ggDrPjWArFS_5WFMyYxfgTO9KC1aXS8Rw0QAMM2uAqlSGRsC9EjHEPUvZQHw9LlqXfzOkRHnMg",
		"realUserStatus": 2
	}
}
 */

type loginAppleParam struct {
	OpenId      string `json:"openId"`
	Name string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
	Code string `json:"code"`
}

type LoginAppleApiController struct {
	Base2ApiController
}

/*
api Apple登录
param:
body:{"openid":"12121212"}
return:{"code":1,"msg":"成功","data":{"id":1,"token":"11111111111111111111","nickname":"用户1","autoLogin":true,"accessToken":"ddfesfsf"}}
*/
func (c *LoginAppleApiController) Post() {
	defer c.RetJSON()
	var p loginAppleParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	if p.OpenId == "" {
		c.Msg = "参数错误"
		return
	}
	
	o := orm.NewOrm()
	var member models.Member
	if err := o.QueryTable(new(models.Member)).Filter("ThirdAuthId", p.OpenId).Limit(1).One(&member); err != nil && err != orm.ErrNoRows {
		logs.Error("LoginAppleApiController QueryTable Member err:", err)
		c.Msg = "查询异常，请重试"
		return
	} else if err == orm.ErrNoRows {
		if p.Token == "" { // TODO 先开放，不验证
			//c.Msg = "获取失败，请重试"
			//return
		}
		var username string
		if p.Email != "" {
			username = p.Email
		} else {
			username = p.OpenId
		}
		if member, err = CreateMemberReg(6, c.AppNo, c.AppChannel, c.AppVersionCode, 0, username, p.OpenId, p.Name, p.OpenId, ""); err != nil {
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
