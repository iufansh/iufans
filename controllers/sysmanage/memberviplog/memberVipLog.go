package memberviplog

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/iufansh/iufans/controllers/sysmanage"
	. "github.com/iufansh/iufans/models"
	"html/template"
)

type MemberVipLogIndexController struct {
	sysmanage.BaseController
}

func (c *MemberVipLogIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *MemberVipLogIndexController) Get() {
	memberId, _ := c.GetInt64("memberId", -1)

	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(MemberVipLog).Paginate(page, limit, memberId)
	c.SetPaginator(limit, total)
	// 返回值
	c.Data["dataList"] = &list
	// 查询条件
	c.Data["condArr"] = map[string]interface{}{"memberId": memberId}

	c.Data["urlMemberVipLogIndexGet"] = c.URLFor("MemberVipLogIndexController.Get")

	if t, err := template.New("tplIndexMemberVipLog.tpl").Funcs(map[string]interface{}{
		"date": beego.Date,
	}).Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}
