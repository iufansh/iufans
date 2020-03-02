package sysapi

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	. "github.com/iufansh/iufans/models"
	utils2 "github.com/iufansh/iufans/utils"
	"github.com/iufansh/iutils"
	"strings"
)

type suggestParam struct {
	Suggest string `json:"suggest"`
}

type MemberSuggestApiController struct {
	BaseApiController
}

/*
api查询用户建议
param:
body:
return:{"code":1,"msg":"成功",[{"creatTime":"2020-01-02 15:04:05","suggest":"建议","status":1,"feedback":"很好！"}]}
desc:
*/
func (c *MemberSuggestApiController) Get() {
	defer c.RetJSON()
	o := orm.NewOrm()
	var suggests []MemberSuggest
	if _, err := o.QueryTable(new(MemberSuggest)).Filter("MemberId", c.LoginMemberId).OrderBy("-Id").Limit(10).All(&suggests); err != nil {
		c.Msg = "查看异常"
		return
	}
	suggestList := make([]map[string]interface{}, 0)
	for _, v := range suggests {
		suggestList = append(suggestList, map[string]interface{}{
			"creatTime": iutils.FormatDatetime(v.CreateDate),
			"suggest":   v.Suggest,
			"status":    v.Status,
			"feedback":  v.Feedback,
		})
	}

	c.Code = utils2.CODE_OK
	c.Msg = "查询成功"
	c.Dta = &suggestList
}

/*
api提交用户建议
param:
body: {"suggest":"用户建议内容"}
return:{"code":1,"msg":"成功"}
desc:
*/
func (c *MemberSuggestApiController) Post() {
	defer c.RetJSON()
	var p suggestParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils2.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	if strings.TrimSpace(p.Suggest) == "" {
		c.Msg = "内容不能为空"
		return
	}
	if len(p.Suggest) > 300 {
		c.Msg = "内容不能超过100个字"
		return
	}
	o := orm.NewOrm()
	var member Member
	o.QueryTable(new(Member)).Filter("Id", c.LoginMemberId).One(&member, "OrgId", "Name", "Mobile")
	ms := MemberSuggest{
		OrgId:    member.OrgId,
		MemberId: c.LoginMemberId,
		Name:     member.Name,
		Mobile:   member.Mobile,
		Suggest:  beego.Htmlquote(p.Suggest),
		Status:   0,
	}
	if _, err := o.Insert(&ms); err != nil {
		c.Msg = "提交失败，请重试"
		return
	}
	c.Code = utils2.CODE_OK
	c.Msg = "提交成功"
}
