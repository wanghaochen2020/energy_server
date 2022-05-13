package middleware

import (
	"energy/utils"
	"energy/utils/errmsg"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var JwtKey = []byte(utils.JwtKey)
var code int

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成token
func SetToken(username string) (string, int) {
	expireTime := time.Now().Add(10 * time.Hour) // token有效时间10小时
	SetClaims := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // token过期时间
			Issuer:    "huanwei",         // 指定token发行人
		},
	}

	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims) // 第一个为签发的方法，即加密函数,第二个为定义的结构体
	// 该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := reqClaim.SignedString(JwtKey)
	if err != nil {
		return "", errmsg.ERROR
	}
	return token, errmsg.SUCCESS
}

// 验证token
func CheckToken(token string) (*MyClaims, int) {
	// 用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	setToken, _ := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	// 从setToken中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的MyClaims
	if key, ok := setToken.Claims.(*MyClaims); ok && setToken.Valid { // 要传入指针，项目中结构体都是用指针传递，节省空间
		return key, errmsg.SUCCESS
	} else {
		return nil, errmsg.ERROR
	}
}

// jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization") // 前端的Authorization存放token,response.addHeader("Authorization", "Bearer " + jwt);
		if tokenHeader == "" {                               // 没有Authorization说明token不存在
			code = errmsg.ERROR_TOKEN_EXIST
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		checkToken := strings.SplitN(tokenHeader, " ", 2)      // 如果有Authorization字段则为"Bearer " + jwt，中间有一个空格，所以按照" "划分字段
		if len(checkToken) != 2 && checkToken[0] != "Bearer" { // token格式错误
			code = errmsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		key, tCode := CheckToken(checkToken[1]) // key是一个MyClaims结构体类型变量
		if tCode == errmsg.ERROR {              // token解析错误
			code = errmsg.ERROR_TOKEN_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		if time.Now().Unix() > key.ExpiresAt { // token过期
			code = errmsg.ERROR_TOKEN_RUNTIME
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		c.Set("username", key.Username)
		c.Next()
	}
}
