package sysapi

import (
	"github.com/astaxie/beego/orm"
	. "github.com/iufansh/iufans/models"
	utils2 "github.com/iufansh/iufans/utils"
)

type AppBannerApiController struct {
	Base2ApiController
}

/*
api app轮播图获取
param:
body:
return:{"code":1,"msg":"成功","data":[{"title":"轮播1","banner":"图片地址","jumpUrl":"http://baidu.com"}]}
desc:
*/
func (c *AppBannerApiController) Get() {
	defer c.RetJSON()
	o := orm.NewOrm()
	qs := o.QueryTable(new(AppBanner))
	qs = qs.Filter("AppNo", c.AppNo)
	qs = qs.Filter("Status", 1)
	qs = qs.OrderBy("Seq", "-Id")
	var appBanners []AppBanner
	if _, err := qs.All(&appBanners); err != nil {
		c.Msg = "获取失败"
		return
	}
	retList := make([]map[string]string, 0)
	for _, v := range appBanners {
		retList = append(retList, map[string]string{
			"title":   v.Title,
			"banner":  v.GetFullBanner(c.Ctx.Input.Site()),
			"jumpUrl": v.JumpUrl,
		})
	}
	c.Dta = &retList
	c.Code = utils2.CODE_OK
	c.Msg = "ok"
}
