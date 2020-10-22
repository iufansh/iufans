package initial

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	utils "github.com/iufansh/iutils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/iufansh/iufans/models"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func InitSql() {
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
	maxIdle := 30
	maxConn := 30
	var dataSource string
	dbDriver := beego.AppConfig.String("dbdriver")
	if dbDriver == "mysql" {
		if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
			logs.Error("orm.RegisterDriver mysql err:", err)
			return
		}
		user := beego.AppConfig.String("mysqluser")
		passwd := beego.AppConfig.String("mysqlpass")
		host := beego.AppConfig.String("mysqlurls")
		port, err := beego.AppConfig.Int("mysqlport")
		dbname := beego.AppConfig.String("mysqldb")

		if nil != err {
			port = 3306
		}
		dataSource = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=%s", user, passwd, host, port, dbname, "Asia%2FShanghai")
	} else if dbDriver == "sqlite3" {
		if err := orm.RegisterDriver("sqlite3", orm.DRSqlite); err != nil {
			logs.Error("orm.RegisterDriver sqlite3 err:", err)
			return
		}
		orm.DefaultTimeLoc = time.Local
		dataSource = "./data.db"
	} else {
		panic("未知数据库驱动类型")
	}
	if err := orm.RegisterDataBase("default", dbDriver, dataSource, maxIdle, maxConn); err != nil {
		logs.Error("orm.RegisterDataBase default err:", err)
	}
}

func InitDbFrameData() {
	// 自动建表
	autoCreateDb := beego.AppConfig.DefaultInt("dbautocreate", 0)
	if autoCreateDb == 1 || autoCreateDb == 2 {
		logs.Info("Auto create db")
		isForce := false
		if autoCreateDb == 2 { // drop table 后再建表
			isForce = true
		}
		// 遇到错误立即返回
		err := orm.RunSyncdb("default", isForce, true)
		if err != nil {
			logs.Error("Auto create db error", err.Error())
			return
		}
		if err := initDbData(); err != nil {
			return
		}
	}
}

