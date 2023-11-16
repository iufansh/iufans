package sysapi

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	. "github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
	"github.com/iufansh/iuplugins/censor/text"
	. "github.com/iufansh/iutils"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type MemberApiController struct {
	BaseApiController
}

/*
api 获取vip状态
param:
body:
return:{"code":1,"msg":"成功","data":}
desc:
*/
func (c *MemberApiController) GetVip() {
	defer c.RetJSON()
	o := orm.NewOrm()
	var member Member
	if err := o.QueryTable(new(Member)).Filter("Id", c.LoginMemberId).One(&member, "Vip", "VipExpire"); err != nil {
		logs.Error("MemberApiController.getVip query Member err:", err)
		c.Msg = "获取异常，请重试"
		return
	}
	var vipEffect int
	if member.Vip > 0 && !member.VipExpire.IsZero() && TimeSub(member.VipExpire, time.Now()) >= 0 {
		vipEffect = 1
	}
	// 特殊处理
	if vipEffect == 0 && member.VipExpire.IsZero() {
		member.Vip = 0
	}

	c.Dta = map[string]interface{}{
		"vipEffect": vipEffect,                    // vip是否有效
		"vip":       member.Vip,                   // vip等级
		"vipExpire": FormatDate(member.VipExpire), // vip过期时间
	}
	c.Code = utils.CODE_OK
	c.Msg = "ok"
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
	valid.Length(model.Username, 11, "errmsg").Message("手机号必须是11位数字")
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
		c.Msg = "当前手机号已被绑定，请换其他手机号"
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
	c.Dta = model.GetFmtMobile()
}

