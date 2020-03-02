package sysfront

import (
	"fmt"
	"github.com/astaxie/beego"
	utils "github.com/iufansh/iutils"
	"net/url"
	"os"
	"strings"
	"time"
)

type CommonFrontController struct {
	BaseFrontController
}

func (c *CommonFrontController) Upload() {
	defer c.RetJSON()
	f, h, err := c.GetFile("file")
	defer f.Close()
	if err != nil {
		beego.Error("Recharge upload file get file error1", err)
		c.Msg = "上传失败，请重试(1)"
		return
	}
	fname := url.QueryEscape(h.Filename)
	suffix := utils.SubString(fname, len(fname), strings.LastIndex(fname, ".")-len(fname))
	fmt.Println(suffix)
	if strings.ToLower(suffix) != ".jpg" && strings.ToLower(suffix) != ".jpeg" && strings.ToLower(suffix) != ".png" {
		c.Msg = "图片必须为jpg、jpeg或png格式"
		return
	}
	uploadPath := fmt.Sprintf("upload/front/%d/%s/%d/", time.Now().Year(), time.Now().Month().String(), time.Now().Day())
	if flag, _ := utils.PathExists(uploadPath); !flag {
		if err2 := os.MkdirAll(uploadPath, 0644); err2 != nil {
			beego.Error("Recharge upload file get file error2", err2)
			c.Msg = "上传失败，请重试(2)"
			return
		}
	}
	var uploadName = fmt.Sprintf("%s%s_%d%s", uploadPath, c.LoginMemberName, time.Now().Unix(), suffix)
	err3 := c.SaveToFile("file", uploadName)
	if err3 != nil {
		beego.Error("Recharge upload file save file error3", err3)
		c.Msg = "上传失败，请重试(3)"
		return
	}
	c.Code = 0
	c.Msg = "上传成功"
	c.Dta = "/" + uploadName
}
