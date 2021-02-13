package sysapi

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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
	var unReadExist bool
	for _, v := range suggests {
		suggestList = append(suggestList, map[string]interface{}{
			"creatTime": iutils.FormatDatetime(v.CreateDate),
			"suggest":   v.Suggest,
			"status":    v.Status,
			"feedback":  v.Feedback,
		})
		if v.Status == 1 || v.Status == 2 {
			unReadExist = true
		}
	}
	if unReadExist {
		if _, err := o.QueryTable(new(MemberSuggest)).Filter("MemberId", c.LoginMemberId).Filter("Status__in", 1, 2).Update(orm.Params{
			"Status": orm.ColValue(orm.ColAdd, 2),
		}); err != nil {
			logs.Error("MemberSuggestApiController.Get Update MemberSuggest err:", err)
		}
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
	if err := o.QueryTable(new(Member)).Filter("Id", c.LoginMemberId).One(&member, "OrgId", "Name", "Mobile"); err != nil {
		c.Msg = "提交失败，请重试"
		return
	}
	ms := MemberSuggest{
		AppInfo:  fmt.Sprintf("%s-%s-%d", c.AppNo, c.AppChannel, c.AppVersionCode),
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

/*
api 查询用户反馈是否已回复未读
param:
body:
return:{"code":1,"msg":"成功", "data": 1}
desc: data > 0 则存在
*/
func (c *MemberSuggestApiController) GetNewFeedback() {
	defer c.RetJSON()
	o := orm.NewOrm()
	if count, err := o.QueryTable(new(MemberSuggest)).Filter("MemberId", c.LoginMemberId).Filter("Status__in", 1, 2).Count(); err != nil {
		logs.Error("MemberSuggestApiController.getNewFeedback QueryTable MemberSuggest err:", err)
		c.Msg = "查询失败"
		return
	} else {
		c.Dta = count
	}
	c.Code = utils2.CODE_OK
	c.Msg = "查询成功"
}
