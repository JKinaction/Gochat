package handler

import (
	"GoChat/model"
	session2 "GoChat/session"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strings"
)

const (
	username = "user"
)

func LoadLogin(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(username)
	if user != nil {
		c.Redirect(http.StatusFound, "/u/chat")
	}
	c.HTML(
		http.StatusOK,
		"login.html",
		gin.H{},
	)
}

func Auth(c *gin.Context) {

	user := c.PostForm("user")
	pwd := c.PostForm("pwd")
	fmt.Println(user, pwd)
	if strings.Contains(user, "<") || strings.Contains(user, ">") || user == "" {
		c.HTML(
			http.StatusOK,
			"login.html",
			gin.H{
				"error": "Use a valid username",
			},
		)
		return
	}
	u, err := model.GetPwd(user)
	fmt.Println(u)
	if err != nil {
		u, err = model.AddUsertoDB(user, pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if u.Pwd != pwd {
		c.HTML(
			http.StatusOK,
			"login.html",
			gin.H{
				"error": "please input correct password",
			},
		)
		return
	}

	err = session2.SaveSession(c, username, user)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": "Unable to save session."})
		return
	} else {
		c.Redirect(http.StatusFound, "/u/chat")
	}
}

func IsLogin(c *gin.Context) {
	user := session2.GetSession(c, username)
	if user == nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.Next()
}

func Chatpage(c *gin.Context) {

	user := session2.GetSession(c, username)
	//SELECT * FROM messages获得数据库中所有消息
	msglist, err := model.GetAllMsgsDB()

	if err != nil {
		log.Fatalln("Unable to retrieve messages from database.")
	}
	//将所有消息拼接到globalmsgList
	msglist = append(model.GlobalmsgList, msglist...)
	users, err := model.GetAllUsers()
	if err != nil {
		log.Fatalln(err)
	}

	c.HTML(
		http.StatusOK,
		"chat.html",
		gin.H{
			"user":     user,
			"msgList":  msglist,
			"userlist": users,
		},
	)
}
func PrivateChat(c *gin.Context) {
	user := session2.GetSession(c, username)

	if c.Query("touser") != "" {
		touser := c.Query("touser")
		err := session2.SaveSession(c, "touser", touser)
		if err != nil {
			log.Fatalln(err)
		}
	}
	touser := session2.GetSession(c, "touser")
	msglist, err := model.GetPriMsgsDB(touser.(string), user.(string))
	if err != nil {
		log.Fatalln(msglist)
	}
	c.HTML(
		http.StatusOK,
		"whisper.html",
		gin.H{
			"user":    user,
			"touser":  touser,
			"msgList": msglist,
		},
	)
}
func Logout(c *gin.Context) {
	session2.ClearSession(c)
	c.Redirect(302, "/")
}
