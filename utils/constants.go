package utils

// 站点配置code
const (
	Scname                  = "NAME"
	Sccompanyname           = "COMPANY_NAME"
	Sccompanyaddress        = "COMPANY_ADDRESS"
	Sccompanyconcattel      = "COMPANY_CONCAT_TEL"
	Sccompanyconcatqq       = "COMPANY_CONCAT_QQ"
	Sccompanyconcatwx       = "COMPANY_CONCAT_WX"
	Scsmssignname           = "sms_signname"
	Scsmsapi                = "sms_api"
	Scsmsuid                = "sms_uid"
	Scsmskey                = "sms_key"
	Scfrontregsmsverify     = "front_reg_sms_verify"
	ScBaiduApiKey           = "baidu_api_key"
	ScBaiduSecretKey        = "baidu_secret_key"
	ScAliyunAccessKeyId     = "aliyun_access_key_id"
	ScAliyunAccessKeySecret = "aliyun_access_key_secret"
	ScAliyunAuthSecret      = "aliyun_auth_secret_" // 阿里云认证服务的秘钥，每个App不同，添加后缀 App编号，如：aliyun_auth_secret_s
)

var SiteConfigCodeMap = map[string]string{
	"DIY":  "自定义",
	Scname: "站点名称",
}

func GetSiteConfigCodeMap() map[string]string {
	return SiteConfigCodeMap
}

const (
	CODE_OK         = 1
	CODE_NEED_LOGIN = 11
	CODE_ERROR      = 21
)

const (
	PayTypeAlipay    = "alipay"
	PayTypeWechatPay = "wechatpay"
	PayTypeApplePay = "applepay"
)

const (
	AppChannelWxa = "wxa" // 微信小程序
)