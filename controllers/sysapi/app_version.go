package sysapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	. "github.com/iufansh/iufans/models"
	utils2 "github.com/iufansh/iufans/utils"
)

type AppVersionApiController struct {
	Base2ApiController
}

/*
api app版本比对
param: os=android&ver=1
body:
return:{"code":1,"msg":"成功","data":{"ver":1,"name":"v1.0","desc":"1.版本描述","url":"http://baidu.com","force":1}}
desc:
*/
func (c *AppVersionApiController) Get() {
	defer c.RetJSON()

	// forbiddenArea := GetSiteConfigValue(utils.ScApiIpForbidden)
	// if allowed := iutils.CheckIpAllowed(forbiddenArea, c.Ctx.Input.IP()); !allowed {
	// 	c.Msg = "没有新版"
	// 	return
	// }

	currentVersion, err := c.GetInt("ver", 0)
	if err != nil || currentVersion == 0 {
		c.Msg = "当前版本号为空"
		return
	}
	osType := c.GetString("os")
	appNo := c.AppNo
	if appNo == "" {
		appNo = c.GetString("app")
	}
	auto, _ := c.GetInt("auto", 0)
	o := orm.NewOrm()
	qs := o.QueryTable(new(AppVersion))
	qs = qs.Filter("AppNo", appNo)
	qs = qs.Filter("OsType", osType)
	qs = qs.Filter("VersionNo__gt", currentVersion)
	if auto == 1 {
		qs = qs.Filter("Ignorable", 0)
	}
	qs = qs.Filter("PublishTime__lt", time.Now())
	qs = qs.OrderBy("-VersionNo", "-Id")
	qs = qs.Limit(1)
	var appVersion AppVersion
	if err := qs.One(&appVersion); err != nil {
		if err == orm.ErrNoRows {
			c.Msg = "已是最新版本"
			return
		} else {
			c.Msg = "版本检查异常，请重试"
			return
		}
	}
	c.Code = utils2.CODE_OK
	c.Msg = "检查到新版本：" + appVersion.VersionName
	c.Dta = map[string]interface{}{
		"ver":   appVersion.VersionNo,
		"name":  appVersion.VersionName,
		"desc":  appVersion.VersionDesc,
		"url":   appVersion.DownloadUrl,
		"force": appVersion.ForceUpdate,
	}
}

/*
api app版本比对, 适配XUI-XUpdate
param: app=a&os=android&ver=1
body:
return:
{
  "Code": 0, //0代表请求成功，非0代表失败
  "Msg": "", //请求出错的信息
  "UpdateStatus": 1, //0代表不更新，1代表有版本更新，不需要强制升级，2代表有版本更新，需要强制升级
  "VersionCode": 3,
  "VersionName": "1.0.2",
  "ModifyContent": "1、优化api接口。\r\n2、添加使用demo演示。\r\n3、新增自定义更新服务API接口。\r\n4、优化更新提示界面。",
  "DownloadUrl": "https://raw.githubusercontent.com/xuexiangjys/XUpdate/master/apk/xupdate_demo_1.0.2.apk",
  "ApkSize": 2048
  "ApkMd5": "..."  //md5值没有的话，就无法保证apk是否完整，每次都会重新下载。框架默认使用的是md5加密。
  "IsIgnorable": 1,
}
desc:
*/
func (c *AppVersionApiController) Post() {
	handleCheck(c, false)
}

func (c *AppVersionApiController) PostAuto() {
	handleCheck(c, true)
}

func handleCheck(c *AppVersionApiController, auto bool) {

	// forbiddenArea := GetSiteConfigValue(utils.ScApiIpForbidden)
	// if allowed := iutils.CheckIpAllowed(forbiddenArea, c.Ctx.Input.IP()); !allowed {
	// 	c.Msg = "没有"
	// 	return
	// }

	var code = 1
	var msg string
	var updateStatus int8
	var versionCode int
	var versionName string
	var modifyContent string
	var downloadUrl string
	var apkSize int
	var apkMd5 string
	var isIgnorable bool
	defer func() {
		ret := make(map[string]interface{})
		ret["Code"] = code
		ret["Msg"] = msg
		ret["UpdateStatus"] = updateStatus
		ret["VersionCode"] = versionCode
		ret["VersionName"] = versionName
		ret["ModifyContent"] = modifyContent
		ret["DownloadUrl"] = downloadUrl
		ret["ApkSize"] = apkSize
		ret["ApkMd5"] = apkMd5
		ret["IsIgnorable"] = isIgnorable
		c.Data["json"] = ret
		b, _ := json.Marshal(ret)

		beego.Info("\r\n----------response---------",
			"\r\n", string(b),
			"\r\n-------------------------", )
		c.ServeJSON()
	}()
	type checkUpdateParam struct {
		AppVer     int    `json:"ver"` // 必填
		DeviceOs   string `json:"os"`  // 必填
		AppNo      string `json:"app"` // 必填
		AppChannel string `json:"cha"` // 必填
	}
	var p checkUpdateParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils2.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	osType := p.DeviceOs
	appNo := p.AppNo
	currentVersion := p.AppVer
	appChannel := p.AppChannel
	if currentVersion == 0 {
		msg = "当前版本号为空"
		return
	}
	if appNo == "" {
		msg = "App编号为空"
		return
	}

	o := orm.NewOrm()
	qs := o.QueryTable(new(AppVersion))
	// 渠道不为空时，拼接上渠道
	if appChannel != "" {
		qs = qs.Filter("AppNo__in", appNo, appNo + "-" + appChannel)
	} else {
		qs = qs.Filter("AppNo", appNo)
	}
	if auto {
		qs = qs.Filter("Ignorable", 0)
	}
	qs = qs.Filter("OsType", osType)
	qs = qs.Filter("VersionNo__gt", currentVersion)
	qs = qs.Filter("PublishTime__lt", time.Now())
	qs = qs.OrderBy("-VersionNo", "-AppNo", "-Id")
	qs = qs.Limit(1)
	var appVersion AppVersion
	if err := qs.One(&appVersion); err != nil {
		if err == orm.ErrNoRows {
			code = 0
			msg = "已是最新版本"
			versionCode = currentVersion
			return
		} else {
			msg = "版本检查异常，请重试"
			return
		}
	}
	code = 0
	msg = "检查到新版本：" + appVersion.VersionName
	if appVersion.ForceUpdate == 1 {
		updateStatus = 2
	} else {
		updateStatus = 1
		// 非强制版本才可忽略
		if appVersion.Ignorable == 1 {
			isIgnorable = true
		}
	}
	if !strings.HasPrefix(strings.ToLower(appVersion.DownloadUrl), "http") {
		if c.Ctx.Input.Port() == 80 {
			downloadUrl = c.Ctx.Input.Site() + appVersion.DownloadUrl
		} else {
			downloadUrl = fmt.Sprintf("%s:%d%s", c.Ctx.Input.Site(), c.Ctx.Input.Port(), appVersion.DownloadUrl)
		}
	} else {
		downloadUrl = appVersion.DownloadUrl
	}
	versionCode = appVersion.VersionNo
	versionName = appVersion.VersionName
	modifyContent = appVersion.VersionDesc
	apkSize = appVersion.AppSize / 1000 // 原单位为B,转为KB
	apkMd5 = appVersion.EncryptValue
}