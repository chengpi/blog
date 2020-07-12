package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"blog/models"
	"blog/utils"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	//需要先获取一下session，以免c.CruSession未空
	if this.CruSession == nil {
		//this.GetSession("loginuser")
		loginuser = this.GetSession("loginuser")
	}
	if loginuser != nil{
		//this.Data["json"] = map[string]interface{}{"code": 1, "message": "登录成功","status":"success"}
		//this.ServeJSON()
		fmt.Println("++++++",loginuser)
		//this.TplName = "home.html"
		this.Redirect("/home",302)
	}else {
		this.TplName = "login.html"
	}

}
func (this *LoginController) Post() {
	username := this.GetString("username")
	password := this.GetString("password")
	fmt.Println("username:", username, ",password:", password)

	id := models.QueryUserWithParam(username, utils.MD5(password))
	fmt.Println("id:",id)
	if id > 0 {
		/*
        设置了session后会将数据处理设置到cookie，然后再浏览器进行网络请求的时候会自动带上cookie
        因为我们可以通过获取这个cookie来判断用户是谁，这里我们使用的是session的方式进行设置
         */
		this.SetSession("loginuser", username)
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "登录成功","status":"success"}
	} else {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "登录失败","status":"fail"}
	}
	this.ServeJSON()
}