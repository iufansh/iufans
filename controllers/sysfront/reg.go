package sysfront

import (
	"github.com/astaxie/beego/logs"
	fm "github.com/iufansh/iufans/models"
	. "github.com/iufansh/iutils"

	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	. "github.com/iufansh/iufans/models"
	"github.com/iufansh/iufans/utils"
	"strings"
)

type RegFrontController struct {
	Base2FrontController
}

func (c *RegFrontController) Get() {
	inviteCode := c.GetString("ic")
	if inviteCode == "" {
		icv := c.GetSession("invite_code")
		if icv != nil {
			inviteCode = icv.(string)
		}
	}
	m := fm.GetSiteConfigMap(utils.Scname, utils.Scfrontregsmsverify)
	c.Data["inviteCode"] = strings.TrimSpace(inviteCode)
	c.Data["siteName"] = m[utils.Scname]
	c.Data["smsVerify"] = m[utils.Scfrontregsmsverify]
	c.TplName = "front/reg.html"
}

func (c *RegFrontController) Post() {
	defer c.RetJSON()
	if agree, _ := c.GetBool("protocol", false); !agree {
		c.Msg = "请阅读并同意用户协议"
		return
	}
	inviteCode := c.GetString("inviteCode", "0")

	model := Member{}
	if err := c.ParseForm(&model); err != nil {
		c.Msg = "异常，请刷新后重试"
		return
	}
	if !utils.GetCpt().VerifyReq(c.Ctx.Request) {
		c.Msg = "请刷新验证码，重新输入"
		return
	}
	scMap := fm.GetSiteConfigMap(utils.Scfrontregsmsverify)
	if scMap[utils.Scfrontregsmsverify] == "1" {
		smsCode := c.GetString("smsCode")
		if smsCode == "" {
			c.Msg = "短信验证码必须填写"
			return
		}
		if ok := utils.VerifySmsVerifyCode(model.Username, smsCode); !ok {
			c.Msg = "短信验证码错误"
			return
		}
	}
	valid := validation.Validation{}
	valid.Required(model.Username, "errmsg").Message("手机号必填")
	valid.MaxSize(model.Username, 11, "errmsg").Message("手机号最长11位")
	valid.Required(model.Password, "errmsg").Message("密码必填")
	valid.MinSize(model.Password, 6, "errmsg").Message("密码最少6位")
	valid.MaxSize(model.Password, 20, "errmsg").Message("密码最多20位")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.Msg = err.Message
			return
		}
	}
	if c.GetString("confirmPass") != model.Password {
		c.Msg = "两次输入的密码不一致"
		return
	}
	// 验证用户名是否存在
	o := orm.NewOrm()
	if isExist := o.QueryTable(new(Member)).Filter("Username", model.Username).Exist(); isExist {
		c.Msg = "手机号已存在"
		return
	}
	// 查询层级
	var refId int64
	if inviteCode != "999" { // 999 默认首层
		refId = utils.ReverseInviteCode(inviteCode)
	} else {
		refId = 0
	}
	// 查询层级
	if refId != 0 {
		var refMem Member
		if err := o.QueryTable(new(Member)).Filter("Id", refId).One(&refMem, "Levels", "LevelsDeep"); err != nil {
			logs.Error("member reg QueryTable Member err", err)
			model.Levels = "0,"
			model.LevelsDeep = 1
			model.RefId = 0
		} else {
			model.Levels = fmt.Sprintf("%s%d,", refMem.Levels, refId)
			model.LevelsDeep = refMem.LevelsDeep + 1
			model.RefId = refId
		}
	} else {
		model.Levels = "0,"
		model.LevelsDeep = 1
		model.RefId = 0
	}
	if model.Name == "" {
		if len(model.Username) == 11 && strings.HasPrefix(model.Username, "1") {
			model.Mobile = model.Username
			model.Name = SubString(model.Username, 0, 3) + "*****" + SubString(model.Username, 8, 3)
		} else {
			model.Name = model.Username
		}
	}
	// 注册送积分
	//integral, _ := strconv.ParseInt(scMap[utils.Scfrontregintegral], 10, 64)
	salt := GetGuid()
	pa := Md5(Md5(model.Password), Pubsalt, salt)
	model.Mobile = model.Username
	model.Password = pa
	model.Salt = salt
	model.Creator = 0
	model.Modifior = 0
	model.Enabled = 1
	model.Locked = 0
	model.LoginFailureCount = 0

	o.Begin()
	var memberId int64
	var err error
	if memberId, err = o.Insert(&model); err != nil {
		o.Rollback()
		logs.Error("memberRegErr Member error", err, memberId)
		c.Msg = "注册失败，请重试(3)"
		return
	}
	o.Commit()

	c.Msg = "注册成功"
	c.Code = 1
	c.Dta = c.URLFor("LoginFrontController.Get")

	// go GenerateRandAvatar(memberId)
	return
}

// 生成一个随机头像
func GenerateRandAvatar(id int64) {
	/*
	avatarRes := httplib.Get("https://api.uomg.com/api/rand.avatar?format=images")
	avatar, err := avatarRes.String() // avatar是一个base64的数据。
	fmt.Println(avatar)
	if err != nil || !strings.HasPrefix(strings.ToLower(avatar), "http") {
		avatar = "/static/front/images/avatar/0.jpg"
	}
	*/
	avatar := fmt.Sprintf("/static/front/images/avatar/%d.jpg", id%12)
	o := orm.NewOrm()
	o.Update(&Member{Id: id, Avatar: avatar}, "Avatar")
}
