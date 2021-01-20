package routers

import (
	"github.com/gin-gonic/gin"
	"myblog/controller"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	//r.Use(getCookie())
	r.StaticFS("/css", http.Dir("templates/css"))
	r.StaticFS("/js", http.Dir("templates/js"))
	r.StaticFS("/img", http.Dir("templates/img"))
	r.LoadHTMLGlob("templates/html/*")
	r.GET("/register", controller.Register)
	r.POST("/register", controller.Register)
	r.GET("/login", controller.Login)
	r.POST("/login", controller.Login)
	r.GET("/quit", controller.Quit)
	r.GET("/", controller.Index)
	r.GET("/user/:userId", controller.PersonalIndex)
	r.GET("/record/read/:recordId", controller.RecordIndex)
	r.POST("/record/read/:recordId", controller.RecordIndex)
	r.GET("/record/create", controller.RecordCreate)
	r.POST("/record/create", controller.RecordCreate)
	r.GET("/record/likes/:id", controller.RecordLikes)
	r.GET("/record/page/:pageId", controller.Index)
	return r
}
