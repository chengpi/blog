package models

import (
	"strconv"
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"github.com/astaxie/beego"
)

type HomeBlockParam struct {
	Id         int
	Title      string
	Tags       [] TagLink
	Short      string
	Content    string
	Author     string
	CreateTime string
	//查看文章的地址
	Link string

	//修改文章的地址
	UpdateLink string
	DeleteLink string

	//记录是否登录
	IsLogin bool

	IsFirst bool
	IsBelong bool
}

type TagLink struct {
	TagName string
	TagUrl  string
}
//----------首页显示内容---------
func MakeHomeBlocks(articles []Article, isLogin bool,loginuser interface{}) template.HTML {
	htmlHome := ""
	for _, art := range articles {
		//将数据库model转换为首页模板所需要的model
		homeParam := HomeBlockParam{}
		homeParam.Id = art.Id
		homeParam.Title = art.Title
		homeParam.Tags = createTagsLinks(art.Tags)
		//fmt.Println("tag-->", art.Tags)
		homeParam.Short = art.Short
		homeParam.Content = art.Content
		homeParam.Author = art.Author
		//homeParam.CreateTime = utils.SwitchTimeStampToData(art.Createtime)
		homeParam.CreateTime =art.Createtime
		homeParam.Link = "/article/" + strconv.Itoa(art.Id)
		homeParam.UpdateLink = "/article/update?id=" + strconv.Itoa(art.Id)
		homeParam.DeleteLink = "/article/delete?id=" + strconv.Itoa(art.Id)
		homeParam.IsLogin = isLogin

		if (art.Id-1)%6 == 0{
			homeParam.IsFirst = true
		}else{
			homeParam.IsFirst = false
		}
		if loginuser != nil{
			if homeParam.Author == loginuser.(string){
				homeParam.IsBelong = true
			}else{
				homeParam.IsBelong = false
			}
		}

		//处理变量
		//ParseFile解析该文件，用于插入变量
		t, _ := template.ParseFiles("views/block/home_block.html")
		buffer := bytes.Buffer{}
		//就是将html文件里面的比那两替换为穿进去的数据
		t.Execute(&buffer, homeParam)
		htmlHome += buffer.String()
	}
	//fmt.Println("htmlHome-->",htmlHome)
	return template.HTML(htmlHome)
}
//将tags字符串转化成首页模板所需要的数据结构
func createTagsLinks(tags string) []TagLink {
	var tagLink [] TagLink
	tagsPamar := strings.Split(tags, "&")//切割字符串
	for _, tag := range tagsPamar {
		tagLink = append(tagLink, TagLink{tag, "/home?tag=" + tag})
	}
	return tagLink
}

type HomeFooterPageCode struct {
	HasPre bool
	HasNext bool
	ShowPage string
	PreLink string
	NextLink string
}
//-----------翻页-----------
//page是当前的页数
func ConfigHomeFooterPageCode(page int) HomeFooterPageCode {
	pageCode := HomeFooterPageCode{}
	//查询出总的条数
	num := GetArticleRowsNum()
	//从配置文件中读取每页显示的条数
	pageRow, _ := beego.AppConfig.Int("articleListPageNum")
	//计算出总页数
	allPageNum := (num-1)/pageRow + 1//不可以是num/pageRow+1,因为当num可以被pageRow整除时，这个式子不能满足

	pageCode.ShowPage = fmt.Sprintf("%d/%d", page, allPageNum)

	//当前页数小于等于1，那么上一页的按钮不能点击
	if page <= 1 {
		pageCode.HasPre = false
	} else {
		pageCode.HasPre = true
	}

	//当前页数大于等于总页数，那么下一页的按钮不能点击
	if page >= allPageNum {
		pageCode.HasNext = false
	} else {
		pageCode.HasNext = true
	}

	pageCode.PreLink = "/home?page=" + strconv.Itoa(page-1)
	pageCode.NextLink = "/home?page=" + strconv.Itoa(page+1)
	return pageCode

}