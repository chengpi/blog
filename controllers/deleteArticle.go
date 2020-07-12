package controllers

import (
	"fmt"
	"blog/models"
	"log"
)

type DeleteArticleController struct{
	BaseController
}
func (this *DeleteArticleController) Get() {
	id, _ := this.GetInt("id")
	//fmt.Println(id)
	var response map[string]interface{}
	_,ok := loginuser.(string)
	if !ok{
		fmt.Println("用户未登录......")
		response = map[string]interface{}{"code": 0, "message": "用户当前为未登录状态","status":"fail"}
		this.Data["json"] = response
		this.ServeJSON()
		return
	}

	_, err := models.DeleteArticle(id)
	if err != nil {
		log.Println(err)
	}

	this.Redirect("/home", 302)
	//this.TplName = "home.html"

}