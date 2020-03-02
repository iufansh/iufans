package utils

import (
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/utils/captcha"
)

var cpt *captcha.Captcha

func InitCaptcha() {
	// use beego cache system store the captcha data
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/captcha/", store)
	cpt.ChallengeNums = 4
}

func GetCpt() *captcha.Captcha {
	if cpt == nil {
		InitCaptcha()
	}
	return cpt
}
