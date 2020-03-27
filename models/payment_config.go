package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type PaymentConfig struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64     // 创建人Id
	Modifior   int64     // 更新人Id
	Version    int       // 版本
	OrgId      int64     // 组织ID
	AppNo      string    // App编号
	AppName    string    // 应用名称（微信App支付必须）
	PayType    string    // 支付类型
	AppId      string    // AppId
	ConfValue  string    `orm:"size(3072)"` // 配置值，二维码等
	Enabled    int8      // 状态(0：禁用；1：启用)
	Remark     string    `orm:"null"` // 备注
}

// 微信支付配置
type WechatVo struct {
	Appid     string
	AppSecret string
	MchNo     string
	MchKey    string
}

// 支付宝官方接口配置
type AlipayVo struct {
	AppId     string
	PartnerId string
	PriKey    string
	PubKey    string
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(PaymentConfig))
}
