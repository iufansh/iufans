package sysapi

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/iGoogle-ink/gopay/wechat"
	"github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
)

type loginHuaweiParam struct {
	AppId string `json:"appId"`
	Code  string `json:"code"`
}

type LoginHuaweiApiController struct {
	Base2ApiController
}

/*
api huawei登录
param:
body:{"appId":"wxappid1111111","code":"111122"}
return:{"code":1,"msg":"成功","data":{"id":1,"token":"11111111111111111111","nickname":"微信用户1","autoLogin":true,"accessToken":"ddfesfsf"}}
*/
func (c *LoginHuaweiApiController) Post() {
	defer c.RetJSON()
	var p loginHuaweiParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	if p.AppId == "" {
		c.Msg = "AppId不能为空"
		return
	}
	if p.Code == "" {
		c.Msg = "Code不能为空"
		return
	}
	o := orm.NewOrm()
	var pc models.PaymentConfig
	o.QueryTable(new(models.PaymentConfig)).Filter("AppId", p.AppId).Limit(1).One(&pc)
	var vo models.WechatVo
	if err := json.Unmarshal([]byte(pc.ConfValue), &vo); err != nil {
		logs.Error("Unmarshal ConfValue err:", err)
		c.Msg = "接口异常(1)"
		return
	}
	// 获取access token
	accessToken, err := wechat.GetAppLoginAccessToken(p.AppId, vo.AppSecret, p.Code)
	//logs.Info("accessToken:", fmt.Sprintf("%+v", accessToken))
	if err != nil {
		logs.Error("wechat.GetAppLoginAccessToken err:", err)
		c.Msg = "Access Token获取异常"
		return
	} else if accessToken.Errcode > 0 || accessToken.Errmsg != "" {
		logs.Error("wechat.GetAppLoginAccessToken err:", accessToken.Errmsg)
		c.Msg = "Access Token获取失败"
		return
	} else if accessToken.Unionid == "" {
		c.Msg = "Access Token获取失败"
		return
	}
	// 获取用户信息
	userInfo, err := GetUserInfo(accessToken.AccessToken, accessToken.Openid)
	//logs.Info("userInfo:", fmt.Sprintf("%+v", userInfo))
	if err != nil {
		logs.Error("wechat.GetUserInfo err:", err)
		c.Msg = "用户信息获取异常"
		return
	} else if userInfo.Errcode > 0 || userInfo.Errmsg != "" {
		logs.Error("wechat.GetUserInfo err:", userInfo.Errmsg)
		c.Msg = "用户信息获取失败"
		return
	} else if userInfo.Unionid == "" {
		c.Msg = "用户信息获取失败"
		return
	}
	var member models.Member
	if err := o.QueryTable(new(models.Member)).Filter("ThirdAuthId", userInfo.Unionid).Limit(1).One(&member); err != nil && err != orm.ErrNoRows {
		logs.Error("QueryTable Member err:", err)
		c.Msg = "用户查询异常"
		return
	} else if err == orm.ErrNoRows {
		if member, err = CreateMemberReg(0, userInfo.Unionid, userInfo.Unionid, userInfo.Nickname, userInfo.Unionid); err != nil {
			c.Msg = "注册失败"
			return
		}
	}

	// 自动登录
	member.LoginIp = c.Ctx.Input.IP()
	_, _, token := UpdateMemberLoginStatus(member)

	c.Code = utils.CODE_OK
	c.Msg = "获取成功"
	c.Dta = map[string]interface{}{
		"id":          member.Id,
		"token":       token,
		"phone":       "",
		"nickname":    member.Name,
		"autoLogin":   true,
		"accessToken": accessToken.AccessToken, // 微信access token
	}
}
