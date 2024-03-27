package sysapi

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	. "github.com/iufansh/iufans/models"
	utils2 "github.com/iufansh/iufans/utils"
	"github.com/iufansh/iutils"
)

type NormalQuestionApiController struct {
	Base2ApiController
}

/*
api 常见问题获取
param:
body:
return:{"code":1,"msg":"成功","data":[{"seq":1,"question":"问题1","answer":"回答1"}]}
desc:
*/
func (c *NormalQuestionApiController) Get() {
	defer c.RetJSON()

	// forbiddenArea := GetSiteConfigValue(utils.ScApiIpForbidden)
	// if allowed := iutils.CheckIpAllowed(forbiddenArea, c.Ctx.Input.IP()); !allowed {
	// 	c.Code = utils2.CODE_OK
	// 	c.Msg = "没有内容"
	// 	c.Dta = []map[string]interface{}{}
	// 	return
	// }

	o := orm.NewOrm()
	qs := o.QueryTable(new(NormalQuestion)).OrderBy("Seq").Limit(10)
	var quetions []NormalQuestion
	if num, err := qs.All(&quetions); err != nil {
		logs.Error("NormalQuestionApiController.query NormalQuestion err:", err)
		c.Msg = "查看失败，请重试"
		return
	} else if num == 0 {
		c.Code = utils2.CODE_OK
		c.Msg = "没有常见问题"
		c.Dta = []map[string]interface{}{}
		return
	}
	list := make([]map[string]interface{}, 0)
	for i, v := range quetions {
		list = append(list, map[string]interface{}{
			"seq":      i+1,
			"question": v.Question,
			"answer":   v.Answer,
		})
	}
	c.Dta = &list
	c.Code = utils2.CODE_OK
	c.Msg = "查询成功"
}
