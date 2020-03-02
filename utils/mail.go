package utils

import (
	. "github.com/iufansh/iutils"
	"gopkg.in/gomail.v2"
	"github.com/astaxie/beego"
)

var (
	sender   string
	host     string
	port     int
	username string
	password string
)

// 使用当前邮件工具时，必须先调用initial文件夹下的InitMailConf初始化。
func InitMail(s string, h string, p int, u string, pw string) {
	sender = s
	host = h
	port = p
	username = u
	password = pw
}

type MailSender struct {
	//From string // 发件人，可空，为空时取app.conf中的配置
	To      []string // 收件人列表，不可空
	Cc      string   // 抄送，可空
	Subject string   // 主题，可空，默认：无主题
	Body    string   // 内容，内容，不可空
	Attach  string   // 附件，可空
}

// 发送邮件，一次性全部发送
func (this *MailSender) Send() error {
	//beego.Info("Mail send, ", this, host, port, username, password)
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", this.To...)
	if this.Cc != "" {
		m.SetAddressHeader("Cc", this.Cc, "cc")
	}
	m.SetHeader("Subject", this.Subject)
	m.SetBody("text/html", this.Body)
	if this.Attach != "" {
		m.Attach(this.Attach)
	}

	d := gomail.NewDialer(host, port, username, password)

	if err := d.DialAndSend(m); err != nil {
		beego.Error("MailSender send error", err.Error())
		return err
	}
	beego.Info("Mail send success")
	return nil
}

// 拆分收件人，一人发送一次；返回成功发送数量
func (this *MailSender) SendSeparate() (int, error) {
	var num int
	d := gomail.NewDialer(host, port, username, password)
	s, err := d.Dial()
	if err != nil {
		beego.Error("MailSender SendSeparate Dial error", err.Error())
		return 0, err
	}

	m := gomail.NewMessage()
	for _, v := range this.To {
		m.SetHeader("From", sender)
		m.SetAddressHeader("To", v, v)
		m.SetHeader("Subject", this.Subject)
		m.SetBody("text/html", this.Body)

		if err := gomail.Send(s, m); err != nil {
			beego.Error("MailSender SendSeparate error", err.Error())
		} else {
			num++
		}
		m.Reset()
	}
	return num, nil
}

// 发送验证码，5分钟有效
func SendMailVerifyCode(to string) error {
	// 生成验证码
	vc := RandStringLower(4)
	ms := MailSender{To: []string{to},
		Subject: "Phage(富吉)系统验证码",
		Body: vc}

	err := ms.Send()
	if err != nil {
		return err
	} else {
		SetCache("MailVerifyCode"+to, vc, 300)
		beego.Info("Send mail verify code to:", to, "verify code:", vc)
		return nil
	}
}

func VerifyMailVerifyCode(key string, vc string) bool {
	if key == "" {
		return false
	}
	var cval string
	if err := GetCache("MailVerifyCode"+key, &cval); err != nil {
		return false
	}
	if cval != vc {
		return false
	}
	return true
}
