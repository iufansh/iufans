package sysapi

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	. "github.com/iufansh/iufans/models"
	utils2 "github.com/iufansh/iufans/utils"
	"strconv"
	"time"
)

type InformationApiController struct {
	Base2ApiController
}

/*
api 系统消息获取
param:
body:
return:{"code":1,"msg":"成功","data":[{"infoId":1,"title":"系统通知","info":"系统通知消息内容","needFeedback":1}]}
desc:
*/
func (c *InformationApiController) Get() {
	defer c.RetJSON()

	receiver := []string{
		"::",                     // 全部空，表示所有人
		c.AppNo + "::",           // App编号为自己的
		":" + c.AppChannel + ":", // App渠道是自己的
		"::" + strconv.FormatInt(c.LoginMemberId, 10), // ID为自己的
		c.AppNo + ":" + c.AppChannel + ":",            // App编号:App渠道是自己的
	}

	o := orm.NewOrm()
	qs := o.QueryTable(new(Information))
	qs = qs.Filter("Receiver__in", receiver)
	qs = qs.Filter("EffectTime__lte", time.Now()).Filter("ExpireTime__gte", time.Now())
	qs = qs.OrderBy("-NeedFeedback", "-Id")
	qs = qs.Limit(5)
	var informations []Information
	if num, err := qs.All(&informations); err != nil {
		logs.Error("InformationApiController.query Information err:", err)
		c.Msg = "消息查询异常"
		return
	} else if num == 0 {
		c.Code = utils2.CODE_OK
		c.Msg = "没有消息"
		c.Dta = []map[string]interface{}{}
		return
	}
	infos := make([]map[string]interface{}, 0)
	informationIds := make([]int64, 0)
	for _, v := range informations {
		informationIds = append(informationIds, v.Id)
	}
	version := 1
	updateFeedBackIds := make([]int64, 0)          // 要更新的反馈记录
	newFeedBacks := make([]InformationFeedback, 0) // 新增的反馈记录
	if !c.IsLogon {
		infos = append(infos, map[string]interface{}{
			"infoId":       informations[0].Id,
			"title":        informations[0].Title,
			"info":         informations[0].Info,
			"needFeedback": informations[0].NeedFeedback,
		})
	} else {
		var feedbacks []InformationFeedback
		if num, err := o.QueryTable(new(InformationFeedback)).Filter("InformationId__in", informationIds).Filter("MemberId", c.LoginMemberId).All(&feedbacks); err != nil {
			logs.Error("InformationApiController.query InformationFeedback err:", err)
		} else if num == 0 { // 都没读
			for _, v := range informations {
				if v.NeedFeedback == 0 {
					version = 3 // 不需要反馈的直接标记为已读
				}
				newFeedBacks = append(newFeedBacks, InformationFeedback{InformationId: v.Id, MemberId: c.LoginMemberId, Version: version})
				infos = append(infos, map[string]interface{}{
					"infoId":       v.Id,
					"title":        v.Title,
					"info":         v.Info,
					"needFeedback": v.NeedFeedback,
				})
				version = 1
				break // 当前版本只取一个消息
			}
		} else { // 部分读
			var status = 0
			for _, v := range informations {
				for _, v2 := range feedbacks {
					if v.Id == v2.InformationId {
						if v2.Version >= 3 {
							status = 1 // 已读
						} else {
							status = 2 // 未读
							infos = append(infos, map[string]interface{}{
								"infoId":       v.Id,
								"title":        v.Title,
								"info":         v.Info,
								"needFeedback": v.NeedFeedback,
							})
							updateFeedBackIds = append(updateFeedBackIds, v2.Id)
						}
						break
					}
				}
				if status == 0 {
					if v.NeedFeedback == 0 {
						version = 3 // 不需要反馈的直接标记为已读
					}
					newFeedBacks = append(newFeedBacks, InformationFeedback{InformationId: v.Id, MemberId: c.LoginMemberId, Version: 1})
					infos = append(infos, map[string]interface{}{
						"infoId":       v.Id,
						"title":        v.Title,
						"info":         v.Info,
						"needFeedback": v.NeedFeedback,
					})
					version = 1
				}
				if status != 1 {
					break // 当前版本只取一个消息
				}
				status = 0
			}
		}
	}
	if len(newFeedBacks) > 0 {
		if _, err := o.InsertMulti(len(newFeedBacks), &newFeedBacks); err != nil {
			logs.Error("InformationApiController.InsertMulti InformationFeedback err:", err)
		}
	}
	if len(updateFeedBackIds) > 0 {
		if _, err := o.QueryTable(new(InformationFeedback)).Filter("Id__in", updateFeedBackIds).Update(orm.Params{
			"Version": orm.ColValue(orm.ColAdd, 1),
		}); err != nil {
			logs.Error("InformationApiController.Update InformationFeedback err:", err)
		}
	}
	c.Code = utils2.CODE_OK
	c.Msg = "查询成功"
	c.Dta = infos
}

type informationParam struct {
	InfoId int64 `json:"infoId"` // 必填
}

/*
api 系统消息更新已读状态
param:
body: {"infoId":1}
return:
{"code":1,"msg":"更新成功","data":"ok"}
desc:
*/
func (c *InformationApiController) Post() {
	defer c.RetJSON()
	if !c.IsLogon {
		c.Code = utils2.CODE_OK
		c.Msg = "更新成功"
		c.Dta = "ok"
	}
	var p informationParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils2.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	o := orm.NewOrm()
	if _, err := o.QueryTable(new(InformationFeedback)).Filter("InformationId", p.InfoId).Filter("MemberId", c.LoginMemberId).Update(orm.Params{
		"Version": 9,
	}); err != nil {
		logs.Error("InformationApiController.Post Update InformationFeedback err:", err)
	}
	if _, err := o.QueryTable(new(Information)).Filter("Id", p.InfoId).Update(orm.Params{
		"ReadNum": orm.ColValue(orm.ColAdd, 1),
	}); err != nil {
		logs.Error("InformationApiController.Post Update Information err:", err)
	}
	c.Code = utils2.CODE_OK
	c.Msg = "更新成功"
	c.Dta = "ok"
}
