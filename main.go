package main

import (
	"energy/model"
)

func main() {
	//引用数据库
	model.InitDb()
	model.InitRedis()
	model.InitMongo()
	//model.LoopQueryUpdate()
	//引用路由组件
	InitRouter()
}

/*
有缺失数据的话，仍然会存储根据已有数据的计算结果，如果进行了数据的补充或更新，要去mongo里面删掉这些结果
记得换路由到token验证里
用的时候更改update.go中的exampleTime
*/
