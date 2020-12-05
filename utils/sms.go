package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/iufansh/iuplugins/sms"
	. "github.com/iufansh/iutils"
	"github.com/pkg/errors"
	"strconv"
)

type SmsSender struct {
	Api     string // api方；1：http://sms.webchinese.cn; 阿里云模板ID（）
	Uid     string
	Key     string
	Mobile  string
	Company string
}

/**
 * 发送验证码，1分钟有效
 * return 发送的验证码，异常
 */
func SendSmsVerifyCode(sender SmsSender) (string, error) {
	// 先验证是否1分钟内发过，如果发过，不允许再发
	var cval int
	if err := GetCache("SmsVerifyCode"+sender.Mobile, &cval); err == nil && cval != 0 {
		return "", errors.New("请稍后再发送")
	}
	// 生成验证码
	vc := strconv.FormatInt(int64(RandNum(1000, 9999)), 10)
	if err := SetCache("SmsVerifyCode"+sender.Mobile, vc, 60); err != nil {
		logs.Error("SendSmsVerifyCode SetCache err:", err)
		return vc, errors.New("发送失败")
	}
	logs.Info("SendSmsVerifyCode to mobile:", sender.Mobile, ",verifyCode:", vc)

	if beego.BConfig.RunMode == "dev" { // 如果是开发模式，直接返回验证码
		//return "测试验证码：" + vc, nil
	}
	smsPam := sms.SmsParam{
		Api:      sender.Api,
		Uid:      sender.Uid,
		Key:      sender.Key,
		SignName: sender.Company,
		Mobile:   sender.Mobile,
		Text:     vc,
	}
	if num, err := sms.SendSms(smsPam); err != nil {
		logs.Error("SendSmsVerifyCode err:", err)
		return vc, errors.New("发送失败")
	} else if num <= 0 {
		logs.Error("SendSmsVerifyCode err: num=", num)
		return vc, errors.New("发送失败")
	}
	return vc, nil
}

func VerifySmsVerifyCode(mobile string, vc string) bool {
	if mobile == "" {
		return false
	}
	var cval string
	if err := GetCache("SmsVerifyCode"+mobile, &cval); err != nil {
		return false
	}
	if cval != vc {
		return false
	}
	return true
}
