package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/iufansh/iufans/utils"
	"net/http"
	"strings"
)

type BaseMainController struct {
	beego.Controller
	// 返回json时用，最后必须调用方法RetJSON或RetJSONOrigin
	Code int
	// 返回Tpl时用，最后必须调用方法
	// 返回json时用，最后必须调用方法RetJSON或RetJSONOrigin
	Msg  string
	// 返回json时用，最后必须调用方法RetJSON或RetJSONOrigin
	Dta  interface{} // 本来叫Data,但是和beego.Controller的冲突，所以改为Dta
}

func (c *BaseMainController) RetTpl(tplName string) {
	c.Data["code"] = c.Code
	c.Data["msg"] = c.Msg
	c.Data["data"] = c.Dta
	c.TplName = tplName
}

func (c *BaseMainController) RetJSON() {
	ret := make(map[string]interface{})
	ret["code"] = c.Code
	ret["msg"] = c.Msg
	switch c.Dta.(type) {
	case string:
		if strings.HasPrefix(c.Dta.(string), "/") {
			ret["url"] = c.Dta
		}
		break
	case *string:
		if strings.HasPrefix(c.Dta.(string), "/") {
			ret["url"] = c.Dta
		}
		break
	}
	if c.Code == utils.CODE_OK && c.Dta == nil {
		ret["data"] = "ok"
	} else {
		ret["data"] = c.Dta
	}
	//b, _ := json.Marshal(ret)

	//logs.Debug("\r\n----------response---------",
		//"\r\n", iutils.SubString(string(b), 0, 300),
		//"\r\n", string(b),
	//	"\r\n-------------------------",)

	c.Data["json"] = ret
	c.ServeJSON()
}

func (c *BaseMainController) RetJSONOrigin() {
	ret := make(map[string]interface{})
	ret["code"] = c.Code
	ret["msg"] = c.Msg
	if c.Code == utils.CODE_OK && c.Dta == nil {
		ret["data"] = "ok"
	} else {
		ret["data"] = c.Dta
	}

	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")

	var err error
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	err = jsonEncoder.Encode(ret)

	if err != nil {
		http.Error(c.Ctx.Output.Context.ResponseWriter, err.Error(), http.StatusInternalServerError)
		logs.Error("BaseMainController marshal err:", err)
		return
	}

	//logs.Debug("\r\n----------response origin---------",
	//	"\r\n", iutils.SubString(bf.String(), 0, 300),
		//"\r\n", bf.String(),
	//	"\r\n-------------------------",)

	c.Ctx.Output.Body(bf.Bytes())
}