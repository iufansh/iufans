package permission

import (
	"html/template"
	"github.com/iufansh/iufans/controllers/sysmanage"
	. "github.com/iufansh/iufans/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

func validate(permission *Permission) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(permission.Name, "errmsg").Message("菜单名必输")
	valid.MaxSize(permission.Name, 30, "errmsg").Message("菜单名最长30位")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type PermissionIndexController struct {
	sysmanage.BaseController
}

func (c *PermissionIndexController) NestPrepare()  {
	c.EnableRender = false
}

func (c *PermissionIndexController) Get() {
	var permissionList []Permission
	o := orm.NewOrm()
	qs := o.QueryTable(new(Permission))
	qs.All(&permissionList)
	// 返回值
	c.Data["dataList"] = &permissionList

	c.Data["urlIpListIndexGet"] = c.URLFor("IpListIndexController.Get")
	c.Data["urlPermissionIndexDelone"] = c.URLFor("PermissionIndexController.Delone")
	c.Data["urlPermissionAddGet"] = c.URLFor("PermissionAddController.Get")
	c.Data["urlPermissionEditGet"] = c.URLFor("PermissionEditController.Get")

	if t, err := template.New("tplPermissionIndex.tpl").Parse(tplIndex); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *PermissionIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, _ := c.GetInt64("id")
	permission := Permission{Id: id}
	o := orm.NewOrm()
	err := o.Read(&permission)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	}
	_, err1 := o.Delete(&permission)
	if err1 != nil {
		beego.Error("Delete permission error", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type PermissionAddController struct {
	sysmanage.BaseController
}

func (c *PermissionAddController) NestPrepare()  {
	c.EnableRender = false
}

func (c *PermissionAddController) Get() {
	c.Data["urlPermissionIndexGet"] = c.URLFor("PermissionIndexController.Get")
	c.Data["urlPermissionAddPost"] = c.URLFor("PermissionAddController.Post")

	if t, err := template.New("tplAddPermission.tpl").Parse(tplAdd); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *PermissionAddController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	permission := Permission{}
	if err := c.ParseForm(&permission); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&permission); hasError {
		msg = errMsg
		return
	}
	permission.Creator = c.LoginAdminId
	permission.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Insert(&permission); err != nil {
		msg = "添加失败"
		beego.Error("Insert permission error", err)
	} else {
		code = 1
		msg = "添加成功"
	}
}

type PermissionEditController struct {
	sysmanage.BaseController
}

func (c *PermissionEditController) NestPrepare()  {
	c.EnableRender = false
}

func (c *PermissionEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	permission := Permission{Id: id}

	err := o.Read(&permission)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(beego.URLFor("PermissionIndexController.get"), 302)
	} else {
		c.Data["data"] = &permission

		c.Data["urlPermissionIndexGet"] = c.URLFor("PermissionIndexController.Get")
		c.Data["urlPermissionEditPost"] = c.URLFor("PermissionEditController.Post")

		if t, err := template.New("tplEditPermission.tpl").Parse(tplEdit); err != nil {
			beego.Error("template Parse err", err)
		} else {
			t.Execute(c.Ctx.ResponseWriter, c.Data)
		}
	}
}

func (c *PermissionEditController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	permission := Permission{}
	if err := c.ParseForm(&permission); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&permission); hasError {
		msg = errMsg
		return
	}
	permission.Modifior = c.LoginAdminId
	cols := []string{"Name", "Description", "Enabled", "Pid", "Url", "Icon", "Sort", "Display", "ModifyDate"}
	o := orm.NewOrm()
	if num, err := o.Update(&permission, cols...); err != nil {
		msg = "更新失败"
		beego.Error("Update permission error", err)
	} else if num == 0 {
		msg = "更新失败"
	} else {
		code = 1
		msg = "更新成功"
	}
}
