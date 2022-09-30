package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type _unit struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	Name    string `json:"name"`
	Time    string `json:"time"`
	User    string `json:"user"`
}
type _user struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Uid      int64  `json:"uid"`
}
type _comment struct {
	Id      int64  `json:"id"`
	Comment string `json:"comment"`
	Uid     string `json:"uid"`
}

var isLogin int
var postNumbers int64    //post数量
var userNumbers int64    //用户数量
var commentNumbers int64 //评论数量
var del []int = []int{0} //记录删除post的id
func main() {
	InitOpen()                                        //初始化数据库
	http.HandleFunc("/save", dataHandler)             //储存推文
	http.HandleFunc("/show", showHandler)             //随机显示一条推文
	http.HandleFunc("/getPosts", GetAllHandler)       //获取全部推文
	http.HandleFunc("/addComment", addCommentHandler) //添加评论
	http.HandleFunc("/comment", commentHandler)       //显示评论
	http.HandleFunc("/change", changeHandler)         //改变推文
	http.HandleFunc("/delete", deleteHandler)         //删除推文
	http.HandleFunc("/register", regHandler)          //用户注册
	http.HandleFunc("/login", loginHandler)           //用户登录
	http.HandleFunc("/admin", adminHandler)           //管理员登录
	http.ListenAndServe("localhost:8080", nil)        //阻塞监听

}

