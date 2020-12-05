package admin

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/iufansh/iufans/controllers/sysmanage"
	. "github.com/iufansh/iufans/models"
	. "github.com/iufansh/iufans/utils"
	. "github.com/iufansh/iutils"
	"html/template"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"strings"
	"time"
)

func validate(admin *Admin) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(admin.Username, "errmsg").Message("用户名必输")
	valid.AlphaDash(admin.Username, "errmsg").Message("用户名必须为字母和数字")
	valid.MinSize(admin.Username, 5, "errmsg").Message("用户名至少5个字符")
	valid.MaxSize(admin.Username, 20, "errmsg").Message("用户名最长20位")
	valid.Required(admin.Name, "errmsg").Message("名称必输")
	valid.MaxSize(admin.Name, 20, "errmsg").Message("名称最长20位")
	valid.MaxSize(admin.Password, 32, "errmsg").Message("密码不符合规范")
	valid.MaxSize(admin.Mobile, 11, "errmsg").Message("手机号最长11位")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type AdminIndexController struct {
	sysmanage.BaseController
}

func (c *AdminIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *AdminIndexController) Get() {
	param1 := strings.TrimSpace(c.GetString("param1"))
	orgId, _ := c.GetInt64("orgId", 0)
	if c.LoginAdminOrgId != 0 {
		orgId = c.LoginAdminOrgId
	}
	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(Admin).Paginate(page, limit, orgId, param1)
	c.SetPaginator(limit, total)
	// 返回值
	c.Data["dataList"] = &list
	// 查询条件
	c.Data["condArr"] = map[string]interface{}{"param1": param1, "orgId": orgId}

	c.Data["urlAdminAddGet"] = c.URLFor("AdminAddController.Get")
	c.Data["urlAdminIndexGet"] = c.URLFor("AdminIndexController.Get")
	c.Data["urlAdminLoginVerify"] = c.URLFor("AdminIndexController.LoginVerify")
	c.Data["urlAdminLocked"] = c.URLFor("AdminIndexController.Locked")
	c.Data["urlAdminEditGet"] = c.URLFor("AdminEditController.Get")
	c.Data["urlAdminDelone"] = c.URLFor("AdminIndexController.Delone")

	if t, err := template.New("tplIndexAdmin.tpl").Funcs(map[string]interface{}{
		"date": beego.Date,
	}).Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *AdminIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, err := c.GetInt64("id")
	if err != nil {
		msg = "数据错误"
		logs.Error("Delete Admin error", err)
		return
	}
	o := orm.NewOrm()
	// 验证数据权限
	if c.LoginAdminOrgId != 0 {
		msg = "无法删除，请使用锁定功能"
	}
	o.Begin()
	// 删除管理员角色关联
	if _, err := o.QueryTable(new(AdminRole)).Filter("AdminId", id).Delete(); err != nil {
		o.Rollback()
		logs.Error("Delete admin error 1", err)
		msg = "删除失败"
		return
	}
	if _, err := o.Delete(&Admin{Id: id}); err != nil {
		o.Rollback()
		logs.Error("Delete admin error 2", err)
		msg = "删除失败"
	} else {
		o.Commit()
		code = 1
		msg = "删除成功"
	}
}

func (c *AdminIndexController) LoginVerify() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, err := c.GetInt64("id")
	if err != nil {
		msg = "数据错误"
		logs.Error("LoginVerify Admin error", err)
		return
	}
	o := orm.NewOrm()
	model := Admin{Id: id}
	if err := o.Read(&model); err != nil {
		logs.Error("Read admin error", err)
		msg = "操作失败，请刷新后重试"
		return
	}
	if c.LoginAdminOrgId != 0 && c.LoginAdminOrgId != model.OrgId {
		msg = "非法操作"
		return
	}
	model.LoginVerify = 0
	model.GaSecret = ""

	if _, err := o.Update(&model, "LoginVerify", "GaSecret"); err != nil {
		logs.Error("Update admin error", err)
		msg = "解除失败，请刷新后重试"
	} else {
		code = 1
		msg = "解除成功"
	}
}

