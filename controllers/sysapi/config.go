package sysapi

import (
	fm "github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
	"strings"
)

type SysConfigApiController struct {
	Base2ApiController
}

/*
api 系统配置信息获取
param: c=tel,qq
body:
return:{"code":1,"msg":"ok","data":{"tel":"13111111111","qq":"2321353"}}
desc: 参数可以输入多个，用逗号隔开
*/
func (c *SysConfigApiController) Get() {
	defer c.RetJSON()
	p := c.GetString("c")
	codes := make([]string, 0)
	if strings.Contains(p, "tel") {
		codes = append(codes, utils.Sccompanyconcattel)
	}
	if strings.Contains(p, "qq") {
		codes = append(codes, utils.Sccompanyconcatqq)
	}
	if len(codes) == 0 {
		c.Msg = "参数未知"
		return
	}
	m := fm.GetSiteConfigMap(codes)
	retMap := map[string]string{}
	if strings.Contains(p, "tel") {
		retMap["tel"] = m[utils.Sccompanyconcattel]
	}
	if strings.Contains(p, "qq") {
		retMap["qq"] = m[utils.Sccompanyconcatqq]
	}
	c.Code = utils.CODE_OK
	c.Msg = "ok"
	c.Dta = retMap
}
