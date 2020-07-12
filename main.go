package main

import (
	_ "blog/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true // 打开session
	//beego.BConfig.WebConfig.StaticDir["/upload"] = "upload"
	beego.SetStaticPath("/static", "static")
	beego.Run()
}

