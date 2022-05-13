package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func Cors() gin.HandlerFunc {
	return cors.New(
		cors.Config{
			//AllowCredentials: true, 												// 是否允许cookie
			AllowOrigins:     []string{"*"},                                       // 允许所有的跨域
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的请求方法
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"Content-Length", "text/plain", "Authorization", "Content-Type", "application/json"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour, // 域请求的持续时间，如果通过，12小时内不再需要域请求
		},
	)
}
