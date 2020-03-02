package iplist

import (
	"github.com/iufansh/iufans/controllers/sysmanage"
	"github.com/astaxie/beego"
	. "github.com/iufansh/iufans/models"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
	"fmt"
	"html/template"
)

type IpListIndexController struct {
	sysmanage.BaseController
}

func (c *IpListIndexController) NestPrepare()  {
	c.EnableRender = false
}

func (c *IpListIndexController) Get() {
	param1 := strings.TrimSpace(c.GetString("param1"))
	black, _ := c.GetInt8("black", -1)
	orgId, _ := c.GetInt64("orgId", 0)
	if c.LoginAdminOrgId != 0 {
		orgId = c.LoginAdminOrgId
	}
	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(IpList).Paginate(page, limit, orgId, param1, black)
	c.SetPaginator(limit, total)
	c.Data["condArr"] = map[string]interface{}{"param1": param1, "orgId": orgId, "black": black}
	c.Data["dataList"] = list

	c.Data["urlIpListIndexGet"] = c.URLFor("IpListIndexController.Get")
	c.Data["urlIpListIndexDelone"] = c.URLFor("IpListIndexController.Delone")
	c.Data["urlIpListAddGet"] = c.URLFor("IpListAddController.Get")

	if t, err := template.New("tplIndexIpList.tpl").Funcs(map[string]interface{}{
		"date": beego.Date,
	}).Parse(tplIndex); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *IpListIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := IpList{}
	model.Id = id
	err := o.Read(&model)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	}
	if c.LoginAdminOrgId != 0 && model.OrgId != c.LoginAdminOrgId {
		msg = "非法请求"
		return
	}
	if _, err1 := o.Delete(&model, "Id"); err1 != nil {
		beego.Error("Delete IpList eror", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type IpListAddController struct {
	sysmanage.BaseController
}

func (c *IpListAddController) NestPrepare()  {
	c.EnableRender = false
}

func (c *IpListAddController) Get() {
	c.Data["orgId"], _ = c.GetInt64("orgId", 0)
	c.Data["curIp"] = c.Ctx.Input.IP()

	c.Data["urlIpListIndexGet"] = c.URLFor("IpListIndexController.Get")
	c.Data["urlIpListAddPost"] = c.URLFor("IpListAddController.Post")

	if t, err := template.New("tplAddIpList.tpl").Parse(tplAdd); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *IpListAddController) Post() {
	var code int
	var msg string
	var orgId int64
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	model := IpList{}
	if err := c.ParseForm(&model); err != nil {
		msg = "参数异常"
		return
	}
	ips := strings.Split(model.Ip, "\r\n")
	if len(ips) == 0 {
		msg = "请配置IP"
		return
	}
	orgId = model.OrgId
	if c.LoginAdminOrgId != 0 {
		orgId = c.LoginAdminOrgId
	}
	lists := make([]IpList, 0)
	for _, v := range ips {
		if strings.TrimSpace(v) == "" {
			continue
		}
		ipList := IpList{
			Creator:    c.LoginAdminId,
			Modifior:   c.LoginAdminId,
			CreateDate: time.Now(),
			ModifyDate: time.Now(),
			OrgId:      orgId,
			Ip:         strings.TrimSpace(v),
			Black:      model.Black,
		}
		lists = append(lists, ipList)
	}

	o := orm.NewOrm()
	if num, err := o.InsertMulti(len(lists), lists); err != nil {
		msg = "添加失败"
		beego.Error("添加失败", err)
	} else {
		code = 1
		msg = fmt.Sprintf("成功添加%d条IP", num)
	}
}
