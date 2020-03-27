package sysapi

import (
	"encoding/json"
	"github.com/iufansh/iufans/utils"
)

type loginHuaweiParam struct {
	Code  string `json:"code"`
}

type LoginHuaweiApiController struct {
	Base2ApiController
}

/*
api huawei登录
param:
body:{"code":"111122"}
return:{"code":1,"msg":"成功","data":{"id":1,"token":"11111111111111111111","nickname":"微信用户1","autoLogin":true,"accessToken":"ddfesfsf"}}
*/
func (c *LoginHuaweiApiController) Post() {
	defer c.RetJSON()
	var p loginHuaweiParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	if p.Code == "" {
		c.Msg = "Code不能为空"
		return
	}
}
