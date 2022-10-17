package main

import "energy/model"

func main() {
	//引用数据库
	model.InitDb()
	model.InitRedis()
	model.InitMongo()
	//引用路由组件
	InitRouter()
}

/*逻辑：
前端请求，以每一个图表为单位，从redis中查询最近的计算结果，如若没有或不是最近的从mongo中查询最近的计算结果，如若没有则找最近的数据进行计算，计算出结果后存入mongo和redis


记得把路由放到token验证里
*/
