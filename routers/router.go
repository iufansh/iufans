package routers

import (
	"github.com/iufansh/iufans/controllers"
	"github.com/iufansh/iufans/controllers/sysmanage"
	"github.com/iufansh/iufans/controllers/sysmanage/admin"
	"github.com/iufansh/iufans/controllers/sysmanage/gift"
	"github.com/iufansh/iufans/controllers/sysmanage/index"
	"github.com/iufansh/iufans/controllers/sysmanage/information"
	"github.com/iufansh/iufans/controllers/sysmanage/normalquestion"
	"github.com/iufansh/iufans/controllers/sysmanage/iplist"
	"github.com/iufansh/iufans/controllers/sysmanage/login"
	"github.com/iufansh/iufans/controllers/sysmanage/organization"
	"github.com/iufansh/iufans/controllers/sysmanage/permission"
	"github.com/iufansh/iufans/controllers/sysmanage/quicknav"
	"github.com/iufansh/iufans/controllers/sysmanage/role"
	"github.com/iufansh/iufans/controllers/sysmanage/siteconfig"

	"github.com/astaxie/beego"
	"github.com/iufansh/iufans/controllers/sysapi"
	"github.com/iufansh/iufans/controllers/sysfront"
	"github.com/iufansh/iufans/controllers/sysmanage/appversion"
	"github.com/iufansh/iufans/controllers/sysmanage/member"
	"github.com/iufansh/iufans/controllers/sysmanage/membersuggest"
	"github.com/iufansh/iufans/controllers/sysmanage/paymentconfig"
)

