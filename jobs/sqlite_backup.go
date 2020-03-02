package jobs

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"github.com/iufansh/iutils"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func InitSqliteBackup() {
	if beego.AppConfig.String("dbdriver") != "sqlite3" {
		return
	}
	backupPath := beego.AppConfig.String("sqlite3backuppath")
	if backupPath == "" {
		return
	}
	if !strings.HasSuffix(backupPath, "/") && !strings.HasSuffix(backupPath, "\\") {
		backupPath = backupPath + "/"
	}
	tk1 := toolbox.NewTask("SysSqliteBackup", "15 15 03 * * *", func() error {
		beego.Info("SysSqliteBackup start")
		size, err := iutils.CopyFile("./data.db", backupPath+time.Now().Format("20060102150405")+".db")
		if err != nil {
			beego.Error("SysSqliteBackup err:", err)
		} else {
			beego.Info("SysSqliteBackup size =", size)
		}
		files, _ := ioutil.ReadDir(backupPath)
		for i,file := range files {
			if i >= len(files) - 7 {
				break
			}
			if err := os.Remove(backupPath+file.Name()); err != nil {
				beego.Error("SysSqliteBackup delete old file =", file.Name(), " err:", err)
			} else {
				beego.Info("SysSqliteBackup delete old file =", file.Name())
			}
		}
		beego.Info("SysSqliteBackup finish")
		return nil
	})
	toolbox.AddTask("SysSqliteBackup", tk1)
	toolbox.StartTask()
}
