package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	request.ParseForm()
	method := request.Method
	if method == "POST" {
		postNumbers++
		unit.Time = time.Now().Format("2006/1/02/ 15:04")
		data, err := ioutil.ReadAll(request.Body)
		checkErr(err)
		json.Unmarshal(data, &unit) //解码json
		unit.Id = strconv.Itoa(int(postNumbers))
		//存入数据库
		db, err := gorm.Open("sqlite3", "wall.db")
		checkErr(err)
		defer db.Close()
		db.Create(&unit)
		writer.WriteHeader(201)
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
	db, err := gorm.Open("sqlite3", "wall.db")
	checkErr(err)
	defer db.Close()
	db.Take(&unit)
	data, err := json.Marshal(unit)
	checkErr(err)
	writer.Header().Set("Content-Type", "application/json") //设置响应头数据类型为json类型
	writer.Write(data)
}
func GetAllHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	username := request.Form.Get("user")
	method := request.Method
	if method == "POST" {
		var units []_unit
		db, err := gorm.Open("sqlite3", "wall.db")
		checkErr(err)
		defer db.Close()
		db.Where("user=?", username).Find(&units)
		writer.Header().Set("Content-Type", "application/json") //设置响应头数据类型为json类型
		data, err := json.Marshal(units)
		writer.Write(data)
	}
}
func addCommentHandler(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	var comment _comment
	if method == "POST" {
		commentNumbers++
		data, _ := ioutil.ReadAll(request.Body)
		json.Unmarshal(data, &comment)
		comment.Id = commentNumbers
		db, err := gorm.Open("sqlite3", "wall.db")
		checkErr(err)
		defer db.Close()
		db.Create(&comment)
		writer.WriteHeader(201) //已创建
	}
}
func commentHandler(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method == "GET" {
		var comments []_comment
		request.ParseForm()
		id := request.Form.Get("id") //所查看评论的推文id
		_id, _ := strconv.Atoi(id)
		db, err := gorm.Open("sqlite3", "wall.db")
		checkErr(err)
		defer db.Close()
		db.Where("uid=?", _id).Find(&comments)
		writer.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(comments)
		writer.Write(data)
	}
}
func changeHandler(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method == "POST" {
		var unit _unit
		data, err := ioutil.ReadAll(request.Body)
		checkErr(err)
		json.Unmarshal(data, &unit)
		db, err := gorm.Open("sqlite3", "wall.db")
		checkErr(err)
		defer db.Close()
		db.Model(&unit).Where("id=?", unit.Id).Update("content", unit.Content)
	}
}
func deleteHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	id := request.Form.Get("id")
	_id, _ := strconv.Atoi(id)
	del = append(del, _id)
	db, err := gorm.Open("sqlite3", "wall.db")
	checkErr(err)
	defer db.Close()
	db.Where("id=?", id).Delete(_unit{})
	writer.WriteHeader(410)
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
		data, err := ioutil.ReadAll(request.Body)
		userNumbers++
		json.Unmarshal(data, &user)
		user.Uid = userNumbers
		{
			fmt.Println(user)
		}
		var userFromDB _user
		db, err := gorm.Open("sqlite3", "wall.db")
		checkErr(err)
		defer db.Close()
		//检测是否名称重复
		db.Where("username=?", user.Username).First(&userFromDB)
		if userFromDB.Username != "" {
			re = 1
			writer.WriteHeader(205)
		}
		if re == 0 { //名称不重复
			//存入数据库
			db.Create(&user)
			writer.WriteHeader(201) //已创建
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
		data, err := ioutil.ReadAll(request.Body)
		checkErr(err)
		json.Unmarshal(data, &user) //解码json
		user.Uid = userNumbers
		db, err := gorm.Open("sqlite3", "wall.db")
		checkErr(err)
		defer db.Close()
		success := 0
		db.Where("username=?", user.Username).First(&userFromDB)
		if user.Username == userFromDB.Username && user.Password == userFromDB.Password {
			success = 1
			isLogin = 1
			writer.Write([]byte(user.Username)) //返回用户名
		}
		if success == 0 {
			writer.WriteHeader(511)
		}
	}
}
func adminHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm() //解析表单
	adminUser := request.Form.Get("username")
	adminPassword := request.Form.Get("password")
	if adminUser == "admin" && adminPassword == "admin" {
		var units []_unit
		db, err := gorm.Open("sqlite3", "wall.db")
		checkErr(err)
		defer db.Close()
		db.Find(&units)
		var buffer bytes.Buffer
		num, _ := json.Marshal(userNumbers)
		u, _ := json.Marshal(units)
		buffer.Write(num)
		buffer.Write(u)
		data := buffer.Bytes()
		writer.Header().Set("Content-Type", "application/json") //设置响应头数据类型为json类型
		writer.Write(data)
	} else { //密码错误
		http.Redirect(writer, request, "localhost:8080/login.html", http.StatusFound)
	}
}
func InitOpen() {

	db, err := gorm.Open("sqlite3", "wall.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&_comment{})
	db.AutoMigrate(&_unit{})
	db.AutoMigrate(&_user{})
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