func dataHandler(writer http.ResponseWriter, request *http.Request) {
	var unit _unit
	request.ParseForm() //解析表单
	method := request.Method
	if method == "POST" {
		//{此处方法为读取表单数据
		//	unit.content = request.Form.Get("content")
		//	unit.name = request.Form.Get("name")
		unit.Time = time.Now().Format("2006/1/02/ 15:04")
		//	unit.user = request.Form.Get("user")
		//}
		data, err := ioutil.ReadAll(request.Body)
		checkErr(err)
		json.Unmarshal(data, &unit) //解码json

		//存入数据库
		db, err := sql.Open("sqlite3", "wall.db")
		checkErr(err)
		insert, err := db.Prepare("INSERT INTO content(content,name,time,user) values (?,?,?,?)")
		checkErr(err)
		res, err := insert.Exec(unit.Content, unit.Name, unit.Time, unit.User)
		checkErr(err)
		postNumbers, err = res.LastInsertId()
		checkErr(err)
		writer.WriteHeader(201)
		db.Close()

	}

}
func showHandler(writer http.ResponseWriter, request *http.Request) {
	var id int
	var re bool = true
	rand.Seed(time.Now().UnixNano())
	id = rand.Intn(int(postNumbers+1)) + 1

	//避免随机到删除post的id
	for i := 0; re == true; i++ {
		re = false
		for j := 0; j < len(del); j++ {
			if id == del[j] {
				re = true
			}
		}

		if re == true {
			rand.Seed(time.Now().UnixNano())
			id = rand.Intn(int(postNumbers))
		}
	}
	var unit _unit
	db, err := sql.Open("sqlite3", "wall.db")
	checkErr(err)
	res, err := db.Query("SELECT * FROM content")
	for res.Next() {
		res.Scan(&unit.Id, &unit.Content, &unit.Name, &unit.Time, &unit.User)
		_id, _ := strconv.Atoi(unit.Id)
		if _id == id {
			data, err := json.Marshal(unit)
			checkErr(err)
			writer.Header().Set("Content-Type", "application/json") //设置响应头数据类型为json类型
			writer.Write(data)
			break
		}
	}
	checkErr(err)
	db.Close()
}
func GetAllHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	username := request.Form.Get("username")
	method := request.Method
	if method == "POST" {
		var units []_unit
		var unit _unit
		db, err := sql.Open("sqlite3", "wall.db")
		checkErr(err)
		rows, err := db.Query("SELECT * FROM content")
		for rows.Next() {
			rows.Scan(&unit)
			if unit.User == username {
				units = append(units, unit)
			}
		}
		writer.Header().Set("Content-Type", "application/json") //设置响应头数据类型为json类型
		data, err := json.Marshal(units)
		writer.Write(data)
		db.Close()
	}
}
func addCommentHandler(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	var comment _comment
	if method == "POST" {
		data, _ := ioutil.ReadAll(request.Body)
		json.Unmarshal(data, &comment)
		db, err := sql.Open("sqlite3", "wall.db")
		checkErr(err)
		insert, err := db.Prepare("INSERT INTO comments(comment, uid) VALUES (?,?)")
		checkErr(err)
		res, err := insert.Exec(comment.Comment, comment.Uid)
		commentNumbers, err = res.LastInsertId()
		checkErr(err)
		writer.WriteHeader(201) //已创建
		db.Close()
	}
}
func commentHandler(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method == "GET" {
		var comments []_comment
		var comment _comment
		request.ParseForm()
		id := request.Form.Get("id") //所查看评论的推文id
		//var Uid string
		//data, err := ioutil.ReadAll(request.Body)
		//checkErr(err)
		//json.Unmarshal(data, &Uid)
		//fmt.Println(Uid)
		_id, _ := strconv.Atoi(id)
		db, err := sql.Open("sqlite3", "wall.db")
		checkErr(err)
		res, err := db.Query("SELECT * FROM comments")
		for res.Next() {
			res.Scan(&comment.Id, &comment.Comment, &comment.Uid)

			_uid, _ := strconv.Atoi(comment.Uid)
			if _uid == _id {
				comments = append(comments, comment)
			}
		}
		writer.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(comments)
		writer.Write(data)
	}
}
func changeHandler(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method == "POST" {
		//request.ParseForm()
		//newContent := request.Form.Get("newContent")
		//id := request.Form.Get("id")
		var unit _unit
		data, err := ioutil.ReadAll(request.Body)
		checkErr(err)
		json.Unmarshal(data, &unit)
		db, err := sql.Open("sqlite3", "wall.db")
		checkErr(err)
		change, err := db.Prepare("update content set content=? where id=?")
		checkErr(err)
		res, err := change.Exec(unit.Content, unit.Id)
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
	user := request.Form.Get("user")
	_id, _ := strconv.Atoi(id)
	del = append(del, _id)
	db, err := sql.Open("sqlite3", "wall.db")
	rows, err := db.Query("SELECT user FROM content")
	checkErr(err)
	var unit _unit
	for rows.Next() {
		rows.Scan(&unit)
		if unit.Id == id {
			if user == unit.User || user == "admin" {
				delete, err := db.Prepare("delete from content where id=?")
				checkErr(err)
				res, err := delete.Exec(id)
				checkErr(err)
				affect, err := res.RowsAffected()
				checkErr(err)
				fmt.Println(affect)
				postNumbers--
			}
		} else {
			writer.WriteHeader(401)
		}
	}
	//return
	db.Close()
}
func regHandler(writer http.ResponseWriter, request *http.Request) {
	var re int
	//var user _user
	user := _user{}

	request.ParseForm() //解析表单
	method := request.Method
	if method == "OPTIONS" {
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		writer.WriteHeader(200)

	}
	if method == "POST" {
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		//writer.Header().Add("Accept-Control-Allow-Origin", "*")
		//writer.Header().Set("Accept-Control-Allow-Origin", "*")
		//{
		//	user.username = request.Form.Get("username")
		//	user.password = request.Form.Get("password")
		//}
		data, err := ioutil.ReadAll(request.Body)
		json.Unmarshal(data, &user)
		{
			fmt.Println(user)
		}
		var userFromDB string
		db, err := sql.Open("sqlite3", "wall.db")
		checkErr(err)
		//检测是否名称重复
		rows, _ := db.Query("SELECT username FROM users")
		for rows.Next() {
			rows.Scan(&userFromDB)
			if user.Username == userFromDB {
				re = 1
				writer.WriteHeader(205) //名称重复，请求重置表单
				break
			}
		}
		if re == 0 { //名称不重复
			//存入数据库
			insert, err := db.Prepare("INSERT INTO users(username,password) values (?,?)")
			checkErr(err)
			res, err := insert.Exec(user.Username, user.Password)
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
	if method == "OPTIONS" {
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		writer.WriteHeader(200)

	}
	if method == "POST" {
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		//request.ParseForm()
		//{
		//	user.username = request.Form.Get("username")
		//	user.password = request.Form.Get("password")
		//}
		data, err := ioutil.ReadAll(request.Body)
		checkErr(err)
		json.Unmarshal(data, &user) //解码json
		//检测username&password
		//if user.Username == "admin" && user.Password == "admin" {
		//	//goto admin
		//	http.Redirect(writer, request, "localhost:8080/admin.html?username=admin&password=admin", http.StatusFound)
		//}
		db, _ := sql.Open("sqlite3", "wall.db")
		rows, _ := db.Query("SELECT * FROM users")
		success := 0
		for rows.Next() {
			rows.Scan(&userFromDB.Uid, &userFromDB.Username, &userFromDB.Password)
			if user.Username == userFromDB.Username && user.Password == userFromDB.Password {
				//http.Redirect(writer, request, "localhost:8080/index.com", http.StatusFound) //跳转到主页面
				success = 1
				//username, err := json.Marshal(user.Username)            //转换为json格式
				checkErr(err)
				isLogin = 1
				{
					fmt.Println("ssss")
				}
				writer.Write([]byte(user.Username)) //返回用户名
				break
			}
		}
		if success == 0 {
			writer.WriteHeader(511)
		}
		db.Close()
	}

}
func adminHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm() //解析表单
	adminUser := request.Form.Get("username")
	adminPassword := request.Form.Get("password")
	if adminUser == "admin" && adminPassword == "admin" {
		var units []_unit
		var unit _unit
		db, err := sql.Open("sqlite3", "wall.db")
		checkErr(err)
		rows, _ := db.Query("SELECT * FROM users")
		for rows.Next() {
			rows.Scan(&unit)
			units = append(units, unit)
		}
		var buffer bytes.Buffer
		num, _ := json.Marshal(userNumbers)
		u, _ := json.Marshal(units)
		buffer.Write(num)
		buffer.Write(u)
		data := buffer.Bytes()
		writer.Header().Set("Content-Type", "application/json") //设置响应头数据类型为json类型
		writer.Write(data)
	} else { //密码错误重定向至登录界面
		http.Redirect(writer, request, "localhost:8080/login.html", http.StatusFound)
	}
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
	sqlTableComment := `
	CREATE TABLE IF NOT EXISTS "comments"(
	    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
	    "comment" VARCHAR(100) NULL,
	    "uid" INTEGER NULL
	)`
	db.Exec(sqlTableComment)
	db.Close()
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
