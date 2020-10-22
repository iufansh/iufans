package utils

// 站点配置code
const (
	Scname              = "NAME"
	Sccompanyname       = "COMPANY_NAME"
	Sccompanyaddress    = "COMPANY_ADDRESS"
	Sccompanyconcattel  = "COMPANY_CONCAT_TEL"
	Sccompanyconcatqq   = "COMPANY_CONCAT_QQ"
	Scsmssignname       = "sms_signname"
	Scsmsapi            = "sms_api"
	Scsmsuid            = "sms_uid"
	Scsmskey            = "sms_key"
	Scfrontregsmsverify = "front_reg_sms_verify"
	ScBaiduApiKey       = "baidu_api_key"
	ScBaiduSecretKey    = "baidu_secret_key"
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
)
