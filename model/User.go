package model

import (
	"encoding/base64"
	"energy/utils/errmsg"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null " json:"username"`
	Password string `gorm:"type:varchar(20);not null " json:"password"`
}

// 新增用户,返回code消息
func CreateUser(data *User) int {
	data.Password = ScryptPw(data.Password) // 对用户密码加密
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCESS // 200
}

// 密码加密,输入用户密码，输出加密后字符串
func ScryptPw(password string) string {
	const KeyLen = 10 // KeyLen表示加密秘钥的字节长度
	salt := make([]byte, 8)
	salt = []byte{11, 23, 43, 63, 49, 64, 43, 54} // 随机加8个salt值用于加密
	// HashPw为加密后的字节切片
	HashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen) //N是CPU/内存成本参数，r和p必须满足r*p<2^30
	if err != nil {
		log.Fatal(err)
	}
	fpw := base64.StdEncoding.EncodeToString(HashPw) // 转成字符串存在数据库中
	return fpw
}

// 登录验证
func CheckLogin(username string, password string) int {
	var user User
	db.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		return errmsg.ERROR_USER_NOT_EXIST
	}
	if ScryptPw(password) != user.Password {
		return errmsg.ERROR_PASSWORD_WRONG
	}
	return errmsg.SUCCESS
}
