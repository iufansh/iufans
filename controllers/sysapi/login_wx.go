package sysapi

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/wechat"
	"github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
)

type loginWxParam struct {
	Code string `json:"code"`
}

type LoginWxApiController struct {
	Base2ApiController
}

/*
api WX登录
param:
body:{"code":"111122"}
return:{"code":1,"msg":"成功","data":{"id":1,"token":"11111111111111111111","nickname":"微信用户1","autoLogin":true,"accessToken":"ddfesfsf"}}
*/
func (c *LoginWxApiController) Post() {
	defer c.RetJSON()
	var p loginWxParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	if p.Code == "" {
		c.Msg = "Code不能为空"
		return
	}
	o := orm.NewOrm()
	var pc models.PaymentConfig
	o.QueryTable(new(models.PaymentConfig)).Filter("AppNo", c.AppNo).Filter("PayType", utils.PayTypeWechatPay).Limit(1).One(&pc)
	var vo models.WechatVo
	if err := json.Unmarshal([]byte(pc.ConfValue), &vo); err != nil {
		logs.Error("Unmarshal ConfValue err:", err)
		c.Msg = "接口异常(1)"
		return
	}
	// 获取access token
	accessToken, err := wechat.GetAppLoginAccessToken(pc.AppId, vo.AppSecret, p.Code)
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
		if member, err = CreateMemberReg(c.AppNo, c.AppChannel, c.AppVersionCode, 0, userInfo.Unionid, userInfo.Unionid, userInfo.Nickname, userInfo.Unionid, userInfo.Headimgurl); err != nil {
			c.Msg = "登录失败，请重试"
			return
		}
	}

	// 自动登录
	member.LoginIp = c.Ctx.Input.IP()
	_, _, token := UpdateMemberLoginStatus(member)

	c.Code = utils.CODE_OK
	c.Msg = "登录成功"
	c.Dta = map[string]interface{}{
		"id":         member.Id,
		"token":      token,
		"phone":      "",
		"nickname":   member.Name,
		"autoLogin":  true,
		"avatar":     member.Avatar,
		"inviteCode": utils.GenInviteCode(member.Id),
		// "accessToken": accessToken.AccessToken, // 微信access token
	}
}

// 获取用户基本信息(UnionID机制)
//    accessToken：接口调用凭据
//    openId：用户的OpenID
//    lang:默认为 zh_CN ，可选填 zh_CN 简体，zh_TW 繁体，en 英语
//    获取用户基本信息(UnionID机制)文档：https://developers.weixin.qq.com/doc/oplatform/Mobile_App/WeChat_Login/Authorized_API_call_UnionID.html
func GetUserInfo(accessToken, openId string, lang ...string) (userInfo *wechat.UserInfo, err error) {
	userInfo = new(wechat.UserInfo)
	url := "https://api.weixin.qq.com/sns/userinfo?access_token=" + accessToken + "&openid=" + openId + "&lang=zh_CN"
	if len(lang) > 0 {
		url = "https://api.weixin.qq.com/sns/userinfo?access_token=" + accessToken + "&openid=" + openId + "&lang=" + lang[0]
	}
	_, errs := gopay.NewHttpClient().Get(url).EndStruct(userInfo)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return userInfo, nil
}
