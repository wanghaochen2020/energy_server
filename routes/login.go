package routes

import (
	"energy/api"
	"github.com/gin-gonic/gin"
)

func LoginRouter(router *gin.RouterGroup) {
	// 用户登录模块路由接口
	router.POST("login", api.Login)
	// 新建用户路由接口
	router.POST("adduser", api.AddUser)
}
