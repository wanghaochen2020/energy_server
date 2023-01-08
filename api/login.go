package api

import (
	"energy/middleware"
	"energy/model"
	"energy/utils/errmsg"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 登录,返回一个token
func Login(c *gin.Context) {
	var data model.User
	_ = c.ShouldBindJSON(&data)
	var code int
	var token string // 返回给用户的toke

	code = model.CheckLogin(data.Username, data.Password)
	if code == errmsg.SUCCESS {
		token, code = middleware.SetToken(data.Username)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   errmsg.GetErrMsg(code),
		"token": token,
	})
}

// 只验证token
func Verify(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  errmsg.GetErrMsg(200),
	})
}

// 添加用户
func AddUser(c *gin.Context) {
	var data model.User
	_ = c.ShouldBindJSON(&data)
	model.CreateUser(&data)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  errmsg.GetErrMsg(200),
	})
}
