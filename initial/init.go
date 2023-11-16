package initial

import (
	"github.com/astaxie/beego"
	. "github.com/iufansh/iufans/utils"
	"strings"
)

func init() {
	InitLog()
	InitSql()
	InitCache()
	InitFilter()
	InitSysTemplateFunc()

	domainUri := beego.AppConfig.String("domainuri")
	if domainUri != "" {
		if !strings.HasPrefix(domainUri, "/") {
			domainUri = "/" + domainUri
		}
		beego.SetStaticPath( domainUri + "/static", "static")
	}
}