/*
api已通过其他方式登录的（如微信）且绑定过手机号，通过此接口解除手机号
param:
body:{"authCode":2356,"username":"13111111111","password":"32md5小写"}
return:{"code":1,"msg":"成功","data":"ok"}
desc:
*/
func (c *MemberApiController) UnBindPhone() {
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
	valid.Length(model.Username, 11, "errmsg").Message("手机号必须是11位数字")
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
	var member Member
	if err := o.QueryTable(new(Member)).Filter("Username", model.Username).One(&member); err != nil {
		if err == orm.ErrNoRows {
			c.Msg = "当前手机号不存在，请确认"
			return
		} else {
			c.Msg = "异常，请重试"
			return
		}
	}
	if member.Password != Md5(p.Password, Pubsalt, member.Salt) {
		c.Msg = "密码错误"
		return
	}
	if member.ThirdAuthId == "" {
		c.Msg = "手机号是唯一登录方式，无法解绑"
		return
	}

	model.Password = ""
	model.Salt = ""
	model.Mobile = ""
	model.Username = member.ThirdAuthId
	model.Id = c.LoginMemberId

	if num, err := o.Update(&model, "Username", "Mobile", "Password", "Salt"); err != nil || num != 1 {
		c.Msg = "解绑失败，请重试"
		return
	}

	c.Msg = "解绑成功"
	c.Code = utils.CODE_OK
	c.Dta = ""
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
	member.Name = member.Username
	member.Mobile = ""
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

/*
api 修改昵称
param:
body:
return:{"code":1,"msg":"成功","data":"ok"}
desc:
*/
func (c *MemberApiController) ModifyName() {
	defer c.RetJSON()

	type param struct {
		Name string `json:"name"`
	}
	var p param
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	// 去除所有空格
	p.Name = strings.ReplaceAll(DeleteExtraSpace(strings.TrimSpace(p.Name)), " ", "")
	p.Name = strings.ReplaceAll(p.Name, "<", "")
	p.Name = strings.ReplaceAll(p.Name, ">", "")
	if p.Name == "" {
		p.Name = "用户" + utils.GenInviteCode(c.LoginMemberId)
	}
	if len(p.Name) > 20 {
		c.Code = utils.CODE_ERROR
		c.Msg = "昵称太长了，请重新输入"
		return
	}

	// 验证昵称合规性
	baiduKeyMap := GetSiteConfigMap(utils.ScBaiduApiKey, utils.ScBaiduSecretKey)
	if apiKey, ok := baiduKeyMap[utils.ScBaiduApiKey]; ok && apiKey != "" {
		if secretKey, ok2 := baiduKeyMap[utils.ScBaiduSecretKey]; ok2 && secretKey != "" {
			logs.Debug("MemberApiController 验证昵称合规性")
			var cacKey = fmt.Sprintf("sys_baidu_access_token%s", apiKey)
			var cacToken string
			if err := utils.GetCache(cacKey, &cacToken); err != nil {
				logs.Error("MemberApiController ModifyName getCache err:", err)
			}
			logs.Debug("MemberApiController cache token=", cacToken)
			censorParam := text.CensorParam{
				ApiKey:      apiKey,
				SecretKey:   secretKey,
				AccessToken: cacToken,
				Text:        p.Name,
			}
			isCensorOk, accessToken, err := text.BaiduCensor(censorParam)
			if cacToken != accessToken {
				if err := utils.SetCache(cacKey, accessToken, 2591000); err != nil { // 本来是2592000，提取1000秒失效
					logs.Error("MemberApiController ModifyName SetCache err:", err)
				}
			}
			if err != nil {
				c.Code = utils.CODE_ERROR
				c.Msg = "昵称验证失败，请联系客服"
				return
			} else if !isCensorOk {
				logs.Debug("MemberApiController 昵称 不合规")
				c.Code = utils.CODE_ERROR
				c.Msg = "昵称不合规，请更换"
				return
			}
			logs.Debug("MemberApiController 昵称合规")
		}
	}

	o := orm.NewOrm()
	member := Member{
		Id:   c.LoginMemberId,
		Name: p.Name,
	}
	if num, err := o.Update(&member, "Name"); err != nil || num != 1 {
		c.Msg = "修改失败，请重试"
		return
	}

	c.Msg = "修改成功"
	c.Code = utils.CODE_OK
	c.Dta = p.Name
}

/*
api 修改头像
param:
body:
return:{"code":1,"msg":"成功","data":"/upload/avatar/m1adddd.jpg"}
desc:
*/
func (c *MemberApiController) UploadAvatar() {
	defer c.RetJSONOrigin()

	f, h, err := c.GetFile("file")
	if err != nil {
		logs.Error("MemberApiController upload file get file error", err)
		c.Msg = "上传失败，请重试-1"
		return
	}
	defer f.Close()
	logs.Info("avatarSize=", h.Size)
	if h.Size > 1024*100 {
		c.Msg = "图片太大，不能超过100KB"
		return
	}

	var uploadPath = "upload/member/avatar/"

	if flag, _ := PathExists(uploadPath); !flag {
		if err2 := os.MkdirAll(uploadPath, 0644); err2 != nil {
			logs.Error("MemberApiController upload file MkdirAll error", err2)
			c.Msg = "无法上传"
			return
		}
	}

	fName := url.QueryEscape(h.Filename)
	suffix := SubString(fName, len(fName), strings.LastIndex(fName, ".")-len(fName))
	var saveName = fmt.Sprintf("m%da%d%s", c.LoginMemberId, time.Now().Unix(), suffix)

	uploadName := uploadPath + saveName
	err3 := c.SaveToFile("file", uploadName)
	if err3 != nil {
		logs.Error("MemberApiController upload file save file error2", err3)
		c.Msg = "上传失败，请重试-3"
		return
	}

	o := orm.NewOrm()
	var member Member
	if err := o.QueryTable(new(Member)).Filter("Id", c.LoginMemberId).One(&member); err != nil {
		c.Msg = "上册失败，请重新登录再试"
		return
	}
	if member.Avatar != "" {
		if err := os.Remove(strings.TrimPrefix(member.Avatar, "/")); err != nil {
			logs.Error("MemberApiController UploadAvatar remove avatrar err:", err)
		}
	}

	member.Avatar = "/" + uploadName
	if num, err := o.Update(&member, "Avatar"); err != nil || num != 1 {
		c.Msg = "修改失败，请重试"
		return
	}

	c.Msg = "修改成功"
	c.Code = utils.CODE_OK
	c.Dta = member.GetFullAvatar(c.Ctx.Input.Site())
}
