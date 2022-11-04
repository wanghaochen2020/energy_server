package utils

import (
	"fmt"

	"gopkg.in/ini.v1"
)

// 全局变量
var (
	AppMode  string // 运行模式
	HttpPort string // 端口
	JwtKey   string // jwt密钥

	Db         string // 数据库类型
	DbHost     string // 数据库主机ip
	DbPort     string // 数据库端口
	DbUser     string // 数据库用户名
	DbPassWord string // 数据库密码
	DbName     string // 数据库表名

	RedisHost     string // redis主机ip
	RedisPort     string // redis端口
	RedisPassWord string // redis密码
	RedisDbNum    string // redis数据库号

	MongoUrl string //Mongodb连接字符串

	DeltaT int //计算服务的统计间隔
)

func init() {
	file, err := ini.Load("config/config.ini") //将配置文件引入，用于后续读取
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	LoadServer(file)
	LoadData(file)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug") //MustString表示默认参数
	HttpPort = file.Section("server").Key("HttpPort").MustString(":6666")
	JwtKey = file.Section("server").Key("JwtKey").MustString("889g5hfs9f")
}

func LoadData(file *ini.File) {
	Db = file.Section("database").Key("Db").MustString("mysql")
	DbHost = file.Section("database").Key("DbHost").MustString("10.112.154.218")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("root")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("Bupt@2021")
	DbName = file.Section("database").Key("DbName").MustString("Energy")

	RedisHost = file.Section("redis").Key("RedisHost").MustString("10.112.154.218")
	RedisPort = file.Section("redis").Key("RedisPort").MustString("6379")
	RedisPassWord = file.Section("redis").Key("RedisPassWord").MustString("123456")
	RedisDbNum = file.Section("redis").Key("RedisDbNum").MustString("2")

	MongoUrl = file.Section("mongodb").Key("MongoUrl").MustString("mongodb://admin:bupt2022@10.112.154.218:27017")

	DeltaT = file.Section("calculator").Key("DeltaT").MustInt(1)
}
