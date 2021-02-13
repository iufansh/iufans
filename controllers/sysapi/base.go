package sysapi

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/controllers"
	. "github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
	"strconv"
	"strings"
	"time"
)

type ApiPreparer interface {
	ApiPreparer()
}

/*
api需要登录的接口父接口
param: mid=1&token=0887a57c630a5a97ee8ded7141754724
desc: 继承该接口的将验证登录状态
*/
type BaseApiController struct {
	controllers.BaseMainController
	LoginMemberId  int64
	AppChannel     string
	AppNo          string
	AppVersionCode int
	DeviceId       string
}

func (c *BaseApiController) Prepare() {
	c.EnableXSRF = false
	logs.Info("\r\n----------request---------",
		"\r\nUri:", c.Ctx.Input.URI(),
		"\r\nMethod:", c.Ctx.Input.Method(),
		"\r\nFrom ip:", c.Ctx.Input.IP(),
		"\r\nUserAgent:", c.Ctx.Input.UserAgent(),
		"\r\nReferer:", c.Ctx.Input.Referer(),
		"\r\nHeader Qx-Api-App:", c.Ctx.Input.Header("Qx-Api-App"),
		"\r\nHeader Qx-Api-Token:", c.Ctx.Input.Header("Qx-Api-Token"),
		"\r\nBody:", string(c.Ctx.Input.RequestBody),
		"\r\n--------------------------")

	// 格式：AppNo;appChannel;deviceId;appVersionCode
	apiApp := c.Ctx.Input.Header("Qx-Api-App")
	// 格式：用户id;token
	apiToken := c.Ctx.Input.Header("Qx-Api-Token")

	var err error
	if apiApp != "" {
		strs := strings.Split(apiApp, ";")
		if len(strs) > 0 {
			c.AppNo = strs[0]
		}
		if len(strs) > 1 {
			c.AppChannel = strs[1]
		}
		if len(strs) > 2 {
			c.DeviceId = strs[2]
		}
		if len(strs) > 3 {
			versionCode, _ := strconv.ParseInt(strs[3], 10, 64)
			c.AppVersionCode = int(versionCode)
		}
	}
	var token string
	if apiToken != "" {
		strs := strings.Split(apiToken, ";")
		/*
			if len(strs) < 2 {
				c.Code = utils.CODE_NEED_LOGIN
				c.Msg = "Api参数异常"
				c.RetJSON()
				return
			}

		*/
		if len(strs) > 0 {
			c.LoginMemberId, err = strconv.ParseInt(strs[0], 10, 64)
		}
		if len(strs) > 1 {
			token = strs[1]
		}
		// TODO 遗留，要删除
		if len(strs) > 2 {
			c.AppChannel = strs[2]
		}
		// TODO 遗留，要删除
		if len(strs) > 3 {
			c.AppNo = strs[3]
		}
	} else {
		c.LoginMemberId, err = c.GetInt64("mid", 0)
		token = c.GetString("token", "")
	}
	if err != nil || c.LoginMemberId <= 0 {
		c.Code = utils.CODE_NEED_LOGIN
		c.Msg = "用户ID为空"
		c.RetJSON()
		return
	}
	if token == "" {
		c.Code = utils.CODE_NEED_LOGIN
		c.Msg = "token为空"
		c.RetJSON()
		return
	}
	o := orm.NewOrm()
	var member Member
	if err := o.QueryTable(new(Member)).Filter("Id", c.LoginMemberId).One(&member, "Enabled", "Locked", "Token", "TokenExpTime"); err != nil {
		if err == orm.ErrNoRows {
			c.Code = utils.CODE_NEED_LOGIN
			c.Msg = "用户不存在"
			c.RetJSON()
			return
		} else {
			logs.Error("BaseApiController QueryTable Member err:", err)
			c.Msg = "异常，请重试"
			c.RetJSON()
			return
		}
	}
	if member.Enabled != 1 {
		c.Msg = "用户已禁用"
		c.RetJSON()
		return
	}
	if member.Locked == 1 {
		c.Msg = "用户已锁定"
		c.RetJSON()
		return
	}
	if token != member.Token {
		c.Code = utils.CODE_NEED_LOGIN
		c.Msg = "登录验证失败，请重新登录"
		c.RetJSON()
		return
	}
	if time.Now().After(member.TokenExpTime) {
		c.Code = utils.CODE_NEED_LOGIN
		c.Msg = "登录过期，请重新登录"
		c.RetJSON()
		return
	}

	if app, ok := c.AppController.(ApiPreparer); ok {
		app.ApiPreparer()
	}
}

