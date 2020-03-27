package controllers

import (
	"bytes"
	"github.com/astaxie/beego"
	"encoding/json"
	"github.com/iufansh/iufans/utils"
	"github.com/iufansh/iutils"
	"net/http"
	"strings"
)

type BaseMainController struct {
	beego.Controller
	Code int
	Msg  string
	Dta  interface{} // 本来叫Data,但是和beego.Controller的冲突，所以改为Dta
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
	b, _ := json.Marshal(ret)

	beego.Info("\r\n----------response---------",
		"\r\n", iutils.SubString(string(b), 0, 300),
		"\r\n-------------------------",)

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
		beego.Error("BaseMainController marshal err:", err)
		return
	}

	beego.Info("\r\n----------response origin---------",
		"\r\n", iutils.SubString(bf.String(), 0, 300),
		"\r\n-------------------------",)

	c.Ctx.Output.Body(bf.Bytes())
}