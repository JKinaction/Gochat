package routes

import (
	"GoChat/handler"
	"GoChat/ws"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func IntializeRoutes() *gin.Engine {
	r := gin.Default()
	hub := ws.NewHub()
	go hub.Run()
	r.LoadHTMLGlob("static/template/*")
	r.Use(sessions.Sessions("chatsession", sessions.NewCookieStore([]byte("secret"))))
	r.Static("/css", "static/css")
	r.GET("/", handler.LoadLogin)
	r.POST("/auth", handler.Auth)
	r.GET("/ws", func(context *gin.Context) {
		ws.ServeWs(hub, context)
	})
	uRoutes := r.Group("/u")
	uRoutes.Use(handler.IsLogin)
	{
		uRoutes.Static("/js", "static/js")
		uRoutes.Static("/css", "static/css")
		uRoutes.GET("/chat", handler.Chatpage)
		uRoutes.GET("/logout", handler.Logout)
		uRoutes.GET("privatechat", handler.PrivateChat)
	}
	return r
}
