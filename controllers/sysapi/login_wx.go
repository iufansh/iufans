package sysapi

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/iGoogle-ink/gopay/wechat"
	"github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
	"github.com/iufansh/iutils"
	"time"
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
	if err := o.QueryTable(new(models.PaymentConfig)).Filter("AppNo", c.AppNo).Filter("PayType", utils.PayTypeWechatPay).Limit(1).One(&pc); err != nil {
		c.Msg = "登录异常(WX01)"
		logs.Error("LoginWxApiController QueryTable PaymentConfig err:", err)
		return
	}
	var vo models.WechatVo
	if err := json.Unmarshal([]byte(pc.ConfValue), &vo); err != nil {
		logs.Error("LoginWxApiController Unmarshal ConfValue err:", err)
		c.Msg = "接口异常(WX02)"
		return
	}
	var unionId string
	var nickname string
	var avatarUrl string
	if c.AppChannel == utils.AppChannelWxa { // 微信小程序
		sessionRsp, err := wechat.Code2Session(pc.AppId, vo.AppSecret, p.Code)
		if err != nil {
			logs.Error("LoginWxApiController wechat.Code2Session err:", err)
			c.Msg = "Session获取异常"
			return
		}
		if sessionRsp.Errcode != 0 {
			logs.Error("LoginWxApiController wechat.Code2Session errCode=", sessionRsp.Errcode)
			c.Msg = "授权失败，请重试"
			return
		}
		unionId = sessionRsp.Unionid
	} else {
		// 获取access token
		accessToken, err := wechat.GetOauth2AccessToken(pc.AppId, vo.AppSecret, p.Code)
		//logs.Info("accessToken:", fmt.Sprintf("%+v", accessToken))
		if err != nil {
			logs.Error("LoginWxApiController wechat.GetOauth2AccessToken err:", err)
			c.Msg = "Access Token获取异常"
			return
		} else if accessToken.Errcode > 0 || accessToken.Errmsg != "" {
			logs.Error("LoginWxApiController wechat.GetOauth2AccessToken err:", accessToken.Errmsg)
			c.Msg = "Access Token获取失败"
			return
		} else if accessToken.Unionid == "" {
			c.Msg = "Access Token获取失败"
			return
		}
		// 获取用户信息
		userInfo, err := wechat.GetOauth2UserInfo(accessToken.AccessToken, accessToken.Openid)
		//logs.Info("userInfo:", fmt.Sprintf("%+v", userInfo))
		if err != nil {
			logs.Error("LoginWxApiController wechat.GetUserInfo err:", err)
			c.Msg = "用户信息获取异常"
			return
		} else if userInfo.Unionid == "" {
			c.Msg = "用户信息获取失败"
			return
		}
		unionId = userInfo.Unionid
		nickname = userInfo.Nickname
		avatarUrl = userInfo.Headimgurl
	}
	var member models.Member
	if err := o.QueryTable(new(models.Member)).Filter("ThirdAuthId", unionId).Limit(1).One(&member); err != nil && err != orm.ErrNoRows {
		logs.Error("LoginWxApiController QueryTable Member err:", err)
		c.Msg = "用户查询异常"
		return
	} else if err == orm.ErrNoRows {
		if member, err = CreateMemberReg(2, c.AppNo, c.AppChannel, c.AppVersionCode, 0, unionId, unionId, nickname, unionId, avatarUrl); err != nil {
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

type wxaUserInfoParam struct {
	Code          string `json:"code"`
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`
}

/*
api 小程序提交用户信息，小程序通过用户主动点击按钮授权的方式，获取用户信息
param:
body:
return:
*/
func (c *LoginWxApiController) PostUserInfo() {
	defer c.RetJSON()
	var p wxaUserInfoParam
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
	if err := o.QueryTable(new(models.PaymentConfig)).Filter("AppNo", c.AppNo).Filter("PayType", utils.PayTypeWechatPay).Limit(1).One(&pc); err != nil {
		c.Msg = "异常(WX01)"
		logs.Error("LoginWxApiController QueryTable PaymentConfig err:", err)
		return
	}
	var vo models.WechatVo
	if err := json.Unmarshal([]byte(pc.ConfValue), &vo); err != nil {
		logs.Error("LoginWxApiController Unmarshal ConfValue err:", err)
		c.Msg = "异常(WX02)"
		return
	}
	sessionRsp, err := wechat.Code2Session(pc.AppId, vo.AppSecret, p.Code)
	if err != nil {
		logs.Error("LoginWxApiController.PostUserInfo wechat.Code2Session err:", err)
		c.Msg = "Session获取异常"
		return
	}
	if sessionRsp.Errcode != 0 {
		logs.Error("LoginWxApiController.PostUserInfo wechat.Code2Session errCode=", sessionRsp.Errcode)
		c.Msg = "授权失败，请重试"
		return
	}
	//小程序获取手机号
	//phone := new(wechat.UserPhone)
	//err = wechat.DecryptOpenDataToStruct(p.EncryptedData, p.Iv, sessionRsp.SessionKey, phone)
	//if err != nil {
	//	logs.Error("LoginWxApiController.PostUserInfo wechat.DecryptOpenDataToStruct UserPhone err:", err)
	//	c.Msg = "提交失败E2"
	//	return
	//}
	// 获取微信小程序用户信息
	userInfo := new(wechat.AppletUserInfo)
	err = wechat.DecryptOpenDataToStruct(p.EncryptedData, p.Iv, sessionRsp.SessionKey, userInfo)
	if err != nil {
		logs.Error("LoginWxApiController.PostUserInfo wechat.DecryptOpenDataToStruct AppletUserInfo err:", err)
		c.Msg = "提交失败E3"
		return
	}
	// mobile := phone.PurePhoneNumber
	unionId := sessionRsp.Unionid
	nickname := userInfo.NickName
	avatarUrl := userInfo.AvatarUrl
	var member models.Member
	if err := o.QueryTable(new(models.Member)).Filter("ThirdAuthId", unionId).Limit(1).One(&member); err != nil && err != orm.ErrNoRows {
		logs.Error("LoginWxApiController QueryTable Member err:", err)
		c.Msg = "用户查询异常"
		return
	} else if err == orm.ErrNoRows {
		if member, err = CreateMemberReg(2, c.AppNo, c.AppChannel, c.AppVersionCode, 0, unionId, unionId, nickname, unionId, avatarUrl); err != nil {
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
	c.Msg = "ok"
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

// 获取用户基本信息(UnionID机制)
//    accessToken：接口调用凭据
//    openId：用户的OpenID
//    lang:默认为 zh_CN ，可选填 zh_CN 简体，zh_TW 繁体，en 英语
//    获取用户基本信息(UnionID机制)文档：https://developers.weixin.qq.com/doc/oplatform/Mobile_App/WeChat_Login/Authorized_API_call_UnionID.html
//func getUserInfo(accessToken, openId string, lang ...string) (userInfo *wechat.Oauth2UserInfo, err error) {
//	wechat.GetOauth2UserInfo()
//
//	url := "https://api.weixin.qq.com/sns/userinfo?access_token=" + accessToken + "&openid=" + openId + "&lang=zh_CN"
//	if len(lang) > 0 {
//		url = "https://api.weixin.qq.com/sns/userinfo?access_token=" + accessToken + "&openid=" + openId + "&lang=" + lang[0]
//	}
//	_, errs := gopay.NewHttpClient().Get(url).EndStruct(userInfo)
//	if len(errs) > 0 {
//		return nil, errs[0]
//	}
//	return userInfo, nil
//}
