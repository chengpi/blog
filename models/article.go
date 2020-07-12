package models

import (
	"blog/utils"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"os"
	"sync"
	//"io"
)

type Article struct {
	Id         int
	Title      string
	Tags       string
	Short      string
	Content    string
	Author     string
	//Createtime int64
	Createtime string
	//Status int //Status=0为正常，1为删除，2为冻结
}
var rwLock sync.RWMutex
//-----------查询文章---------

//根据页码查询文章
func FindArticleWithPage(page int) ([]Article, error) {
	//从配置文件中获取每页的文章数量
	num, _ := beego.AppConfig.Int("articleListPageNum")
	//fmt.Println("num值:",num)
	page--
	//fmt.Println("---------->page", page)
	return QueryArticleWithPage(page, num)
}
/**
分页查询数据库
limit分页查询语句，
    语法：limit m，n

    m代表从多少位开始获取，与id值无关
    n代表获取多少条数据

注意limit前面没有where
 */
func QueryArticleWithPage(page, num int) ([]Article, error) {
	sql := fmt.Sprintf("limit %d,%d", page*num, num)//分页查询其实就是当前端每显示一个分页时就对应到数据库里查询其中某一部分
	return QueryArticlesWithCon(sql)
}
//第一个分页 0*6 6 =>0 6 代表从数据库里查询序号为0到5 这6条显示到第1个分页上
//第二个分页 1*6 6 =>6 6 代表从数据库里查询序号为6到11 这6条显示到第2个分页上
//第三个分页 2*6 6 =>12 6 代表从数据库里查询序号为12到17 这6条显示到第3个分页上
//...
func QueryArticlesWithCon(sql string) ([]Article, error) {
	sql = "select id,title,tags,short,content,author,createtime from article " + sql
	rows, err := utils.QueryDB(sql)
	if err != nil {
		return nil, err
	}
	var artList []Article
	for rows.Next() {
		id := 0
		title := ""
		tags := ""
		short := ""
		content := ""
		author := ""
		var createtime string
		createtime = ""
		rows.Scan(&id, &title, &tags, &short, &content, &author, &createtime)
		art := Article{id, title, tags, short, content, author, createtime}
		artList = append(artList, art)
	}
	return artList, nil
}

//---------数据处理-----------
//func AddArticle(article Article) (int64, error) {
//	i, err := insertArticle(article)
//	return i, err
//}

//-----------数据库操作---------------

//插入一篇文章
func insertArticle(article Article) (int64, error) {
	return utils.ModifyDB("insert into article(title,tags,short,content,author,createtime) values(?,?,?,?,?,?)",
		article.Title, article.Tags, article.Short, article.Content, article.Author, article.Createtime)
}
//------翻页------

//存储表的行数，只有自己可以更改，当文章新增或者删除时需要更新这个值
var artcileRowsNum = 0

//只有首次获取行数的时候采取统计表里的行数
func GetArticleRowsNum() int {
	if artcileRowsNum == 0 {
		artcileRowsNum = QueryArticleRowNum()
	}
	return artcileRowsNum
}

//查询文章的总条数
func QueryArticleRowNum() int {
	row := utils.QueryRowDB("select count(id) from article")
	num := 0
	row.Scan(&num)
	return num
}
//设置页数
func SetArticleRowsNum(){
	artcileRowsNum = QueryArticleRowNum()
}
//---------添加文章-----------
func AddArticle(article Article) (int64, error) {
	i, err := insertArticle(article)
	SetArticleRowsNum()
	return i, err
}
//----------删除文章---------
func DeleteArticle(artID int) (int64, error) {
	row := utils.QueryRowDB("select content from article where id="+ strconv.Itoa(artID))
	url := ""
	row.Scan(&url)

	//fmt.Println("删除文件的路径:",url)
	rwLock.Lock()//写锁
	err := os.RemoveAll(url)
	if err != nil{
		fmt.Println("删除文件出错:",err)
	}
	i, err := deleteArticleWithArtId(artID)
	SetArticleRowsNum()
	rwLock.Unlock()
	return i, err
}
func deleteArticleWithArtId(artID int) (int64, error) {
	return utils.ModifyDB("delete from article where id=?", artID)
}
//----------查询文章-------------

func QueryArticleWithId(id int,) Article {
	row := utils.QueryRowDB("select id,title,tags,short,content,author,createtime from article where id=" + strconv.Itoa(id))
	title := ""
	tags := ""
	short := ""
	content := ""
	author := ""
	var createtime string
	createtime = ""
	row.Scan(&id, &title, &tags, &short, &content, &author, &createtime)
	art := Article{id, title, tags, short, content, author, createtime}
	return art
}
//----------修改数据----------

func UpdateArticle(article Article,content string) (int64, error) {
	//数据库操作
	row := utils.QueryRowDB("select content from article where id="+ strconv.Itoa(article.Id))
	url := ""
	row.Scan(&url)
	//file,err:=os.OpenFile(url,
	//	os.O_CREATE|os.O_WRONLY|os.O_APPEND,os.ModePerm)
	//if err != nil {
	//	fmt.Println("打开文件有误：", err.Error())
	//	return 0,err
	//}
	//_, err = file.WriteString(content)
	//if err != nil{
	//	fmt.Println("写入内容有误",err)
	//	return 0,err
	//}
	fmt.Println("文件保存路径:",url)
	//_, err :=copyFile(url,url)
	//if err != nil{
	//	fmt.Println("更新覆盖文件出错:",err)
	//}
	err := os.RemoveAll(url)
	if err != nil{
		fmt.Println("删除文件出错:",err)
	}
	filename := "static/file/"+ article.Title +".md"
	file,err:=os.OpenFile(filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,os.ModePerm)
	if err != nil {
		fmt.Println("打开新文件有误:", err.Error())
		return 0,err
	}
	_, err = file.WriteString(content)
	if err != nil{
		fmt.Println("写入新内容有误:",err)
		return 0,err
	}

	return utils.ModifyDB("update article set title=?,tags=?,short=?,content=? where id=?",
		article.Title, article.Tags, article.Short,filename, article.Id)
}
//func copyFile(srcFile, destFile string)(int64,error){
//	file1,err:=os.Open(srcFile)
//	if err != nil{
//		return 0,err
//	}
//	file2,err:=os.OpenFile(destFile,os.O_WRONLY|os.O_CREATE,os.ModePerm)
//	if err !=nil{
//		return 0,err
//	}
//	defer file1.Close()
//	defer file2.Close()
//
//	return io.Copy(file2,file1)
//}