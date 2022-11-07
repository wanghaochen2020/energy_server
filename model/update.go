package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

const (
	sleepTime    = 1           // 轮询间隔
	spanHalfHour = 30          // 30分钟
	spanDay      = 60 * 24     // 一天
	spanWeek     = 60 * 24 * 7 // 一周
	//spanTest     = 5
)

// 保存各个数据表的更新时间，保存时间戳单位为ms
type Update struct {
	Id        int
	UpdatedAt string
	TableName string
}

// 检查时间是否超时
func checkTime(pTime string, nowTime time.Time, timeout int) bool {
	target, err := time.ParseInLocation("2006-01-02 15:04:05", pTime, time.Local)
	CheckErr(err)
	if t1 := nowTime.Sub(target).Minutes(); int(t1) >= timeout {
		return true
	}
	return false
}

// 输入表名判断是否需要更新
func checkUpdate(table Update, now time.Time) {
	//var collection *mongo.Collection
	switch table.TableName {
	case "atmosphere":
		/*
			if checkTime(table.UpdatedAt, now, 1) {
				// 超过时间间隔，更新
				newData := receive.GetAtmosphere()
				collection = Db.Collection("atmosphere") //环境水电客流消耗品
				_, err := collection.InsertOne(context.TODO(), bson.D{{"data", newData}})
				if err != nil {
					log.Println("mongodb insert fail")
				}
				// 更新update表
				_, err = Db.Collection("update").UpdateOne(context.TODO(), bson.D{{
					Key:   "tableName",
					Value: table.TableName,
				}}, bson.D{
					{"$set", bson.D{{
						Key:   "updatedAt",
						Value: now.Format("2006-01-02 15:04:05"),
					}}},
				})
				CheckErr(err)
			}
		*/

	}
}

// 循环查询是否更新
func LoopQueryUpdate() {
	func() {
		var updates []Update
		for {
			fmt.Println("1111")
			// 根据规则过滤数据，这里过滤条件为空
			data, err := Db.Collection("update").Find(context.TODO(), bson.D{})
			CheckErr(err)
			// 解析数据到数组中
			err = data.All(context.TODO(), &updates)
			CheckErr(err)
			now := time.Now()
			for _, v := range updates {
				checkUpdate(v, now)
			}
			_ = data.Close(context.TODO())
			// 睡眠一分钟释放cpu
			time.Sleep(time.Minute * time.Duration(sleepTime))
		}
	}()
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}
