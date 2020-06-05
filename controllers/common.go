package controllers

import (
	fu "github.com/iufansh/iufans/utils"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/iufansh/iufans/models"
	license "github.com/iufansh/iuplugins/iu-authorize"
)

type CommonController struct {
	beego.Controller
}

// 发送短信验证码，front和api使用
// 请求：get
// 参数：mobile=手机号&mode=1
// mode=1验证手机号不存在则返回异常；mode=2不验证手机号是否存在；mode=3验证手机号已存在则返回异常
func (c *CommonController) SendSmsCode() {
	beego.Info("\r\n----------request---------",
		"\r\nUri:", c.Ctx.Input.URI(),
		"\r\nMethod:", c.Ctx.Input.Method(),
		"\r\nFrom ip:", c.Ctx.Input.IP(),
		"\r\nUserAgent:", c.Ctx.Input.UserAgent(),
		"\r\nBody:", string(c.Ctx.Input.RequestBody),
		"\r\n--------------------------")
	var code int
	var msg string
	defer func(msg *string, code *int) {
		ret := make(map[string]interface{})
		ret["code"] = code
		ret["msg"] = msg
		ret["data"] = ""
		c.Data["json"] = ret
		c.ServeJSON()
	}(&msg, &code)
	to := c.GetString("mobile")
	if to == "" {
		msg = "手机号不能为空"
		return
	}
	mode, _ := c.GetInt("mode", 1)
	if mode == 1 || mode == 3 { // 验证手机号是否在系统中
		o := orm.NewOrm()
		if exist := o.QueryTable(new(models.Member)).Filter("Username", to).Exist(); mode == 1 && !exist {
			msg = "手机号不存在"
			return
		} else if mode == 3 && exist {
			msg = "手机号已存在"
			return
		}
	}
	sc := models.GetSiteConfigMap(fu.Scsmssignname, fu.Scsmsapi, fu.Scsmsuid, fu.Scsmskey, fu.Scname)
	var companyName string
	if sc[fu.Scsmssignname] != "" {
		companyName = sc[fu.Scsmssignname]
	} else {
		companyName = sc[fu.Scname]
	}
	ms := fu.SmsSender{
		Api:     sc[fu.Scsmsapi],
		Uid:     sc[fu.Scsmsuid],
		Key:     sc[fu.Scsmskey],
		Mobile:  to,
		Company: companyName,
	}
	msg, err := fu.SendSmsVerifyCode(ms)
	if err != nil {
		msg = err.Error()
	} else {
		code = 1
	}
}

func (c *CommonController) HealthCheck() {
	c.Ctx.ResponseWriter.Write([]byte("1"))
}

func (c *CommonController) SystemInfo() {
	var code int
	var msg string
	ret := make(map[string]interface{})
	token := c.GetString("t")
	if token != "" {
		t := time.Now().Format("20060102")
		if token == t {
			code = 1
			ret["data"] = license.GetMachineData()
			ret["code"] = code
			ret["msg"] = msg
			c.Data["json"] = ret
			c.ServeJSON()
			return
		}
	}
	c.Abort("404")
}
