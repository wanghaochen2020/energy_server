package errmsg

// TODO 修改错误类型
const (
	SUCCESS = 200
	ERROR   = 500

	// code= 10x... 用户模块错误
	ERROR_PASSWORD_WRONG   = 101
	ERROR_USER_NOT_EXIST   = 102
	ERROR_TOKEN_EXIST      = 103
	ERROR_TOKEN_RUNTIME    = 104
	ERROR_TOKEN_WRONG      = 105
	ERROR_TOKEN_TYPE_WRONG = 106
)

var codeMsg = map[int]string{
	SUCCESS:                "OK",
	ERROR:                  "FAIL",
	ERROR_PASSWORD_WRONG:   "密码错误",
	ERROR_USER_NOT_EXIST:   "用户不存在",
	ERROR_TOKEN_EXIST:      "TOKEN不存在,请重新登陆",
	ERROR_TOKEN_RUNTIME:    "TOKEN已过期,请重新登陆",
	ERROR_TOKEN_WRONG:      "TOKEN不正确,请重新登陆",
	ERROR_TOKEN_TYPE_WRONG: "TOKEN格式错误,请重新登陆",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
