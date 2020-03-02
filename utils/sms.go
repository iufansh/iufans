package utils

import (
	"github.com/astaxie/beego"
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

// 发送验证码，2分钟有效
func SendSmsVerifyCode(sender SmsSender) error {
	// 先验证是否1分钟内发过，如果发过，不允许再发
	var cval int
	if err := GetCache("SmsVerifyCode"+sender.Mobile, &cval); err == nil && cval != 0 {
		return errors.New("请1分钟后再发送")
	}
	// 生成验证码
	vc := strconv.FormatInt(int64(RandNum(1000, 9999)), 10)
	SetCache("SmsVerifyCode"+sender.Mobile, vc, 58)
	beego.Info("Send sms verify code to mobile no:", sender.Mobile, ",verify code:", vc)
	if sender.Company == "" {
		sender.Company = "且学科技"
	}
	smsPam := sms.SmsParam{
		Api:    sender.Api,
		Uid:    sender.Uid,
		Key:    sender.Key,
		SignName: sender.Company,
		Mobile: sender.Mobile,
		Text:   vc,
	}
	if num, err := sms.SendSms(smsPam); err != nil {
		beego.Error("SendSmsVerifyCode err:", err)
		return errors.New("发送失败(1)")
	} else if num <= 0 {
		beego.Error("SendSmsVerifyCode err: num=", num)
		return errors.New("发送失败(2)")
	}
	return nil
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