package controllers

import (
	"fmt"
	"blog/models"
)

type HomeController struct {
	//beego.Controller
	BaseController
}

func (this *HomeController)Get(){
	//this.Prepare()
	//loginuser = this.GetSession("loginuser")
	//if loginuser != nil {
	//	this.IsLogin = true
	//	this.Loginuser = loginuser
	//} else {
	//	this.IsLogin = false
	//}
	//this.Data["IsLogin"] = this.IsLogin

	page, _ := this.GetInt("page")//get请求时获取url后面的参数，即127.0.0.1:8080/home?page=...
	//fmt.Println("==========",page)
	if page <= 0 {
		page = 1
	}
	var artList []models.Article
	artList, _ = models.FindArticleWithPage(page)
	this.Data["PageCode"] = models.ConfigHomeFooterPageCode(page)
	this.Data["HasFooter"] = true

	fmt.Println("IsLogin:", this.IsLogin, this.Loginuser)
	this.Data["Content"] = models.MakeHomeBlocks(artList, this.IsLogin,this.Loginuser)

	this.TplName = "home.html"
	//fmt.Println("IsLogin:",this.IsLogin,this.Loginuser)
	//this.IsLogin = false
	//this.Data["IsLogin"] = false
	//this.TplName="home.html"
}