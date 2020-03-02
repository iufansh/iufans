# 基于beego的web框架

项目依赖：
go get github.com/astaxie/beego
go get github.com/go-sql-driver/mysql
go get github.com/skip2/go-qrcode
go get github.com/gomodule/redigo/redis
go get gopkg.in/gomail.v2
go get github.com/mattn/go-sqlite3
 

# 使用说明，后台管理系统必须引用，api和front根据需要引用。
a.后台管理系统引用方法
    1.项目中加载初始化
        InitLog()
        InitSql()
        InitBeeCache()
        InitFilter()
        InitSysTemplateFunc()

    2.routers中import _ "github.com/iufansh/iufans/routers"

b.api和front引用方法
    1.项目router中init api和front的路由
    2.static目录下要放front文件夹，front下内容：css、js、theme-xxx、images/avatar、
    3.views目录下要放front文件夹，内容根据front的路由页面创建


## 关于api返回code说明
0：普通错误、异常，终端提示即可
1：成功处理
11：需要登录的接口未登录，终端必须登录才能调用
21：接口调用异常，对接时就要处理好，上线后归为终端系统异常