func init() {
	// 禁止使用的前缀 i
	beego.Router("/i/:appNo", &sysfront.AppDownloadFrontController{}, "get:DownloadRedirect")
	beego.Router("/i/privacy", &sysfront.PrivacyFrontController{})
	beego.Router("/i/protocol", &sysfront.ProtocolFrontController{})

	beego.Router("/sendsmscode", &controllers.CommonController{}, "get:SendSmsCode")
	beego.Router("/healthcheck", &controllers.CommonController{}, "get:HealthCheck")
	beego.Router("/serversysteminfo", &controllers.CommonController{}, "get:SystemInfo")
	var adminRouter = beego.AppConfig.String("adminrouter")
	beego.ErrorController(&controllers.ErrorController{})
	beego.Router(adminRouter+"/sys/main", &sysmanage.BaseIndexController{}, "get:Get")
	beego.Router(adminRouter+"/sys/index", &index.SysIndexController{})
	beego.Router(adminRouter+"/sys/getauth", &index.SysIndexController{}, "get:GetAuth")
	beego.Router(adminRouter+"/sys/postauth", &index.SysIndexController{}, "post:PostAuth")

	beego.Router(adminRouter+"/sys/upload", &sysmanage.SyscommonController{}, "post:Upload")
	beego.Router(adminRouter+"/sys/uploadmulti", &sysmanage.SyscommonController{}, "post:UploadMulti")

	beego.Router(adminRouter+"/org/index", &organization.OrganizationIndexController{})
	beego.Router(adminRouter+"/org/delone", &organization.OrganizationIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/org/add", &organization.OrganizationAddController{})
	beego.Router(adminRouter+"/org/edit", &organization.OrganizationEditController{})

	beego.Router(adminRouter+"/admin/index", &admin.AdminIndexController{})
	beego.Router(adminRouter+"/admin/delone", &admin.AdminIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/admin/locked", &admin.AdminIndexController{}, "post:Locked")
	beego.Router(adminRouter+"/admin/LoginVerify", &admin.AdminIndexController{}, "post:LoginVerify")
	beego.Router(adminRouter+"/admin/add", &admin.AdminAddController{})
	beego.Router(adminRouter+"/admin/edit", &admin.AdminEditController{})
	beego.Router(adminRouter+"/changepwd/index", &admin.ChangePwdController{})

	beego.Router(adminRouter+"/role/index", &role.RoleIndexController{})
	beego.Router(adminRouter+"/role/delone", &role.RoleIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/role/add", &role.RoleAddController{})
	beego.Router(adminRouter+"/role/edit", &role.RoleEditController{})

	beego.Router(adminRouter+"/permission/index", &permission.PermissionIndexController{})
	beego.Router(adminRouter+"/permission/delone", &permission.PermissionIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/permission/add", &permission.PermissionAddController{})
	beego.Router(adminRouter+"/permission/edit", &permission.PermissionEditController{})

	beego.Router(adminRouter+"/login", &login.LoginController{})
	beego.Router(adminRouter+"/loginverify", &login.LoginController{}, "post:LoginVerify")
	beego.Router(adminRouter+"/logout", &login.LoginController{}, "get:Logout")

	beego.Router(adminRouter+"/site/index", &siteconfig.SiteConfigIndexController{})
	beego.Router(adminRouter+"/site/delone", &siteconfig.SiteConfigIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/site/add", &siteconfig.SiteConfigAddController{})
	beego.Router(adminRouter+"/site/edit", &siteconfig.SiteConfigEditController{})

	beego.Router(adminRouter+"/information/index", &information.InformationIndexController{})
	beego.Router(adminRouter+"/information/delone", &information.InformationIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/information/add", &information.InformationAddController{})
	beego.Router(adminRouter+"/information/edit", &information.InformationEditController{})

	beego.Router(adminRouter+"/normalquestion/index", &normalquestion.NormalQuestionIndexController{})
	beego.Router(adminRouter+"/normalquestion/delone", &normalquestion.NormalQuestionIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/normalquestion/add", &normalquestion.NormalQuestionAddController{})
	beego.Router(adminRouter+"/normalquestion/edit", &normalquestion.NormalQuestionEditController{})

	beego.Router(adminRouter+"/qicknav/index", &quicknav.QuickNavIndexController{})
	beego.Router(adminRouter+"/qicknav/delone", &quicknav.QuickNavIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/qicknav/add", &quicknav.QuickNavAddController{})
	beego.Router(adminRouter+"/qicknav/edit", &quicknav.QuickNavEditController{})

	beego.Router(adminRouter+"/iplist/index", &iplist.IpListIndexController{})
	beego.Router(adminRouter+"/iplist/delone", &iplist.IpListIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/iplist/add", &iplist.IpListAddController{})

	beego.Router(adminRouter+"/paymentconfig/index", &paymentconfig.PaymentConfigIndexController{})
	beego.Router(adminRouter+"/paymentconfig/delone", &paymentconfig.PaymentConfigIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/paymentconfig/enabled", &paymentconfig.PaymentConfigIndexController{}, "post:Enabled")
	beego.Router(adminRouter+"/paymentconfig/add", &paymentconfig.PaymentConfigAddController{})
	beego.Router(adminRouter+"/paymentconfig/edit", &paymentconfig.PaymentConfigEditController{})
	/* 会员管理 */
	beego.Router(adminRouter+"/member/index", &member.MemberIndexController{})
	beego.Router(adminRouter+"/member/delone", &member.MemberIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/member/locked", &member.MemberIndexController{}, "post:Locked")
	beego.Router(adminRouter+"/member/edit", &member.MemberEditController{})

	beego.Router(adminRouter+"/membersuggest/index", &membersuggest.MemberSuggestIndexController{})
	beego.Router(adminRouter+"/membersuggest/status", &membersuggest.MemberSuggestIndexController{}, "post:Status")

	/* 应用管理 */
	beego.Router(adminRouter+"/appversion/index", &appversion.AppVersionIndexController{})
	beego.Router(adminRouter+"/appversion/delone", &appversion.AppVersionIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/appversion/add", &appversion.AppVersionAddController{})
	beego.Router(adminRouter+"/appversion/edit", &appversion.AppVersionEditController{})

	beego.Router(adminRouter+"/gift/index", &gift.GiftIndexController{})
	beego.Router(adminRouter+"/gift/delone", &gift.GiftIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/gift/add", &gift.GiftAddController{})

	// 前端
	var frontRouter = beego.AppConfig.String("frontrouter")
	beego.Router(frontRouter+"/uploadb", &sysfront.CommonFrontController{}, "post:Upload")
	beego.Router(frontRouter+"/login", &sysfront.LoginFrontController{})
	beego.Router(frontRouter+"/logout", &sysfront.LoginFrontController{}, "get:Logout")
	beego.Router(frontRouter+"/reg", &sysfront.RegFrontController{})
	beego.Router(frontRouter+"/forgetpwd", &sysfront.ForgetPwdFrontController{})
	beego.Router(frontRouter+"/changepwd", &sysfront.ChangePwdFrontController{})

	// api
	var apiRouter = beego.AppConfig.String("apirouter")
	beego.Router(apiRouter+"/login", &sysapi.LoginApiController{})
	beego.Router(apiRouter+"/loginwx", &sysapi.LoginWxApiController{})
	beego.Router(apiRouter+"/loginalipay", &sysapi.LoginAlipayApiController{})
	beego.Router(apiRouter+"/logout", &sysapi.LoginApiController{}, "get:Logout")
	beego.Router(apiRouter+"/bindphone", &sysapi.MemberApiController{}, "post:BindPhone")
	beego.Router(apiRouter+"/cancelaccount", &sysapi.MemberApiController{}, "post:CancelAccount")
	beego.Router(apiRouter+"/refreshlogin", &sysapi.RefreshLoginApiController{})
	beego.Router(apiRouter+"/reg", &sysapi.RegApiController{})
	beego.Router(apiRouter+"/forgetpwd", &sysapi.ForgetPwdApiController{})
	beego.Router(apiRouter+"/changepwd", &sysapi.ChangePwdApiController{})
	beego.Router(apiRouter+"/suggest", &sysapi.MemberSuggestApiController{})
	beego.Router(apiRouter+"/suggest/unread", &sysapi.MemberSuggestApiController{}, "get:GetNewFeedback")
	beego.Router(apiRouter+"/checkupdate", &sysapi.AppVersionApiController{})
	beego.Router(apiRouter+"/info", &sysapi.InformationApiController{})
	beego.Router(apiRouter+"/normalqa", &sysapi.NormalQuestionApiController{})
	beego.Router(apiRouter+"/sysconf", &sysapi.SysConfigApiController{})
	beego.Router(apiRouter+"/member/modifyname", &sysapi.MemberApiController{}, "post:ModifyName")
	beego.Router(apiRouter+"/member/uploadavatar", &sysapi.MemberApiController{}, "post:UploadAvatar")
}
