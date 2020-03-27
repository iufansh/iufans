package paymentconfig

import (
	. "github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/controllers/sysmanage"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/iufansh/iufans/utils"
	"encoding/json"
	"fmt"
	"html/template"
)

func validatePaymentConfig(model *PaymentConfig) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(model.AppNo, "errmsg").Message("应用编号必填")
	valid.Required(model.AppName, "errmsg").Message("应用名称必填")
	valid.MaxSize(model.AppName, 100, "errmsg").Message("应用名称最长100个字符")
	valid.Required(model.AppId, "errmsg").Message("AppId必填")
	valid.MaxSize(model.AppId, 50, "errmsg").Message("AppId最长50个字符")
	valid.MaxSize(model.Remark, 255, "errmsg").Message("备注最长255个字符")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type PaymentConfigIndexController struct {
	sysmanage.BaseController
}

func (c *PaymentConfigIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *PaymentConfigIndexController) Get() {
	o := orm.NewOrm()
	var list []PaymentConfig
	o.QueryTable(new(PaymentConfig)).OrderBy("-Id").All(&list)
	// 返回值
	c.Data["dataList"] = &list

	if t, err := template.New("tplIndexPaymentConfig.tpl").Funcs(map[string]interface{}{
		"date": beego.Date,
		"urlfor": beego.URLFor,
	}).Parse(tplIndex); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *PaymentConfigIndexController) Enabled() {
	defer c.RetJSON()
	id, _ := c.GetInt64("id")
	model := PaymentConfig{Id: id}
	o := orm.NewOrm()
	err := o.Read(&model)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Msg = "数据不存在，请确认"
		return
	}
	if model.Enabled == 1 {
		model.Enabled = 0
	} else {
		model.Enabled = 1
	}
	model.Modifior = c.LoginAdminId

	if _, err := o.Update(&model, "Enabled", "ModifyDate", "Modifior"); err != nil {
		beego.Error("Enabled PaymentConfig error", err)
		c.Msg = "操作失败"
	} else {
		c.Code = 1
		c.Msg = "操作成功"
		c.Dta = c.URLFor("PaymentConfigIndexController.Get")
		wLog := fmt.Sprintf("Action:%s, Id:%d, enable:%d", "PaymentConfig Enabled", model.Id, model.Enabled)
		beego.Warn(wLog, "User:", c.LoginAdminUsername, "Ip:", c.Ctx.Input.IP())
	}
}

func (c *PaymentConfigIndexController) Delone() {
	defer c.RetJSON()
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	if _, err := o.Delete(&PaymentConfig{Id: id}, "Id"); err != nil {
		beego.Error("Delete PaymentConfig error", err)
		c.Msg = "删除失败"
	} else {
		c.Code = 1
		c.Msg = "删除成功"
		c.Dta = c.URLFor("PaymentConfigIndexController.Get")
		wLog := fmt.Sprintf("Action:%s, Id:%d", "PaymentConfig Delone", id)
		beego.Warn(wLog, "User:", c.LoginAdminUsername, "Ip:", c.Ctx.Input.IP())
	}
}

type PaymentConfigAddController struct {
	sysmanage.BaseController
}

func (c *PaymentConfigAddController) NestPrepare() {
	c.EnableRender = false
}
func (c *PaymentConfigAddController) Get() {
	payType := c.GetString("pt", utils.PayTypeWechatPay)

	c.Data["data"] = PaymentConfig{PayType: payType}
	if payType == utils.PayTypeWechatPay {
		c.Data["vo"] = WechatVo{}
		if t, err := template.New("tplAddWechatpay.tpl").Funcs(map[string]interface{}{
			"urlfor": beego.URLFor,
		}).Parse(tplAddWechatpay); err != nil {
			beego.Error("template Parse err", err)
		} else {
			t.Execute(c.Ctx.ResponseWriter, c.Data)
		}
	} else {
		c.Data["vo"] = AlipayVo{}
		if t, err := template.New("tplAddAlipay.tpl").Funcs(map[string]interface{}{
			"urlfor": beego.URLFor,
		}).Parse(tplAddAlipay); err != nil {
			beego.Error("template Parse err", err)
		} else {
			t.Execute(c.Ctx.ResponseWriter, c.Data)
		}
	}
}