/*
api不需要登录的父接口
param: mid=1 (可不传)
desc: 如果mid参数不为空，则记录为已登录状态
*/
type Base2ApiController struct {
	controllers.BaseMainController
	IsLogon        bool
	LoginMemberId  int64
	AppChannel     string
	AppNo          string
	AppVersionCode int
	DeviceId       string
}

func (c *Base2ApiController) Prepare() {
	c.EnableXSRF = false
	logs.Info("\r\n----------request---------",
		"\r\nUri:", c.Ctx.Input.URI(),
		"\r\nMethod:", c.Ctx.Input.Method(),
		"\r\nFrom ip:", c.Ctx.Input.IP(),
		"\r\nUserAgent:", c.Ctx.Input.UserAgent(),
		"\r\nReferer:", c.Ctx.Input.Referer(),
		"\r\nHeader Qx-Api-App:", c.Ctx.Input.Header("Qx-Api-App"),
		"\r\nHeader Qx-Api-Token:", c.Ctx.Input.Header("Qx-Api-Token"),
		"\r\nBody:", string(c.Ctx.Input.RequestBody),
		"\r\n--------------------------")

	// 格式：AppNo;appChannel;deviceId;appVersionCode
	apiApp := c.Ctx.Input.Header("Qx-Api-App")
	// 格式：用户id;token
	apiToken := c.Ctx.Input.Header("Qx-Api-Token")

	var err error
	if apiApp != "" {
		strs := strings.Split(apiApp, ";")
		if len(strs) > 0 {
			c.AppNo = strs[0]
		}
		if len(strs) > 1 {
			c.AppChannel = strs[1]
		}
		if len(strs) > 2 {
			c.DeviceId = strs[2]
		}
		if len(strs) > 3 {
			versionCode, _ := strconv.ParseInt(strs[3], 10, 64)
			c.AppVersionCode = int(versionCode)
		}
	}
	var token string
	if apiToken != "" {
		strs := strings.Split(apiToken, ";")
		if len(strs) > 0 {
			c.LoginMemberId, err = strconv.ParseInt(strs[0], 10, 64)
		}
		if len(strs) > 1 {
			token = strs[1]
		}
		// TODO 遗留，要删除
		if len(strs) > 2 {
			c.AppChannel = strs[2]
		}
		// TODO 遗留，要删除
		if len(strs) > 3 {
			c.AppNo = strs[3]
		}
	} else {
		c.LoginMemberId, err = c.GetInt64("mid", 0)
	}
	if err == nil && c.LoginMemberId > 0 && token != "" {
		c.IsLogon = true
	}
	if app, ok := c.AppController.(ApiPreparer); ok {
		app.ApiPreparer()
	}
}

/*
func Retjson(ctx *context.Context, msg *string, code *int, data ...interface{}) {
	ret := make(map[string]interface{})
	ret["code"] = code
	ret["msg"] = msg
	if len(data) > 0 {
		ret["data"] = data[0]
	}
	b, _ := json.Marshal(ret)
	logs.Info("\r\n----------response---------",
		"\r\n", string(b),
		"\r\n-------------------------",)
	ctx.Output.JSON(ret, false, false)
}
*/