func (c *AdminIndexController) Locked() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, err := c.GetInt64("id")
	if err != nil {
		msg = "数据错误"
		logs.Error("Locked Admin error", err)
		return
	}
	o := orm.NewOrm()
	model := Admin{Id: id}
	if err := o.Read(&model); err != nil {
		logs.Error("Read admin error", err)
		msg = "操作失败，请刷新后重试"
		return
	}
	if c.LoginAdminOrgId != 0 && c.LoginAdminOrgId != model.OrgId {
		msg = "非法操作"
		return
	}
	if model.Locked == 1 {
		model.Locked = 0
		model.LoginFailureCount = 0
		model.LockedDate = time.Now()
	} else {
		model.Locked = 1
		model.LockedDate = time.Now()
	}

	if _, err := o.Update(&model, "Locked", "LockedDate", "LoginFailureCount"); err != nil {
		logs.Error("Update admin error", err)
		msg = "操作失败，请刷新后重试"
	} else {
		code = 1
		msg = "操作成功"
		if model.Locked == 1 { // 如果是锁定，则一并清楚登录token，强制用户退出
			DelCache(fmt.Sprintf("loginAdminId%d", id))
		}
	}
}

type AdminAddController struct {
	sysmanage.BaseController
}

func (c *AdminAddController) NestPrepare() {
	c.EnableRender = false
}

func (c *AdminAddController) Get() {
	orgId, _ := c.GetInt64("orgId", 0)
	if orgId == 0 {
		orgId = c.LoginAdminOrgId
	}
	c.Data["orgId"] = orgId
	c.Data["isOrg"] = c.LoginAdminOrgId != 0
	c.Data["prefix"] = ""
	c.Data["roleList"] = GetRoleList(c.LoginAdminOrgId != 0)

	c.Data["urlAdminIndexGet"] = c.URLFor("AdminIndexController.Get")
	c.Data["urlAdminAddPost"] = c.URLFor("AdminAddController.Post")

	if t, err := template.New("tplAddAdmin.tpl").Parse(tplAdd); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *AdminAddController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	admin := Admin{}
	if err := c.ParseForm(&admin); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&admin); hasError {
		msg = errMsg
		return
	} else if admin.Password == "" {
		msg = "密码不能为空"
		return
	} else if admin.Password != c.GetString("repassword") {
		msg = "两次输入的密码不一致"
		return
	}
	roles := c.GetStrings("roles")
	if len(roles) == 0 {
		msg = "请选择所属权限组"
		return
	}
	o := orm.NewOrm()
	if c.LoginAdminOrgId != 0 {
		if admin.OrgId == 0 {
			admin.OrgId = c.LoginAdminOrgId
		} else if admin.OrgId != c.LoginAdminOrgId {
			var org Organization
			o.QueryTable(new(Organization)).Filter("Id", c.LoginAdminOrgId).One(&org)
			levels := fmt.Sprintf("%s%d,", org.Levels, c.LoginAdminOrgId)
			if exists := o.QueryTable(new(Organization)).Filter("Levels", levels).Filter("Id", admin.OrgId).Exist(); !exists {
				msg = "组织获取异常，请刷新后重试"
				return
			}
		}

		if count, err := o.QueryTable(new(Role)).Filter("IsOrg", 1).Filter("Id__in", roles).Count(); err != nil || int(count) != len(roles) {
			msg = "权限获取异常，请刷新后重试"
			return
		}
	}
	if admin.OrgId != 0 {
		// 添加用户名前缀
		// admin.Username = Num2Letters(admin.OrgId) + "_" + admin.Username
	}
	salt := GetGuid()
	admin.Password = Md5(admin.Password, Pubsalt, salt)
	admin.Salt = salt
	admin.Creator = c.LoginAdminId
	admin.Modifior = c.LoginAdminId
	o.Begin()
	if created, _, err := o.ReadOrCreate(&admin, "Username"); err != nil {
		o.Rollback()
		msg = "添加失败"
		logs.Error("Insert admin error 1", err)
	} else if created {
		adminRoles := make([]AdminRole, 0)
		for _, v := range roles {
			roleId, _ := strconv.ParseInt(v, 10, 64)
			ar := AdminRole{AdminId: admin.Id, RoleId: roleId}
			adminRoles = append(adminRoles, ar)
		}
		if _, err := o.InsertMulti(len(adminRoles), adminRoles); err != nil {
			o.Rollback()
			msg = "添加失败"
			logs.Error("Insert admin error 3", err)
			return
		}
		o.Commit()
		code = 1
		msg = "添加成功，账号：" + admin.Username
		if admin.OrgId != c.LoginAdminOrgId {
			msg = msg + "；下级组织账号仅显示一次，请复制账号，以免忘记！"
		}
	} else {
		msg = "账号已存在或不可用，请更换"
	}
}

