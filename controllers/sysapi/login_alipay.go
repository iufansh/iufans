package sysapi

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/alipay"
	"github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
	"github.com/iufansh/iutils"
	"hash"
	"strings"
)

type loginAlipayParam struct {
	Code string `json:"code"`
}

type LoginAlipayApiController struct {
	Base2ApiController
}

/*
api Alipay登录，获取authInfo
param:
body:
return:{"code":1,"msg":"成功","data":"authInfo"}
*/
func (c *LoginAlipayApiController) Get() {
	defer c.RetJSON()
	o := orm.NewOrm()
	var pc models.PaymentConfig
	o.QueryTable(new(models.PaymentConfig)).Filter("AppNo", c.AppNo).Filter("PayType", utils.PayTypeAlipay).Limit(1).One(&pc)
	var vo models.AlipayVo
	if err := json.Unmarshal([]byte(pc.ConfValue), &vo); err != nil {
		logs.Error("Unmarshal ConfValue err:", err)
		c.Msg = "接口异常(ALI)"
		return
	}
	bm := make(gopay.BodyMap)
	bm.Set("app_id", pc.AppId)
	bm.Set("pid", vo.PartnerId)
	bm.Set("apiname", "com.alipay.account.auth")
	bm.Set("methodname", "alipay.open.auth.sdk.code.get")
	bm.Set("app_name", "mc")
	bm.Set("biz_type", "openservice")
	bm.Set("product_id", "APP_FAST_LOGIN")
	bm.Set("scope", "kuaijie")
	bm.Set("target_id", c.AppNo+c.AppChannel+iutils.RandStringLower(20))
	bm.Set("auth_type", "AUTHACCOUNT")
	bm.Set("sign_type", "RSA2")

	urlParam := alipay.FormatURLParam(bm)
	priKey := alipay.FormatPrivateKey(vo.PriKey)
	sign, err := getRsaSign(bm, "RSA2", priKey)
	if err != nil {
		logs.Error("LoginAlipayApiController.Get getRsaSign err:", err)
		c.Msg = "获取签名失败"
		return
	}
	logs.Info("LoginAlipayApiController.Get authInfo=", urlParam+"&"+sign)
	c.Code = utils.CODE_OK
	c.Msg = "获取成功"
	c.Dta = urlParam + "&" + sign
}

// 获取参数签名
// copy 自alipay.param
func getRsaSign(bm gopay.BodyMap, signType, privateKey string) (sign string, err error) {
	var (
		block          *pem.Block
		h              hash.Hash
		key            *rsa.PrivateKey
		hashs          crypto.Hash
		encryptedBytes []byte
	)

	if block, _ = pem.Decode([]byte(privateKey)); block == nil {
		return gopay.NULL, errors.New("pem.Decode：privateKey decode error")
	}
	if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		return
	}
	switch signType {
	case "RSA":
		h = sha1.New()
		hashs = crypto.SHA1
	case "RSA2":
		h = sha256.New()
		hashs = crypto.SHA256
	default:
		h = sha256.New()
		hashs = crypto.SHA256
	}
	if _, err = h.Write([]byte(bm.EncodeAliPaySignParams())); err != nil {
		return
	}
	if encryptedBytes, err = rsa.SignPKCS1v15(rand.Reader, key, hashs, h.Sum(nil)); err != nil {
		return
	}
	sign = base64.StdEncoding.EncodeToString(encryptedBytes)
	return
}

/*
api Alipay登录，获取用户基本信息
param:
body:{"code":"111122"}
return:{"code":1,"msg":"成功","data":{"id":1,"token":"11111111111111111111","nickname":"支付宝用户1","autoLogin":true,"accessToken":"ddfesfsf"}}
*/
func (c *LoginAlipayApiController) Post() {
	defer c.RetJSON()
	var p loginAlipayParam
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
	o.QueryTable(new(models.PaymentConfig)).Filter("AppNo", c.AppNo).Filter("PayType", utils.PayTypeAlipay).Limit(1).One(&pc)
	var vo models.AlipayVo
	if err := json.Unmarshal([]byte(pc.ConfValue), &vo); err != nil {
		logs.Error("Unmarshal ConfValue err:", err)
		c.Msg = "接口异常(ALI)"
		return
	}
	// 获取access token
	rsp, err := alipay.SystemOauthToken(pc.AppId, vo.PriKey, "authorization_code", p.Code, "RSA2")
	logs.Info("alipay.SystemOauthToken rsp=", fmt.Sprintf("%+v", rsp))
	if err != nil {
		logs.Error("alipay.SystemOauthToken err:", err)
		c.Msg = "授权异常(ALI1)"
		return
	} else if rsp.ErrorResponse != nil && rsp.ErrorResponse.Code != "10000" {
		logs.Error("alipay.SystemOauthToken err:", rsp.ErrorResponse)
		c.Msg = "信息获取失败(ALI3)"
		return
	} else if rsp.Response == nil || rsp.Response.UserId == "" {
		logs.Error("alipay.SystemOauthToken err:", rsp.ErrorResponse)
		c.Msg = "信息获取失败(ALI4)"
		return
	}
	var member models.Member
	if err := o.QueryTable(new(models.Member)).Filter("ThirdAuthId", rsp.Response.UserId).Limit(1).One(&member); err != nil && err != orm.ErrNoRows {
		logs.Error("QueryTable Member err:", err)
		c.Msg = "用户查询异常"
		return
	} else if err == orm.ErrNoRows {
		var isProd bool
		if beego.BConfig.RunMode == "prod" {
			isProd = true
		}
		client := alipay.NewClient(pc.AppId, vo.PriKey, isProd)
		client.SetCharset("utf-8").SetSignType("RSA2").SetAuthToken(rsp.Response.AccessToken)
		resp, err := client.UserInfoShare()
		logs.Info("alipay.UserInfoShare rsp=", resp)
		if err != nil {
			logs.Error("alipay.UserInfoShare err:", err)
			c.Msg = "授权异常(ALI2)"
			return
		} else if resp.Response.Code != "10000" {
			logs.Error("alipay.UserInfoShare err:", resp.Response)
			c.Msg = "信息获取失败(ALI3)"
			return
		}
		var nickName string
		if resp.Response.NickName == "" || strings.HasPrefix(resp.Response.NickName, "2088") {
			nickName = "2088***" + iutils.SubString(resp.Response.UserId, len(resp.Response.UserId)-4, 4)
		} else {
			nickName = resp.Response.NickName
		}
		if member, err = CreateMemberReg(c.AppNo, c.AppChannel, c.AppVersionCode, 0, resp.Response.UserId, resp.Response.UserId, nickName, resp.Response.UserId, resp.Response.Avatar); err != nil {
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
		// "accessToken": rsp.Response.AccessToken, // access token
	}
}
