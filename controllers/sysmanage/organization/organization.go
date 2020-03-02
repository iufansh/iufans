package organization

import (
	"github.com/iufansh/iufans/controllers/sysmanage"
	"github.com/astaxie/beego"
	. "github.com/iufansh/iufans/models"
	"html/template"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"time"
	"strings"
	utils "github.com/iufansh/iutils"
	"strconv"
	"fmt"
)

func validate(org *Organization) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(org.Name, "errmsg").Message("组织名称必填")
	valid.MaxSize(org.BindDomain, 127, "errmsg").Message("绑定域名最长127个字符")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type OrganizationIndexController struct {
	sysmanage.BaseController
}

func (c *OrganizationIndexController) NestPrepare()  {
	c.EnableRender = false
}

func (c *OrganizationIndexController) Get() {
	param1 := strings.TrimSpace(c.GetString("param1"))
	orgId, _ := c.GetInt64("orgId", -1)
	if c.LoginAdminOrgId != 0 {
		orgId = c.LoginAdminOrgId
	}
	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(Organization).Paginate(page, limit, orgId, param1, []int{})
	c.SetPaginator(limit, total)
	c.Data["dataList"] = &list
	// 查询条件
	var orgIdStr string
	if orgId >= 0 {
		orgIdStr = strconv.FormatInt(orgId, 10)
	}
	c.Data["condArr"] = map[string]interface{}{"param1": param1, "orgId": orgIdStr}

	c.Data["urlOrgIndexGet"] = c.URLFor("OrganizationIndexController.Get")
	c.Data["urlOrgAddGet"] = c.URLFor("OrganizationAddController.Get")
	c.Data["urlOrgEditGet"] = c.URLFor("OrganizationEditController.Get")
	c.Data["urlOrgDelone"] = c.URLFor("OrganizationIndexController.Delone")
	c.Data["urlAdminAddGet"] = c.URLFor("AdminAddController.Get")
	c.Data["urlAdminIndexGet"] = c.URLFor("AdminIndexController.Get")
	c.Data["urlIpListIndexGet"] = c.URLFor("IpListIndexController.Get")

	if t, err := template.New("tplIndexOrg.tpl").Funcs(map[string]interface{}{
		"date": beego.Date,
		"formatAmount": func(a int64) string {
			return strconv.FormatFloat(float64(a)/100, 'f', 2, 64)
		},
	}).Parse(tplOrgIndex); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *OrganizationIndexController) Delone() {
	var code int
	var msg string
	url := beego.URLFor("OrganizationIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	if c.LoginAdminOrgId != 0 {
		msg = "没有权限"
		return
	}
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	_, err1 := o.Delete(&Organization{Id: id}, "Id")
	if err1 != nil {
		beego.Error("Delete Organization error", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type OrganizationAddController struct {
	sysmanage.BaseController
}

func (c *OrganizationAddController) NestPrepare()  {
	c.EnableRender = false
}

func (c *OrganizationAddController) Get() {
	c.Data["urlOrgIndexGet"] = c.URLFor("OrganizationIndexController.Get")
	c.Data["urlOrgAddPost"] = c.URLFor("OrganizationAddController.Post")

	if t, err := template.New("tplAddOrg.tpl").Parse(tplOrgAdd); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *OrganizationAddController) Post() {
	var code int
	var msg string
	var url = beego.URLFor("OrganizationIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	model := Organization{}
	if err := c.ParseForm(&model); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&model); hasError {
		msg = errMsg
		return
	}
	var zeroTime time.Time
	if model.ExpireTime == zeroTime {
		model.ExpireTime = time.Now().Add(time.Hour * 26280) // 默认三年
	}
	model.Enabled = 1
	model.EncryptKey = utils.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))
	model.Creator = c.LoginAdminId
	model.Modifior = c.LoginAdminId
	if model.BindDomain != "" && !strings.HasSuffix(model.BindDomain, ",") {
		model.BindDomain = model.BindDomain + ","
	}
	model.OrgId = c.LoginAdminOrgId
	o := orm.NewOrm()

	// 查询层级
	var levels string
	if c.LoginAdminOrgId != 0 {
		var org = Organization{Id: c.LoginAdminOrgId}
		if err := o.Read(&org); err != nil {
			msg = "添加失败"
			beego.Error("添加失败o.Read(&org)", err)
			return
		}
		levels = org.Levels
	} else {
		levels = "0,"
	}
	model.Levels = fmt.Sprintf("%s%d,", levels, c.LoginAdminOrgId)

	if _, err := o.Insert(&model); err != nil {
		msg = "添加失败"
		beego.Error("添加失败", err)
		return
	}
	code = 1
	msg = "添加成功"

}

type OrganizationEditController struct {
	sysmanage.BaseController
}

func (c *OrganizationEditController) NestPrepare()  {
	c.EnableRender = false
}

func (c *OrganizationEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := Organization{}
	model.Id = id
	err := o.Read(&model)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(beego.URLFor("OrganizationIndexController.get"), 302)
		return
	}
	c.Data["data"] = &model

	c.Data["urlOrgIndexGet"] = c.URLFor("OrganizationIndexController.Get")
	c.Data["urlOrgEditPost"] = c.URLFor("OrganizationEditController.Post")

	if t, err := template.New("tplEditOrg.tpl").Funcs(map[string]interface{}{
		"date": beego.Date,
	}).Parse(tplOrgEdit); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *OrganizationEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("OrganizationIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	model := Organization{}
	if err := c.ParseForm(&model); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&model); hasError {
		msg = errMsg
		return
	}
	if model.BindDomain != "" && !strings.HasSuffix(model.BindDomain, ",") {
		model.BindDomain = model.BindDomain + ","
	}
	cols := []string{"Name", "Vip", "BindDomain", "ModifyDate", "Modifior", "ExpireTime"}
	model.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Update(&model, cols...); err != nil {
		msg = "更新失败"
		beego.Error("更新失败", err)
	} else {
		code = 1
		msg = "更新成功"
	}
}
