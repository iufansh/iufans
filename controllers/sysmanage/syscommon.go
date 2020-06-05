package sysmanage

import (
	"fmt"
	utils2 "github.com/iufansh/iufans/utils"
	utils "github.com/iufansh/iutils"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type SyscommonController struct {
	BaseController
}

// 后台通用文件上传接口，上传的文件默认保存于upload目录下。静态资源。
func (c *SyscommonController) Upload() {
	defer c.RetJSON()
	nameMode := c.GetString("nameMode", "0")
	saveDir := c.GetString("saveDir")
	isMd5, _ := c.GetInt8("md5", 0)
	f, h, err := c.GetFile("file")
	defer f.Close()
	if err != nil {
		beego.Error("Syscommon upload file get file error", err)
		c.Msg = "上传失败，请重试(1)"
		return
	}
	var uploadPath string
	if saveDir != "" {
		uploadPath = strings.TrimPrefix(strings.TrimSuffix(saveDir, "/"), "/") + "/"
	} else {
		uploadPath = fmt.Sprintf("upload/%d/%s/%d/", time.Now().Year(), time.Now().Month().String(), time.Now().Day())
	}
	if flag, _ := utils.PathExists(uploadPath); !flag {
		if err2 := os.MkdirAll(uploadPath, 0644); err2 != nil {
			beego.Error("Syscommon upload file get file error", err2)
			c.Msg = "上传失败，请重试(2)"
			return
		}
	}

	fName := url.QueryEscape(h.Filename)
	suffix := utils.SubString(fName, len(fName), strings.LastIndex(fName, ".")-len(fName))
	var saveName string
	if nameMode == "0" { // 随机名称
		saveName = utils.Md5(utils.GetGuid(), strconv.FormatInt(time.Now().UnixNano(), 16)) + suffix
	} else if nameMode == "1" {
		saveName = fName
	} else {
		saveName = nameMode + suffix
	}
	uploadName := uploadPath + saveName
	err3 := c.SaveToFile("file", uploadName)
	if err3 != nil {
		beego.Error("Syscommon upload file save file error2", err3)
		c.Msg = "上传失败，请重试(3)"
		return
	}

	var md5Value string
	if isMd5 == 1 {
		md5Value = utils.Md5File(f)
	}
	c.Msg = "上传成功"
	c.Code = utils2.CODE_OK
	c.Dta = map[string]interface{}{
		"src":      "/" + uploadName,
		"name":     saveName,
		"size":     h.Size, // 单位：B
		"md5Value": md5Value,
	}
}

// 后台通用多文件上传接口，上传的文件默认保存于upload目录下。静态资源。
func (c *SyscommonController) UploadMulti() {
	defer c.RetJSON()
	nameMode := c.GetString("nameMode", "0")
	saveDir := c.GetString("saveDir")
	fs, err := c.GetFiles("file")
	if err != nil {
		beego.Error("Syscommon UploadMulti file get file error", err)
		c.Msg = "上传失败，请重试(1)"
		return
	}
	var uploadPath string
	if saveDir != "" {
		uploadPath = strings.TrimPrefix(strings.TrimSuffix(saveDir, "/"), "/") + "/"
	} else {
		uploadPath = fmt.Sprintf("upload/%d/%s/%d/", time.Now().Year(), time.Now().Month().String(), time.Now().Day())
	}
	if flag, _ := utils.PathExists(uploadPath); !flag {
		if err2 := os.MkdirAll(uploadPath, 0644); err2 != nil {
			beego.Error("Syscommon UploadMulti file get file error", err2)
			c.Msg = "上传失败，请重试(2)"
			return
		}
	}

	for _, file := range fs {
		fName := url.QueryEscape(file.Filename)
		var saveName string
		if nameMode == "0" { // 随机名称
			suffix := utils.SubString(fName, len(fName), strings.LastIndex(fName, ".")-len(fName))
			saveName = utils.Md5(utils.GetGuid(), strconv.FormatInt(time.Now().UnixNano(), 16)) + suffix
		} else {
			saveName = fName
		}

		toFile := uploadPath + saveName

		f, err := os.OpenFile(toFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			beego.Error("Syscommon UploadMulti OpenFile error", err)
			c.Msg = "上传失败，请重试(3)"
			return
		}
		defer f.Close()
		fi, err := file.Open()
		if err != nil {
			beego.Error("Syscommon UploadMulti file Open error", err)
			c.Msg = "上传失败，请重试(4)"
			return
		}
		defer fi.Close()
		io.Copy(f, fi)
	}

	c.Msg = "上传成功"
	c.Code = utils2.CODE_OK
}

func (c *SyscommonController) Download() {

}
