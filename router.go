package main

import (
	"energy/api"
	"energy/middleware"
	"energy/routes"
	"energy/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() { //可以返回一个*gin.Engine
	gin.SetMode(utils.AppMode)
	r := gin.New()
	// 加载中间件
	r.Use(middleware.Cors())     // 跨域中间件
	r.Use(middleware.Logger())   // 日志中间件
	r.Use(gin.Recovery())        // 恢复恐慌中间件
	_ = r.SetTrustedProxies(nil) // 信任所有ip端口

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该路由",
		})
	})
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该方法",
		})
	})

	// rPublic组内的路由不需要token验证身份
	rPublic := r.Group("api")
	{
		// 登录路由
		routes.LoginRouter(rPublic)

		// 页面数据
		rPublic.GET("pageData", api.GetPageData)
	}
	// rAuth组内的路由需要有jwt的token
	rAuth := r.Group("api")
	rAuth.Use(middleware.JwtToken())
	{
		// 验证token
		rAuth.POST("userJwt", api.Verify)
		// TODO 后续所有接口都在这里，验证token

	}

	_ = r.Run(utils.HttpPort)
}
