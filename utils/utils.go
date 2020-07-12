package utils

import (
	"database/sql"
	"fmt"
	//"log"
	_ "github.com/Go-SQL-Driver/MySQL"
	//"github.com/jmoiron/sqlx"
)
func init(){
	InitMysql()
}

var db *sql.DB
var err error
func InitMysql() {
	fmt.Println("InitMysql....")
	if db == nil {
		db ,err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/myblogweb?charset=utf8")
		if err != nil{
			fmt.Println("设置连接数据库的参数有误...",err)
			//db.Close()
			return
		}
		fmt.Println("连接数据库成功...")
		CreateTableWithUser()
		CreateTableWithArticle()
	}
}
//创建用户表
func CreateTableWithUser() {
	sql := `CREATE TABLE IF NOT EXISTS users(
        id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
        username VARCHAR(64),
        password VARCHAR(64),
        status INT(4),
        createtime INT(10)
        );`

	ModifyDB(sql)
	//isSuccess,_:= ModifyDB(sql)
	//fmt.Println(isSuccess)
}
//创建文章表
func CreateTableWithArticle(){
	sql:=`create table if not exists article(
        id int(4) primary key auto_increment not null,
        title varchar(120),
        author varchar(20),
        tags varchar(30),
        short varchar(255),
        content longtext,
        createtime varchar(255)
        );`
	//createtime int(10)
	ModifyDB(sql)
}
//操作数据库
func ModifyDB(sql string, args ...interface{}) (int64, error) {
	result, err := db.Exec(sql, args...)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	//fmt.Println("count的值:",count)
	return count, nil
}
//查询
func QueryRowDB(sql string) *sql.Row{
	return db.QueryRow(sql)
}
//分页查询
func QueryDB(sql string)(*sql.Rows,error){
	rows,err := db.Query(sql)
	if err != nil {
		fmt.Println("查询有误...")
		return nil,err
	}
	return rows,nil
}