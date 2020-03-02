package role

import (
	"html/template"
	"github.com/iufansh/iufans/controllers/sysmanage"
	. "github.com/iufansh/iufans/models"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

func validate(role *Role) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(role.Name, "errmsg").Message("角色名必输")
	valid.MaxSize(role.Name, 50, "errmsg").Message("角色名最长50位")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type RoleIndexController struct {
	sysmanage.BaseController
}

func (c *RoleIndexController) NestPrepare()  {
	c.EnableRender = false
}

func (c *RoleIndexController) Get() {
	var roleList []Role
	o := orm.NewOrm()
	qs := o.QueryTable(new(Role))
	qs.All(&roleList)
	// 返回值
	c.Data["dataList"] = &roleList

	c.Data["urlRoleIndexDelone"] = c.URLFor("RoleIndexController.Delone")
	c.Data["urlRoleAddGet"] = c.URLFor("RoleAddController.Get")
	c.Data["urlRoleEditGet"] = c.URLFor("RoleEditController.Get")

	if t, err := template.New("tplPermissionIndex.tpl").Funcs(map[string]interface{}{
		"date": beego.Date,
	}).Parse(tplIndex); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *RoleIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, _ := c.GetInt64("id")
	role := Role{Id: id}
	o := orm.NewOrm()
	err := o.Read(&role)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	} else if role.IsSystem == 1 {
		msg = "系统内置角色，不能删除"
		return
	}
	// 先删除角色权限关联
	o.Begin()
	if _, err := o.QueryTable(new(RolePermission)).Filter("RoleId", role.Id).Delete(); err != nil {
		o.Rollback()
		beego.Error("Delete role error 1", err)
		msg = "删除失败"
		return
	}

	if _, err := o.Delete(&Role{Id: id}); err != nil {
		o.Rollback()
		beego.Error("Delete role error 2", err)
		msg = "删除失败"
		return
	}
	o.Commit()
	code = 1
	msg = "删除成功"
}

type RoleAddController struct {
	sysmanage.BaseController
}

func (c *RoleAddController) NestPrepare()  {
	c.EnableRender = false
}

func (c *RoleAddController) Get() {
	c.Data["permissionList"] = GetPermissionList()

	c.Data["urlRoleIndexGet"] = c.URLFor("RoleIndexController.Get")
	c.Data["urlRoleAddPost"] = c.URLFor("RoleAddController.Post")

	if t, err := template.New("tplAddRole.tpl").Parse(tplAdd); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *RoleAddController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	role := Role{}
	if err := c.ParseForm(&role); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&role); hasError {
		msg = errMsg
		return
	}
	role.Creator = c.LoginAdminId
	role.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Insert(&role); err != nil {
		o.Rollback()
		msg = "添加失败"
		beego.Error("Insert role error 1", err)
	} else {
		permissions := c.GetStrings("permissions")
		rolePermissions := make([]RolePermission, 0)
		for _, v := range permissions {
			permissionId, _ := strconv.ParseInt(v, 10, 64)
			ar := RolePermission{RoleId: role.Id, PermissionId: permissionId}
			rolePermissions = append(rolePermissions, ar)
		}
		if _, err := o.InsertMulti(len(rolePermissions), rolePermissions); err != nil {
			o.Rollback()
			msg = "添加失败"
			beego.Error("Insert role error 2", err)
			return
		}
		o.Commit()
		code = 1
		msg = "添加成功"
	}
}

type RoleEditController struct {
	sysmanage.BaseController
}

func (c *RoleEditController) NestPrepare()  {
	c.EnableRender = false
}

func (c *RoleEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	role := Role{Id: id}

	err := o.Read(&role)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(beego.URLFor("RoleIndexController.get"), 302)
	} else {
		// 当前角色包含的权限
		var rpList orm.ParamsList
		o.QueryTable(new(RolePermission)).Filter("RoleId", id).ValuesFlat(&rpList, "PermissionId")
		rpMap := make(map[int64]bool)
		for _, v := range rpList {
			rpId, ok := v.(int64)
			if ok {
				rpMap[rpId] = true
			}
		}
		c.Data["data"] = &role
		c.Data["rolePermissionMap"] = rpMap
		c.Data["permissionList"] = GetPermissionList()

		c.Data["urlRoleIndexGet"] = c.URLFor("RoleIndexController.Get")
		c.Data["urlRoleEditPost"] = c.URLFor("RoleEditController.Post")

		if t, err := template.New("tplEditRole.tpl").Parse(tplEdit); err != nil {
			beego.Error("template Parse err", err)
		} else {
			t.Execute(c.Ctx.ResponseWriter, c.Data)
		}
	}
}

func (c *RoleEditController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	role := Role{}
	if err := c.ParseForm(&role); err != nil {
		msg = "参数异常"
		return
	}
	cols := []string{"Name", "Description", "Enabled", "ModifyDate"}
	role.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	o.Begin()
	if _, err := o.Update(&role, cols...); err != nil {
		o.Rollback()
		msg = "更新失败"
		beego.Error("Update role error 1", err)
	} else {
		// 删除旧权限
		if _, err := o.QueryTable(new(RolePermission)).Filter("RoleId", role.Id).Delete(); err != nil {
			o.Rollback()
			msg = "更新失败"
			beego.Error("Update role error 2", err)
			return
		}
		// 重新插入新权限
		permissions := c.GetStrings("permissions")
		rolePermissions := make([]RolePermission, 0)
		for _, v := range permissions {
			permissionId, _ := strconv.ParseInt(v, 10, 64)
			ar := RolePermission{RoleId: role.Id, PermissionId: permissionId}
			rolePermissions = append(rolePermissions, ar)
		}
		if _, err := o.InsertMulti(len(rolePermissions), rolePermissions); err != nil {
			o.Rollback()
			msg = "更新失败"
			beego.Error("Update role error 3", err)
			return
		}
		o.Commit()
		code = 1
		msg = "更新成功"
	}
}
