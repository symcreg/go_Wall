package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type _unit struct {
	uid     int64
	content string
	name    string
	time    string
	user    string
}
type _user struct {
	uid      int64
	username string
	password string
}

var postNumbers int64 //post数量
var userNumbers int64 //用户数量
func main() {
	InitOpen() //初始化数据库
	http.HandleFunc("/save", dataHandler)
	http.HandleFunc("/change", changeHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/register", regHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/admin", adminHandler)
	http.ListenAndServe("localhost:8080", nil) //阻塞监听
}

func dataHandler(writer http.ResponseWriter, request *http.Request) {
	var unit _unit
	request.ParseForm() //解析表单
	method := request.Method
	if method == "POST" {
		unit.content = request.Form.Get("content")
		unit.name = request.Form.Get("name")
		unit.time = time.Now().Format("2006/1/02/ 15:04")
		unit.user = request.Form.Get("user")
		//存入数据库
		db, err := sql.Open("sqlite3", "wall.db")
		checkErr(err)
		insert, err := db.Prepare("INSERT INTO content(content,name,time,user) values (?,?,?,?)")
		checkErr(err)
		res, err := insert.Exec(unit.content, unit.name, unit.time, unit.user)
		checkErr(err)
		postNumbers, err = res.LastInsertId()
		checkErr(err)

		db.Close()

		{
			body, _ := ioutil.ReadAll(request.Body)
			var s map[string]interface{}
			json.Unmarshal(body, &s) //读取json
			msg, _ := json.Marshal(unit)
			io.WriteString(writer, string(msg)) ////////////////////////返回数据
		}

	}

}
func changeHandler(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method == "POST" {
		request.ParseForm()
		newContent := request.Form.Get("newContent")
		id := request.Form.Get("id")
		db, err := sql.Open("sqlite3", "wall.db")
		checkErr(err)
		change, err := db.Prepare("update content set content=? where id=?")
		checkErr(err)
		res, err := change.Exec(newContent, id)
		checkErr(err)
		affect, err := res.RowsAffected()
		checkErr(err)
		fmt.Println(affect)
		//return
		db.Close()
	}
}
func deleteHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	id := request.Form.Get("id")
	db, err := sql.Open("sqlite3", "wall.db")
	checkErr(err)
	delete, err := db.Prepare("delete from content where id=?")
	checkErr(err)
	res, err := delete.Exec(id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)
	postNumbers--
	//return
	db.Close()
}
func regHandler(writer http.ResponseWriter, request *http.Request) {
	var re int
	var user _user
	request.ParseForm() //解析表单
	method := request.Method
	if method == "POST" {
		//{
		//	user.username = request.Form.Get("username")
		//	user.password = request.Form.Get("password")
		//}
		data, err := ioutil.ReadAll(request.Body)
		checkErr(err)
		json.Unmarshal(data, &user) //解码json

		var userFromDB string
		db, err := sql.Open("sqlite3", "wall.db")
		checkErr(err)
		//检测是否名称重复
		rows, _ := db.Query("SELECT username FROM users")
		for rows.Next() {
			rows.Scan(&userFromDB)
			if user.username == userFromDB {
				re = 1
				writer.WriteHeader(205) //名称重复，请求重置表单
				break
			}
		}
		if re == 0 { //名称不重复
			//存入数据库
			insert, err := db.Prepare("INSERT INTO users(username,password) values (?,?)")
			checkErr(err)
			res, err := insert.Exec(user.username, user.password)
			checkErr(err)
			userNumbers, err = res.LastInsertId()
			checkErr(err)
			writer.WriteHeader(201) //已创建
			//goto login
			db.Close()
		}
	}

}
func loginHandler(writer http.ResponseWriter, request *http.Request) {
	var user _user
	var userFromDB _user
	method := request.Method

	if method == "POST" {
		request.ParseForm()
		//{
		//	user.username = request.Form.Get("username")
		//	user.password = request.Form.Get("password")
		//}
		data, err := ioutil.ReadAll(request.Body)
		checkErr(err)
		json.Unmarshal(data, &user) //解码json
		//检测username&password
		db, _ := sql.Open("sqlite2", "wall.db")
		rows, _ := db.Query("SELECT * FROM users")
		success := 0
		for rows.Next() {
			rows.Scan(&userFromDB.uid, &userFromDB.username, &userFromDB.password)
			if user.username == userFromDB.username && user.password == userFromDB.password {
				//goto index
				//http.Redirect(writer, request, "localhost:8080/index.com", http.StatusFound) //跳转到主页面
				success = 1
				writer.Header().Set("Content-Type", "application/json") //设置响应头数据类型为json类型
				username, err := json.Marshal(user.username)            //转换为json格式
				checkErr(err)
				writer.Write([]byte(username)) //返回用户名

				break
			}
		}
		if success == 0 {
			writer.WriteHeader(511)
		}
	}
	if user.username == "admin" && user.password == "admin" {
		//goto admin

	}
}
func adminHandler(writer http.ResponseWriter, request *http.Request) {
	//
}
func InitOpen() {
	db, err := sql.Open("sqlite3", "wall.db")
	if err != nil {
		panic(err)
	}
	sqlTableContent := `CREATE TABLE IF NOT EXISTS "content"(
	    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
	    "content" VARCHAR(1024) NULL,
	    "name" VARCHAR(20) NULL,
	    "time" VARCHAR(50) NULL,
	    "user" VARCHAR(100) NULL
	)`
	db.Exec(sqlTableContent)
	sqlTableUser := `
	CREATE TABLE IF NOT EXISTS "users"(
	    "uid" INTEGER PRIMARY KEY AUTOINCREMENT,
	    "username" VARCHAR(100) NULL,
	    "password" VARCHAR(100) NULL
	)`
	db.Exec(sqlTableUser)

	db.Close()
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
