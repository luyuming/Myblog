package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"myblog/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LoginInfo struct {
	Islogin  bool
	Username string
	Userid   uint
}

func getLogInfo(info *LoginInfo, c *gin.Context) {
	name, err := c.Cookie("name")
	if err != nil {
		info.Islogin = false
	} else {
		info.Islogin = true
		info.Username = name
		u, _ := models.GetUser(name)
		info.Userid = u.ID
	}
}

func Index(c *gin.Context) {
	var loginfo LoginInfo
	getLogInfo(&loginfo, c)
	var recordusers []models.RecordUser
	var ShowRecords []models.Record
	var recordnum int
	var pageNow int
	if strings.Contains(c.Request.RequestURI, "page") {
		pageNow, _ = strconv.Atoi(c.Param("pageId"))
	} else {
		pageNow = 1
	}
	models.GetCount("records", &recordnum)
	pager := CreatePaginator(pageNow, defaultPageSize, recordnum)
	offset := defaultPageSize * (pageNow - 1)
	models.GetLimitRecordsOffset(defaultPageSize, offset, &ShowRecords)
	for _, record := range ShowRecords {
		var u models.User
		models.RecordRelatedUser(&record, &u)
		key := "record" + strconv.Itoa(int(record.ID)) + "_likes"
		val, _ := models.GetKV(key)
		var temp models.RecordUser
		temp.Record = record
		temp.UserName = u.Name
		temp.Star, _ = strconv.Atoi(val)
		recordusers = append(recordusers, temp)
	}
	c.HTML(http.StatusOK, "index.html", gin.H{"records": recordusers, "logininfo": loginfo, "paginator": pager})
}

func Register(c *gin.Context) {
	if c.Request.Method == "POST" {
		name := c.PostForm("username")
		//判断是否已注册
		u, _ := models.GetUser(name)
		//未注册,将用户信息存入数据库
		if u == nil || u.ID == 0 {
			var u1 models.User
			c.Bind(&u1)
			models.AddUser(&u1)
			msg, _ := json.Marshal(u1)
			ProductMsg(msg)
			time.Sleep(1 * time.Second)
			c.Redirect(http.StatusMovedPermanently, "/login")
		} else { //已注册,提示该用户名已注册，
			c.String(http.StatusOK, "该用户名已注册")
		}
	} else {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	}
}

func Login(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	} else {
		name := c.PostForm("username")
		password := c.PostForm("password")
		u, _ := models.GetUser(name)
		if u == nil || u.ID == 0 { //该用户不存在
			c.String(http.StatusOK, "该用户不存在")
		} else { //该用户存在
			if password == u.Password { //密码正确,重定向到首页
				c.SetCookie("name", name, 3600, "/", "127.0.0.1", false, true)
				c.Redirect(http.StatusMovedPermanently, "/")
				//使用rabbitmq消息队列，把用户邮箱名发送到消息队列，通过rabbitmq给用户发送邮件
				//product(name)
			} else { //密码错误
				c.String(http.StatusOK, "密码错误")
			}
		}
	}
}

func Quit(c *gin.Context) {
	name, err := c.Cookie("name")
	if err == nil {
		c.SetCookie("name", name, -1, "/", "127.0.0.1", false, true)
	}
	c.Redirect(http.StatusFound, "/")
}

func PersonalIndex(c *gin.Context) {
	var loginfo LoginInfo
	getLogInfo(&loginfo, c)
	id := c.Param("userId")
	u, _ := models.GetUserRecord(id)
	var visits int64
	key := "user" + id + "_visits"
	hasVisit, _ := models.CheckKeyExist(key)
	if hasVisit == 0 {
		visits = 1
		models.SetKV(key, 0, 0)
	} else {
		visits, _ = models.IncrKey(key)
	}
	c.HTML(http.StatusOK, "user.html", gin.H{"userinfo": u, "logininfo": loginfo, "visits": visits})
}

func RecordIndex(c *gin.Context) {
	var loginfo LoginInfo
	getLogInfo(&loginfo, c)
	id := c.Param("recordId")
	record, _ := models.GetRecordComment(id)
	var user models.User
	models.RecordRelatedUser(record, &user)
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "record.html", gin.H{"userinfo": user, "recordinfo": record, "logininfo": loginfo})
	} else {
		if loginfo.Islogin == false {
			c.String(http.StatusOK, "请先登录后再评论")
		} else {
			comment := new(models.Comment)
			comment.Content = c.PostForm("comment")
			comment.Name = user.Name
			models.RecordAssocComment(record, comment)
			c.Redirect(http.StatusMovedPermanently, c.Request.RequestURI)
		}
	}
}

func RecordCreate(c *gin.Context) {
	var loginfo LoginInfo
	getLogInfo(&loginfo, c)
	name, _ := c.Cookie("name")
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "create.html", gin.H{"name": name, "logininfo": loginfo})
	} else {
		var article models.Record
		c.Bind(&article)
		models.AddRecord(&article)
		u, _ := models.GetUser(name)
		models.UserAssocRecord(u, &article)
		location := fmt.Sprintf("%s%d", "/record/read/", article.ID)
		c.Redirect(http.StatusMovedPermanently, location)
	}
}

func RecordLikes(c *gin.Context) {
	islike := c.Query("islike")
	var val int64
	var err error
	key := "record" + c.Param("id") + "_likes"
	if islike == "true" {
		val, err = models.IncrKey(key)
	} else {
		val, err = models.DecrKey(key)
	}
	if err != nil {
		panic(err)
	}
	c.String(http.StatusOK, strconv.Itoa(int(val)))
}
