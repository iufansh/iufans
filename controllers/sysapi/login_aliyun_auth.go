package sysapi

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
	"github.com/iufansh/iutils"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type loginAliyunAuthParam struct {
	Token string `json:"token"`
}

type LoginAliyunAuthApiController struct {
	Base2ApiController
}

/*
api Aliyun认证登录，注册
param:
body:
return:{"code":1,"msg":"成功","data":"authInfo"}
*/
func (c *LoginAliyunAuthApiController) Post() {
	defer c.RetJSON()
	var p loginAliyunAuthParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	if p.Token == "" {
		c.Msg = "Token不能为空"
		return
	}
	regionId := beego.AppConfig.String("aliyunregionid")
	if regionId == "" {
		regionId = "cn-hangzhou"
	}
	scs := models.GetSiteConfigMap(utils.ScAliyunAccessKeyId, utils.ScAliyunAccessKeySecret)
	client, err := sdk.NewClientWithAccessKey(regionId, scs[utils.ScAliyunAccessKeyId], scs[utils.ScAliyunAccessKeySecret])
	if err != nil {
		c.Msg = "认证失败(E1)"
		logs.Error("LoginAliyunAuthApiController NewClientWithAccessKey err:", err)
		return
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dypnsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "GetMobile"
	request.QueryParams["RegionId"] = regionId
	request.QueryParams["AccessToken"] = p.Token

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		c.Msg = "认证失败(E2)"
		logs.Error("LoginAliyunAuthApiController ProcessCommonRequest err:", err)
		return
	}
	// {"Message":"OK","RequestId":"CC702B1A-524C-42CC-A674-09CAE4260277","Code":"OK","GetMobileResultDTO":{"Mobile":"13711111111"}}
	// fmt.Print(response.GetHttpContentString())
	var respJson aliyunAuthResponse
	if err := json.Unmarshal([]byte(response.GetHttpContentString()), &respJson); err != nil {
		c.Msg = "认证失败(ER1)"
		return
	}
	if respJson.Code != "OK" {
		c.Msg = "认证失败(ER2)"
		return
	}
	o := orm.NewOrm()
	var member models.Member
	if err := o.QueryTable(new(models.Member)).Filter("Username", respJson.GetMobileResultDTO.Mobile).Limit(1).One(&member); err != nil && err != orm.ErrNoRows {
		logs.Error("LoginAliyunAuthApiController QueryTable Member err:", err)
		c.Msg = "用户查询异常，请重试"
		return
	} else if err == orm.ErrNoRows {
		if member, err = CreateMemberReg(5, c.AppNo, c.AppChannel, c.AppVersionCode, 0, respJson.GetMobileResultDTO.Mobile, respJson.RequestId, "", "", ""); err != nil {
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

type aliyunAuthMobile struct {
	Mobile string `json:"Mobile"`
}

type aliyunAuthResponse struct {
	Code string `json:"Code"`
	Message string `json:"Message"`
	RequestId string `json:"RequestId"`
	GetMobileResultDTO aliyunAuthMobile `json:"GetMobileResultDTO"`
}