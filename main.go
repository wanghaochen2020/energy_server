package main

import "energy/model"

func main() {
	//引用数据库
	model.InitDb()
	//引用路由组件
	InitRouter()
}

