package controllers

import (
	"fmt"
	"blog/models"
	"time"
	"os"
)

type AddArticleController struct {
	BaseController
}

/*
当访问/add路径的时候回触发AddArticleController的Get方法
响应的页面是通过TpName
 */
func (this *AddArticleController) Get() {
	
	this.TplName = "write_article.html"
}

//通过this.ServerJSON()这个方法去返回json字符串
func (this *AddArticleController) Post() {

	//获取浏览器传输的数据，通过表单的name属性获取值
	title := this.GetString("title")
	tags := this.GetString("tags")
	short := this.GetString("short")
	content := this.GetString("content")
	fmt.Printf("title:%s,tags:%s\n", title, tags)
	
	var response map[string]interface{}
	author,ok := loginuser.(string)
	if !ok{
		fmt.Println("用户未登录...")
		response = map[string]interface{}{"code": 0, "message": "用户未登录","status":"fail"}
		this.Data["json"] = response
		this.ServeJSON()
		return
	}
	//filename := "F:/go/src/blog/static/file/"+ title +".md"
	filename := "static/file/"+ title +".md"
	file,err:=os.OpenFile(filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,os.ModePerm)
	if err != nil {
		fmt.Println("打开文件有误：", err.Error())
		return
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil{
		fmt.Println("写入内容有误",err)
		return
	}
	//实例化model，将它出入到数据库中
	art := models.Article{
		0, 
		title, 
		tags, 
		short, 
		filename, 
		author, 
		time.Now().Format("2006/1/02 15:04:05")}//time.Now().Unix()
	_, err = models.AddArticle(art)

	//返回数据给浏览器
	if err == nil {
		//无误
		response = map[string]interface{}{"code": 1, "message": "ok","status":"success"}
	} else {
		response = map[string]interface{}{"code": 0, "message": "error","status":"fail"}
	}

	this.Data["json"] = response
	this.ServeJSON()

}