type AdminEditController struct {
	sysmanage.BaseController
}

func (c *AdminEditController) NestPrepare() {
	c.EnableRender = false
}

func (c *AdminEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	admin := Admin{Id: id}

	err := o.Read(&admin)

	if c.LoginAdminOrgId != 0 && c.LoginAdminOrgId != admin.OrgId {
		c.Redirect(beego.URLFor("AdminIndexController.get"), 302)
		return
	}
	arMap := make(map[int64]bool)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(beego.URLFor("AdminIndexController.get"), 302)
		return
	} else if c.LoginAdminOrgId == 0 || c.LoginAdminOrgId == admin.OrgId {
		// 当前管理员所属角色
		var arList orm.ParamsList
		o.QueryTable(new(AdminRole)).Filter("AdminId", id).ValuesFlat(&arList, "RoleId")
		for _, v := range arList {
			arId, ok := v.(int64)
			if ok {
				arMap[arId] = true
			}
		}
	}
	c.Data["data"] = &admin
	c.Data["adminRoleMap"] = arMap
	c.Data["roleList"] = GetRoleList(c.LoginAdminOrgId != 0)

	c.Data["urlAdminIndexGet"] = c.URLFor("AdminIndexController.Get")
	c.Data["urlAdminEditPost"] = c.URLFor("AdminEditController.Post")

	if t, err := template.New("tplEditAdmin.tpl").Parse(tplEdit); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *AdminEditController) Post() {
	var code int
	var msg string
	var reurl = c.URLFor("AdminIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &reurl)
	admin := Admin{}
	if err := c.ParseForm(&admin); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&admin); hasError {
		msg = errMsg
		return
	} else if admin.Password != "" && admin.Password != c.GetString("repassword") {
		msg = "两次输入的密码不一致"
		return
	}
	roles := c.GetStrings("roles")
	o := orm.NewOrm()
	// 验证数据权限
	if c.LoginAdminOrgId != 0 {
		if exists := o.QueryTable(new(Admin)).Filter("Id", admin.Id).Filter("OrgId", c.LoginAdminOrgId).Exist(); !exists {
			msg = "非法操作"
			return
		}
		if len(roles) > 0 {
			if count, err := o.QueryTable(new(Role)).Filter("IsOrg", 1).Filter("Id__in", roles).Count(); err != nil || int(count) != len(roles) {
				msg = "权限获取异常，请刷新后重试"
				return
			}
		}
	}
	cols := []string{"Name", "Enabled", "Mobile", "ModifyDate"}
	isChangePwd := false
	if admin.Password != "" {
		salt := GetGuid()
		admin.Password = Md5(admin.Password, Pubsalt, salt)
		admin.Salt = salt
		cols = append(cols, "Password", "Salt")
		isChangePwd = true
	}
	admin.Modifior = c.LoginAdminId
	o.Begin()
	if _, err := o.Update(&admin, cols...); err != nil {
		o.Rollback()
		msg = "更新失败"
		logs.Error("Update admin error 1", err)
	} else {
		// 删除旧角色
		if _, err := o.QueryTable(new(AdminRole)).Filter("AdminId", admin.Id).Delete(); err != nil {
			o.Rollback()
			msg = "更新失败"
			logs.Error("Update admin error 2", err)
		}
		// 重新插入角色
		if len(roles) > 0 {
			adminRoles := make([]AdminRole, 0)
			for _, v := range roles {
				roleId, _ := strconv.ParseInt(v, 10, 64)
				ar := AdminRole{AdminId: admin.Id, RoleId: roleId}
				adminRoles = append(adminRoles, ar)
			}

			if _, err := o.InsertMulti(len(adminRoles), adminRoles); err != nil {
				o.Rollback()
				msg = "更新失败"
				logs.Error("Update admin error 3", err)
			}
		}
		o.Commit()
		// 如修改了密码，则重置登录，让用户必须重新登录
		if isChangePwd {
			DelCache(fmt.Sprintf("loginAdminId%d", admin.Id))
		}

		code = 1
		msg = "更新成功"
	}
}