// 初始化数据库数据--框架数据
func initDbData() error {
	o := orm.NewOrm()
	// 通过 admin表判断数据库是否已初始化，已存在数据，则说明已初始化过。
	if isExist := o.QueryTable(new(Admin)).Exist(); isExist {
		return nil
	}
	logs.Info("Init frame data")
	if err := o.Begin(); err != nil {
		logs.Error("db initDbData orm begin transaction err:", err)
		return err
	}
	// 系统配置
	sc := SiteConfig{Id: 1, Code: "NAME", Value: "公司名称", IsSystem: 1}
	if _, err := o.Insert(&sc); err != nil {
		logs.Warn("Init SiteConfig data error", err)
	}
	pwd := utils.Md5(utils.Md5("111111"), utils.Pubsalt, "17b007bdb8e7af362a1167bcce7277c9")
	// 管理员
	admins := []Admin{
		{Id: 1, Enabled: 1, Locked: 0, IsSystem: 1, LoginFailureCount: 0, Salt: "17b007bdb8e7af362a1167bcce7277c9", Name: "超级管理员", Password: pwd, Username: "superadmin", LoginVerify: 0},
		{Id: 2, Enabled: 1, Locked: 0, IsSystem: 0, LoginFailureCount: 0, Salt: "17b007bdb8e7af362a1167bcce7277c9", Name: "管理员", Password: pwd, Username: "admin", LoginVerify: 0},
	}
	if num, err := o.InsertMulti(len(admins), admins); err != nil {
		logs.Warn("Init Admin data success num:", num, " error:", err)
	}
	// 角色
	roles := []Role{
		{Id: 1, Enabled: 1, Description: "后台管理最高权限", IsSystem: 1, Name: "超级管理员"},
		{Id: 2, Enabled: 1, Description: "后台总管理权限", IsSystem: 0, Name: "后台总管理员"},
		{Id: 10, Enabled: 1, Description: "普通管理权限", IsSystem: 0, Name: "普通管理员", IsOrg: 1},
	}
	if num, err := o.InsertMulti(len(roles), roles); err != nil {
		logs.Warn("Init Role data success num:", num, " error", err)
	}
	// 管理员--角色关联
	adminRoles := []AdminRole{
		{Id: 1, AdminId: 1, RoleId: 1},
		{Id: 2, AdminId: 2, RoleId: 2},
		{Id: 10, AdminId: 2, RoleId: 10},
	}
	if num, err := o.InsertMulti(len(adminRoles), adminRoles); err != nil {
		logs.Warn("Init AdminRole data success num:", num, " error", err)
	}
	// 菜单权限配置
	permissions := []Permission{
		{Id: 1, Pid: 0, Enabled: 1, Display: 0, Description: "系统框架", Url: "BaseIndexController.Get", Name: "系统框架", Icon: "", Sort: 1},
		{Id: 2, Pid: 0, Enabled: 1, Display: 0, Description: "修改密码", Url: "ChangePwdController.Get", Name: "修改密码", Icon: "", Sort: 2},
		{Id: 3, Pid: 0, Enabled: 1, Display: 0, Description: "系统信息", Url: "SysIndexController.Get", Name: "系统信息", Icon: "", Sort: 3},
		{Id: 4, Pid: 3, Enabled: 1, Display: 0, Description: "登录谷歌验证页面", Url: "SysIndexController.GetAuth", Name: "登录谷歌验证页面", Icon: "", Sort: 3},
		{Id: 5, Pid: 3, Enabled: 1, Display: 0, Description: "登录谷歌验证提交", Url: "SysIndexController.PostAuth", Name: "登录谷歌验证提交", Icon: "", Sort: 3},
		{Id: 10, Pid: 0, Enabled: 1, Display: 0, Description: "系统通用-单文件上传", Url: "SyscommonController.Upload", Name: "系统通用-单文件上传", Icon: "", Sort: 10},
		{Id: 11, Pid: 0, Enabled: 1, Display: 0, Description: "系统通用-多文件上传", Url: "SyscommonController.UploadMulti", Name: "系统通用-多文件上传", Icon: "", Sort: 10},
		{Id: 20, Pid: 0, Enabled: 1, Display: 1, Description: "系统设置", Url: "", Name: "系统设置", Icon: "#xe716;", Sort: 100},
		{Id: 21, Pid: 20, Enabled: 1, Display: 1, Description: "管理员", Url: "AdminIndexController.Get", Name: "管理员", Icon: "", Sort: 100},
		{Id: 22, Pid: 21, Enabled: 1, Display: 0, Description: "添加管理员", Url: "AdminAddController.Get", Name: "添加管理员", Icon: "", Sort: 100},
		{Id: 23, Pid: 21, Enabled: 1, Display: 0, Description: "编辑管理员", Url: "AdminEditController.Get", Name: "编辑管理员", Icon: "", Sort: 100},
		{Id: 24, Pid: 21, Enabled: 1, Display: 0, Description: "删除管理员", Url: "AdminIndexController.Delone", Name: "删除管理员", Icon: "", Sort: 100},
		{Id: 25, Pid: 21, Enabled: 1, Display: 0, Description: "锁定解锁管理员", Url: "AdminIndexController.Locked", Name: "锁定解锁管理员", Icon: "", Sort: 100},
		{Id: 26, Pid: 21, Enabled: 1, Display: 0, Description: "解除管理员登录验证", Url: "AdminIndexController.LoginVerify", Name: "解除管理员登录验证", Icon: "", Sort: 100},
		{Id: 30, Pid: 20, Enabled: 1, Display: 1, Description: "角色管理", Url: "RoleIndexController.Get", Name: "角色管理", Icon: "", Sort: 100},
		{Id: 31, Pid: 30, Enabled: 1, Display: 0, Description: "添加角色", Url: "RoleAddController.Get", Name: "添加角色", Icon: "", Sort: 100},
		{Id: 32, Pid: 30, Enabled: 1, Display: 0, Description: "编辑角色", Url: "RoleEditController.Get", Name: "编辑角色", Icon: "", Sort: 100},
		{Id: 33, Pid: 30, Enabled: 1, Display: 0, Description: "删除角色", Url: "RoleIndexController.Delone", Name: "删除角色", Icon: "", Sort: 100},
		{Id: 40, Pid: 20, Enabled: 1, Display: 1, Description: "菜单管理", Url: "PermissionIndexController.Get", Name: "菜单管理", Icon: "", Sort: 100},
		{Id: 41, Pid: 40, Enabled: 1, Display: 0, Description: "添加菜单", Url: "PermissionAddController.Get", Name: "添加菜单", Icon: "", Sort: 100},
		{Id: 42, Pid: 40, Enabled: 1, Display: 0, Description: "编辑菜单", Url: "PermissionEditController.Get", Name: "编辑菜单", Icon: "", Sort: 100},
		{Id: 43, Pid: 40, Enabled: 1, Display: 0, Description: "删除菜单", Url: "PermissionIndexController.Delone", Name: "删除菜单", Icon: "", Sort: 100},
		{Id: 50, Pid: 20, Enabled: 1, Display: 1, Description: "站点配置", Url: "SiteConfigIndexController.Get", Name: "站点配置", Icon: "", Sort: 100},
		{Id: 51, Pid: 50, Enabled: 1, Display: 0, Description: "添加站点配置", Url: "SiteConfigAddController.Get", Name: "添加站点配置", Icon: "", Sort: 100},
		{Id: 52, Pid: 50, Enabled: 1, Display: 0, Description: "编辑站点配置", Url: "SiteConfigEditController.Get", Name: "编辑站点配置", Icon: "", Sort: 100},
		{Id: 53, Pid: 50, Enabled: 1, Display: 0, Description: "删除站点配置", Url: "SiteConfigIndexController.Delone", Name: "删除站点配置", Icon: "", Sort: 100},
		{Id: 60, Pid: 20, Enabled: 1, Display: 1, Description: "快捷导航", Url: "QuickNavIndexController.Get", Name: "快捷导航", Icon: "", Sort: 100},
		{Id: 61, Pid: 60, Enabled: 1, Display: 0, Description: "添加快捷导航", Url: "QuickNavAddController.Get", Name: "添加快捷导航", Icon: "", Sort: 100},
		{Id: 62, Pid: 60, Enabled: 1, Display: 0, Description: "编辑快捷导航", Url: "QuickNavEditController.Get", Name: "编辑快捷导航", Icon: "", Sort: 100},
		{Id: 63, Pid: 60, Enabled: 1, Display: 0, Description: "删除快捷导航", Url: "QuickNavIndexController.Delone", Name: "删除快捷导航", Icon: "", Sort: 100},
		{Id: 70, Pid: 20, Enabled: 1, Display: 1, Description: "组织管理", Url: "OrganizationIndexController.Get", Name: "组织管理", Icon: "", Sort: 90},
		{Id: 71, Pid: 70, Enabled: 1, Display: 0, Description: "添加组织", Url: "OrganizationAddController.Get", Name: "添加组织", Icon: "", Sort: 100},
		{Id: 72, Pid: 70, Enabled: 1, Display: 0, Description: "编辑组织", Url: "OrganizationEditController.Get", Name: "编辑组织", Icon: "", Sort: 100},
		{Id: 73, Pid: 70, Enabled: 1, Display: 0, Description: "删除组织", Url: "OrganizationIndexController.Delone", Name: "删除组织", Icon: "", Sort: 100},
		{Id: 90, Pid: 20, Enabled: 1, Display: 1, Description: "IP黑白名单", Url: "IpListIndexController.Get", Name: "IP黑白名单", Icon: "", Sort: 100},
		{Id: 91, Pid: 90, Enabled: 1, Display: 0, Description: "添加IP黑白名单", Url: "IpListAddController.Get", Name: "添加IP黑白名单", Icon: "", Sort: 100},
		{Id: 92, Pid: 90, Enabled: 1, Display: 0, Description: "删除IP黑白名单", Url: "IpListIndexController.Delone", Name: "删除IP黑白名单", Icon: "", Sort: 100},
		{Id: 100, Pid: 20, Enabled: 1, Display: 1, Description: "支付配置", Url: "PaymentConfigIndexController.Get", Name: "支付配置", Icon: "", Sort: 100},
		{Id: 101, Pid: 100, Enabled: 1, Display: 0, Description: "添加支付配置", Url: "PaymentConfigAddController.Get", Name: "添加支付配置", Icon: "", Sort: 100},
		{Id: 102, Pid: 100, Enabled: 1, Display: 0, Description: "编辑支付配置", Url: "PaymentConfigEditController.Get", Name: "编辑支付配置", Icon: "", Sort: 100},
		{Id: 103, Pid: 100, Enabled: 1, Display: 0, Description: "删除支付配置", Url: "PaymentConfigIndexController.Delone", Name: "删除支付配置", Icon: "", Sort: 100},
		{Id: 104, Pid: 100, Enabled: 1, Display: 0, Description: "启用禁用支付配置", Url: "PaymentConfigIndexController.Enabled", Name: "启用禁用支付配置", Icon: "", Sort: 100},

		{Id: 110, Pid: 20, Enabled: 1, Display: 1, Description: "系统消息", Url: "InformationIndexController.Get", Name: "系统消息", Icon: "", Sort: 100},
		{Id: 111, Pid: 110, Enabled: 1, Display: 0, Description: "添加系统消息", Url: "InformationAddController.Get", Name: "添加系统消息", Icon: "", Sort: 100},
		{Id: 112, Pid: 110, Enabled: 1, Display: 0, Description: "编辑系统消息", Url: "InformationEditController.Get", Name: "编辑系统消息", Icon: "", Sort: 100},
		{Id: 113, Pid: 110, Enabled: 1, Display: 0, Description: "删除系统消息", Url: "InformationIndexController.Delone", Name: "删除系统消息", Icon: "", Sort: 100},

		{Id: 120, Pid: 20, Enabled: 1, Display: 1, Description: "常见问题", Url: "NormalQuestionIndexController.Get", Name: "常见问题", Icon: "", Sort: 100},
		{Id: 121, Pid: 120, Enabled: 1, Display: 0, Description: "添加常见问题", Url: "NormalQuestionAddController.Get", Name: "添加常见问题", Icon: "", Sort: 100},
		{Id: 122, Pid: 120, Enabled: 1, Display: 0, Description: "编辑常见问题", Url: "NormalQuestionEditController.Get", Name: "编辑常见问题", Icon: "", Sort: 100},
		{Id: 123, Pid: 120, Enabled: 1, Display: 0, Description: "删除常见问题", Url: "NormalQuestionIndexController.Delone", Name: "删除常见问题", Icon: "", Sort: 100},

		/* 应用管理 */
		{Id: 200, Pid: 0, Enabled: 1, Display: 1, Description: "应用管理", Url: "", Name: "应用管理", Icon: "#xe653;", Sort: 100},
		{Id: 210, Pid: 200, Enabled: 1, Display: 1, Description: "App版本列表", Url: "AppVersionIndexController.Get", Name: "App版本列表", Icon: "", Sort: 100},
		{Id: 211, Pid: 210, Enabled: 1, Display: 0, Description: "添加App版本", Url: "AppVersionAddController.Get", Name: "添加App版本", Icon: "", Sort: 100},
		{Id: 212, Pid: 210, Enabled: 1, Display: 0, Description: "编辑App版本", Url: "AppVersionEditController.Get", Name: "编辑App版本", Icon: "", Sort: 100},
		{Id: 213, Pid: 210, Enabled: 1, Display: 0, Description: "删除App版本", Url: "AppVersionIndexController.Delone", Name: "删除App版本", Icon: "", Sort: 100},
		// 礼包码
		{Id: 220, Pid: 200, Enabled: 1, Display: 1, Description: "礼包列表", Url: "GiftIndexController.Get", Name: "礼包列表", Icon: "", Sort: 100},
		{Id: 221, Pid: 220, Enabled: 1, Display: 0, Description: "添加礼包", Url: "GiftAddController.Get", Name: "添加礼包", Icon: "", Sort: 100},
		{Id: 222, Pid: 220, Enabled: 1, Display: 0, Description: "删除礼包", Url: "GiftIndexController.Delone", Name: "删除礼包", Icon: "", Sort: 100},

		/* 会员管理 */
		{Id: 300, Pid: 0, Enabled: 1, Display: 1, Description: "会员管理", Url: "", Name: "会员管理", Icon: "#xe770;", Sort: 100},
		{Id: 310, Pid: 300, Enabled: 1, Display: 1, Description: "会员列表", Url: "MemberIndexController.Get", Name: "会员列表", Icon: "", Sort: 100},
		{Id: 311, Pid: 310, Enabled: 1, Display: 0, Description: "编辑会员", Url: "MemberEditController.Get", Name: "编辑会员", Icon: "", Sort: 100},
		{Id: 312, Pid: 310, Enabled: 1, Display: 0, Description: "删除会员", Url: "MemberIndexController.Delone", Name: "删除会员", Icon: "", Sort: 100},
		{Id: 313, Pid: 310, Enabled: 1, Display: 0, Description: "锁定解锁会员", Url: "MemberIndexController.Locked", Name: "锁定解锁会员", Icon: "", Sort: 100},
		// 用户建议
		{Id: 330, Pid: 300, Enabled: 1, Display: 1, Description: "会员反馈列表", Url: "MemberSuggestIndexController.Get", Name: "会员反馈列表", Icon: "", Sort: 100},
		{Id: 331, Pid: 330, Enabled: 1, Display: 0, Description: "会员反馈状态设置", Url: "MemberSuggestIndexController.Status", Name: "会员反馈状态设置", Icon: "", Sort: 100},
	}
	if num, err := o.InsertMulti(len(permissions), permissions); err != nil {
		logs.Warn("Init Permission data success num:", num, " error", err)
	}
	// 角色--权限关联
	rolePermissions := []RolePermission{
		{Id: 1, RoleId: 1, PermissionId: 1},
		{Id: 2, RoleId: 1, PermissionId: 2},
		{Id: 3, RoleId: 1, PermissionId: 3},
		{Id: 4, RoleId: 1, PermissionId: 4},
		{Id: 5, RoleId: 1, PermissionId: 5},
		{Id: 20, RoleId: 1, PermissionId: 20},
		{Id: 21, RoleId: 1, PermissionId: 21},
		{Id: 22, RoleId: 1, PermissionId: 22},
		{Id: 23, RoleId: 1, PermissionId: 23},
		{Id: 24, RoleId: 1, PermissionId: 24},
		{Id: 25, RoleId: 1, PermissionId: 25},
		{Id: 26, RoleId: 1, PermissionId: 26},
		{Id: 30, RoleId: 1, PermissionId: 30},
		{Id: 31, RoleId: 1, PermissionId: 31},
		{Id: 32, RoleId: 1, PermissionId: 32},
		{Id: 33, RoleId: 1, PermissionId: 33},
		{Id: 40, RoleId: 1, PermissionId: 40},
		{Id: 41, RoleId: 1, PermissionId: 41},
		{Id: 42, RoleId: 1, PermissionId: 42},
		{Id: 43, RoleId: 1, PermissionId: 43},
		{Id: 50, RoleId: 1, PermissionId: 50},
		{Id: 51, RoleId: 1, PermissionId: 51},
		{Id: 52, RoleId: 1, PermissionId: 52},
		{Id: 53, RoleId: 1, PermissionId: 53},
		{Id: 60, RoleId: 1, PermissionId: 60},
		{Id: 61, RoleId: 1, PermissionId: 61},
		{Id: 62, RoleId: 1, PermissionId: 62},
		{Id: 63, RoleId: 1, PermissionId: 63},
		{Id: 101, RoleId: 2, PermissionId: 1},
		{Id: 102, RoleId: 2, PermissionId: 2},
		{Id: 103, RoleId: 2, PermissionId: 3},
		{Id: 104, RoleId: 2, PermissionId: 4},
		{Id: 105, RoleId: 2, PermissionId: 5},
		{Id: 110, RoleId: 2, PermissionId: 10},
		{Id: 111, RoleId: 2, PermissionId: 11},
		{Id: 120, RoleId: 2, PermissionId: 20},
		{Id: 121, RoleId: 2, PermissionId: 21},
		{Id: 122, RoleId: 2, PermissionId: 22},
		{Id: 123, RoleId: 2, PermissionId: 23},
		{Id: 124, RoleId: 2, PermissionId: 24},
		{Id: 125, RoleId: 2, PermissionId: 25},
		{Id: 126, RoleId: 2, PermissionId: 26},
		{Id: 150, RoleId: 2, PermissionId: 50},
		{Id: 151, RoleId: 2, PermissionId: 51},
		{Id: 152, RoleId: 2, PermissionId: 52},
		{Id: 153, RoleId: 2, PermissionId: 53},
		/*
		{Id: 160, RoleId: 2, PermissionId: 60},
		{Id: 161, RoleId: 2, PermissionId: 61},
		{Id: 162, RoleId: 2, PermissionId: 62},
		{Id: 163, RoleId: 2, PermissionId: 63},
		 */
		{Id: 170, RoleId: 2, PermissionId: 70},
		{Id: 171, RoleId: 2, PermissionId: 71},
		{Id: 172, RoleId: 2, PermissionId: 72},
		{Id: 173, RoleId: 2, PermissionId: 73},
		{Id: 190, RoleId: 2, PermissionId: 90},
		{Id: 191, RoleId: 2, PermissionId: 91},
		{Id: 192, RoleId: 2, PermissionId: 92},
		{Id: 200, RoleId: 2, PermissionId: 100},
		{Id: 201, RoleId: 2, PermissionId: 101},
		{Id: 202, RoleId: 2, PermissionId: 102},
		{Id: 203, RoleId: 2, PermissionId: 103},
		{Id: 204, RoleId: 2, PermissionId: 104},
		{Id: 210, RoleId: 2, PermissionId: 110},
		{Id: 211, RoleId: 2, PermissionId: 111},
		{Id: 212, RoleId: 2, PermissionId: 112},
		{Id: 213, RoleId: 2, PermissionId: 113},
		{Id: 220, RoleId: 2, PermissionId: 120},
		{Id: 221, RoleId: 2, PermissionId: 121},
		{Id: 222, RoleId: 2, PermissionId: 122},
		{Id: 223, RoleId: 2, PermissionId: 123},

		/* 应用管理 */
		/* 默认不开启
		{Id: 300, RoleId: 2, PermissionId: 200},
		{Id: 310, RoleId: 2, PermissionId: 210},
		{Id: 311, RoleId: 2, PermissionId: 211},
		{Id: 312, RoleId: 2, PermissionId: 212},
		{Id: 313, RoleId: 2, PermissionId: 213},

		{Id: 320, RoleId: 2, PermissionId: 220},
		{Id: 321, RoleId: 2, PermissionId: 221},
		{Id: 322, RoleId: 2, PermissionId: 222},
		*/
		/* 会员管理 */
		{Id: 400, RoleId: 2, PermissionId: 300},
		{Id: 410, RoleId: 2, PermissionId: 310},
		{Id: 411, RoleId: 2, PermissionId: 311},
		{Id: 412, RoleId: 2, PermissionId: 312},
		{Id: 413, RoleId: 2, PermissionId: 313},
		/*
		{Id: 430, RoleId: 2, PermissionId: 330},
		{Id: 431, RoleId: 2, PermissionId: 331},
		 */
		/* 普通管理员*/
		{Id: 1001, RoleId: 10, PermissionId: 1},
		{Id: 1002, RoleId: 10, PermissionId: 2},
		{Id: 1003, RoleId: 10, PermissionId: 3},
		{Id: 1004, RoleId: 10, PermissionId: 4},
		{Id: 1005, RoleId: 10, PermissionId: 5},
		{Id: 1010, RoleId: 10, PermissionId: 10},
		{Id: 1020, RoleId: 10, PermissionId: 20},
		{Id: 1021, RoleId: 10, PermissionId: 21},
		{Id: 1022, RoleId: 10, PermissionId: 22},
		{Id: 1023, RoleId: 10, PermissionId: 23},
		{Id: 1024, RoleId: 10, PermissionId: 24},
		{Id: 1025, RoleId: 10, PermissionId: 25},
		{Id: 1026, RoleId: 10, PermissionId: 26},

		{Id: 1090, RoleId: 10, PermissionId: 90},
		{Id: 1091, RoleId: 10, PermissionId: 91},
		{Id: 1092, RoleId: 10, PermissionId: 92},
	}
	if num, err := o.InsertMulti(len(rolePermissions), rolePermissions); err != nil {
		logs.Warn("Init RolePermission data success num:", num, " error", err)
	}
	if err := o.Commit(); err != nil {
		logs.Error("db initDbData orm Commit transaction err:", err)
	}
	return nil
}
