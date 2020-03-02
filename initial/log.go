package initial

import (
	"fmt"

	"github.com/astaxie/beego"
)

func InitLog() {
	filename := beego.AppConfig.String("logfilename")
	maxdays, err1 := beego.AppConfig.Int("logmaxdays")
	level, err2 := beego.AppConfig.Int("loglevel")
	if nil != err1 {
		maxdays = 7
	}
	if nil != err2 {
		level = beego.LevelInformational
	}
	beego.SetLogger("file", fmt.Sprintf(`{"filename":"%s","daily":true,"maxdays":%d}`, filename, maxdays))
	beego.SetLevel(level)
	beego.SetLogFuncCall(true)
}
