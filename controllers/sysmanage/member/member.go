package member

import (
	"fmt"
	"html/template"
	"github.com/iufansh/iufans/controllers/sysmanage"
	. "github.com/iufansh/iufans/models"
	. "github.com/iufansh/iufans/utils"
	. "github.com/iufansh/iutils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"time"
	"strings"
)

func validate(member *Member) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(member.Username, "errmsg").Message("用户名必填")
	valid.AlphaDash(member.Username, "errmsg").Message("用户名必须为字母和数字")
	valid.MinSize(member.Username, 5, "errmsg").Message("用户名至少5个字符")
	valid.MaxSize(member.Username, 20, "errmsg").Message("用户名最长20位")
	valid.Required(member.Name, "errmsg").Message("名称必填")
	valid.MaxSize(member.Name, 20, "errmsg").Message("名称最长20位")
	valid.MaxSize(member.Password, 32, "errmsg").Message("密码不符合规范")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type MemberIndexController struct {
	sysmanage.BaseController
}

func (c *MemberIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *MemberIndexController) Get() {
	param1 := strings.TrimSpace(c.GetString("param1"))
	refId, _ := c.GetInt64("refId", -1)

	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(Member).Paginate(page, limit, refId, param1)
	c.SetPaginator(limit, total)
	// 返回值
	c.Data["dataList"] = &list
	// 查询条件
	c.Data["condArr"] = map[string]interface{}{"param1": param1, "refId": refId}

	c.Data["urlMemberIndexGet"] = c.URLFor("MemberIndexController.Get")
	c.Data["urlMemberLocked"] = c.URLFor("MemberIndexController.Locked")
	c.Data["urlMemberEditGet"] = c.URLFor("MemberEditController.Get")

	if t, err := template.New("tplIndexMember.tpl").Funcs(map[string]interface{}{
		"date": beego.Date,
	}).Parse(tplIndex); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *MemberIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, err := c.GetInt64("id")
	if err != nil {
		msg = "数据错误"
		beego.Error("Delete Member error", err)
		return
	}
	o := orm.NewOrm()
	if _, err := o.Delete(&Member{Id: id}); err != nil {
		beego.Error("Delete admin error 2", err)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

func (c *MemberIndexController) Locked() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, err := c.GetInt64("id")
	if err != nil {
		msg = "数据错误"
		beego.Error("Locked Member error", err)
		return
	}
	o := orm.NewOrm()
	model := Member{Id: id}
	if err := o.Read(&model); err != nil {
		beego.Error("Read admin error", err)
		msg = "操作失败，请刷新后重试"
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
		beego.Error("Update admin error", err)
		msg = "操作失败，请刷新后重试"
	} else {
		code = 1
		msg = "操作成功"
		if model.Locked == 1 { // 如果是锁定，则一并清楚登录token，强制用户退出
			DelCache(fmt.Sprintf("loginMemberId%d", id))
		}
	}
}

type MemberEditController struct {
	sysmanage.BaseController
}

func (c *MemberEditController) NestPrepare() {
	c.EnableRender = false
}

func (c *MemberEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	admin := Member{Id: id}

	err := o.Read(&admin)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(beego.URLFor("MemberIndexController.get"), 302)
		return
	}
	c.Data["data"] = &admin

	c.Data["urlMemberIndexGet"] = c.URLFor("MemberIndexController.Get")
	c.Data["urlMemberEditPost"] = c.URLFor("MemberEditController.Post")

	if t, err := template.New("tplEditMember.tpl").Parse(tplEdit); err != nil {
		beego.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *MemberEditController) Post() {
	var code int
	var msg string
	var reurl = c.URLFor("MemberIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &reurl)
	admin := Member{}
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
	o := orm.NewOrm()

	cols := []string{"Name", "Vip", "Enabled", "ModifyDate"}
	isChangePwd := false
	if admin.Password != "" {
		salt := GetGuid()
		admin.Password = Md5(admin.Password, Pubsalt, salt)
		admin.Salt = salt
		cols = append(cols, "Password", "Salt")
		isChangePwd = true
	}
	admin.Modifior = c.LoginAdminId
	if _, err := o.Update(&admin, cols...); err != nil {
		msg = "更新失败"
		beego.Error("Update admin error 1", err)
	} else {
		// 如修改了密码，则重置登录，让用户必须重新登录
		if isChangePwd {
			DelCache(fmt.Sprintf("loginMemberId%d", admin.Id))
		}
		code = 1
		msg = "更新成功"
	}
}
