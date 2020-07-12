package controllers

import (
	"strconv"
	"fmt"
	"blog/models"
	"blog/utils"
	"io/ioutil"
	"html/template"
	//"html"
)
type ShowArticleController struct {
	//beego.Controller
	BaseController
}

func (this *ShowArticleController) Get() {
	idStr := this.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	fmt.Println("id:", id)

	//获取id所对应的文章信息
	art := models.QueryArticleWithId(id)
	if art.Content == ""{
		this.Data["Title"] = art.Title
		//this.Data["Content"] = `<div><h1>文件已被删除</h1></div>`//c.Ctx.WriteString(string(data))
		//this.Ctx.WriteString("<div><h1>文件已被删除</h1></div>")
		//this.Data["Content"] = utils.SwitchMarkdownToHtml("<div><h1>文件已被删除</h1></div>")
		this.Data["Content"] = template.HTML("<div><h1>文件已被删除</h1></div>")

	}else{
		fileread, err := ioutil.ReadFile(art.Content)
		if err != nil{
			fmt.Println("读取文件出错",err)
		}

		this.Data["Title"] = art.Title
		//this.Data["Content"] = art.Content

		this.Data["Content"] = utils.SwitchMarkdownToHtml(string(fileread))
	}

	this.TplName="show_article.html"
}