func (c *PaymentConfigAddController) Post() {
	defer c.RetJSON()
	model := PaymentConfig{}
	if err := c.ParseForm(&model); err != nil {
		c.Msg = "参数异常"
		return
	} else if hasError, errMsg := validatePaymentConfig(&model); hasError {
		c.Msg = errMsg
		return
	}
	// 配置信息
	if model.PayType == utils.PayTypeWechatPay {
		var vo = WechatVo{}
		if err := c.ParseForm(&vo); err != nil {
			beego.Error("PaymentConfig parse config value vo err", err)
			c.Msg = "参数异常(2)"
			return
		}
		if bytes, err := json.Marshal(&vo); err != nil {
			c.Msg = "参数异常(3)"
			return
		} else {
			model.ConfValue = string(bytes)
		}
	} else if model.PayType == utils.PayTypeAlipay {
		var vo = AlipayVo{}
		if err := c.ParseForm(&vo); err != nil {
			beego.Error("PaymentConfig parse config value vo err", err)
			c.Msg = "参数异常(2)"
			return
		}
		if bytes, err := json.Marshal(&vo); err != nil {
			c.Msg = "参数异常(3)"
			return
		} else {
			model.ConfValue = string(bytes)
		}
	}
	o := orm.NewOrm()

	model.OrgId = c.LoginAdminOrgId
	model.Creator = c.LoginAdminId
	model.Modifior = c.LoginAdminId
	model.Version = 1
	var wLog string

	if id, err := o.Insert(&model); err != nil {
		c.Msg = "添加失败"
		beego.Error("Insert PaymentConfig error", err)
		return
	} else {
		c.Code = 1
		c.Msg = "添加成功"
		wLog = fmt.Sprintf("Action:%s, Id:%d, Model:%+v", "PaymentConfig Add", id, model)
	}

	c.Dta = c.URLFor("PaymentConfigIndexController.Get")
	beego.Warn(wLog, "User:", c.LoginAdminUsername, "Ip:", c.Ctx.Input.IP())
}

type PaymentConfigEditController struct {
	sysmanage.BaseController
}

func (c *PaymentConfigEditController) NestPrepare() {
	c.EnableRender = false
}
func (c *PaymentConfigEditController) Get() {
	id, _ := c.GetInt64("id")

	o := orm.NewOrm()
	model := PaymentConfig{Id: id}

	err := o.Read(&model)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(beego.URLFor("PaymentConfigIndexController.get"), 302)
	} else {
		var vo interface{}
		json.Unmarshal([]byte(model.ConfValue), &vo)

		c.Data["data"] = &model
		c.Data["vo"] = vo
		if model.PayType == utils.PayTypeWechatPay {
			if t, err := template.New("tplEditWechatpay.tpl").Funcs(map[string]interface{}{
				"urlfor": beego.URLFor,
			}).Parse(tplEditWechatpay); err != nil {
				beego.Error("template Parse err", err)
			} else {
				t.Execute(c.Ctx.ResponseWriter, c.Data)
			}
		} else {
			if t, err := template.New("tplEditAlipay.tpl").Funcs(map[string]interface{}{
				"urlfor": beego.URLFor,
			}).Parse(tplEditAlipay); err != nil {
				beego.Error("template Parse err", err)
			} else {
				t.Execute(c.Ctx.ResponseWriter, c.Data)
			}
		}
	}
}

func (c *PaymentConfigEditController) Post() {
	defer c.RetJSON()
	model := PaymentConfig{}
	if err := c.ParseForm(&model); err != nil {
		c.Msg = "参数异常"
		return
	} else if hasError, errMsg := validatePaymentConfig(&model); hasError {
		c.Msg = errMsg
		return
	}
	// 配置信息
	if model.PayType == utils.PayTypeWechatPay {
		var vo = WechatVo{}
		if err := c.ParseForm(&vo); err != nil {
			beego.Error("PaymentConfig parse config value vo err", err)
			c.Msg = "参数异常(2)"
			return
		}
		if bytes, err := json.Marshal(&vo); err != nil {
			c.Msg = "参数异常(3)"
			return
		} else {
			model.ConfValue = string(bytes)
		}
	} else if model.PayType == utils.PayTypeAlipay {
		var vo = AlipayVo{}
		if err := c.ParseForm(&vo); err != nil {
			beego.Error("PaymentConfig parse config value vo err", err)
			c.Msg = "参数异常(2)"
			return
		}
		if bytes, err := json.Marshal(&vo); err != nil {
			c.Msg = "参数异常(3)"
			return
		} else {
			model.ConfValue = string(bytes)
		}
	}
	o := orm.NewOrm()
	model.Modifior = c.LoginAdminId
	if _, err := o.Update(&model, "AppId", "AppNo", "AppName", "ConfValue", "Enabled", "Remark", "ModifyDate"); err != nil {
		c.Msg = "更新失败"
		beego.Error("Update PaymentConfig error", err)
	} else {
		c.Code = 1
		c.Msg = "更新成功"
		c.Dta = c.URLFor("PaymentConfigIndexController.Get")
		wLog := fmt.Sprintf("Action:%s, Id:%d, Model:%+v", "PaymentConfig Edit", model.Id, model)
		beego.Warn(wLog, "User:", c.LoginAdminUsername, "Ip:", c.Ctx.Input.IP())
	}
}
