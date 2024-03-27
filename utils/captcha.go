package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/utils/captcha"
	"strings"
)

var cpt *captcha.Captcha

func InitCaptcha() {
	// use beego cache system store the captcha data
	var domainUri = beego.AppConfig.String("domainuri")
	if domainUri != "" && !strings.HasPrefix(domainUri, "/") {
		domainUri = "/" + domainUri
	}
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter(domainUri + "/captcha/", store)
	cpt.ChallengeNums = 4
}

func GetCpt() *captcha.Captcha {
	if cpt == nil {
		InitCaptcha()
	}
	return cpt
}
