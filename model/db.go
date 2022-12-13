package model

import (
	"context"
	"energy/defs"
	"energy/utils"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db           *gorm.DB
	err          error
	RedisClient  *redis.Client
	MongoClient  *mongo.Client
	MongoOPC     *mongo.Collection
	MongoResult  *mongo.Collection
	MongoLoukong *mongo.Collection
)
var Db *mongo.Database

func InitDb() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.DbUser,
		utils.DbPassWord,
		utils.DbHost,
		utils.DbPort,
		utils.DbName,
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// gorm日志模式：silent
		Logger: logger.Default.LogMode(logger.Silent),
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Println("连接数据库失败，请检查参数：", err)
		os.Exit(1)
	}

	// 迁移数据表，在没有数据表结构变更时候，建议注释不执行
	_ = db.AutoMigrate(&User{})

	sqlDB, _ := db.DB()
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置连接的最大可复用时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)
}

func InitRedis() {
	db, err := strconv.Atoi(utils.RedisDbNum)
	if err != nil {
		log.Panicf("Redis DB num is not int type: %s", err)
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", utils.RedisHost, utils.RedisPort),
		Password: utils.RedisPassWord,
		DB:       db,
	})
	_, err = RedisClient.Ping().Result()
	if err != nil {
		log.Panicf("Redis connection error: %s", err)
	}
}

func InitMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	MongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(utils.MongoUrl))
	Db = MongoClient.Database(utils.MongoName)
	if err != nil {
		log.Panicf("MongoDB connection error: %s", err)
	}

	err = MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Panicf("MongoDB ping error: %s", err)
	}

	MongoOPC = MongoClient.Database("energy").Collection("opc_data")
	MongoResult = MongoClient.Database("energy").Collection("calculation_result")
	MongoLoukong = MongoClient.Database("energy").Collection("loukong")
}

func MongoUpdateList(timeStr string, index int, name string, value float64) {
	var result defs.CalculationResultFloatList
	err = MongoResult.FindOne(context.TODO(), bson.D{{"time", timeStr}, {"name", name}}).Decode(&result)
	l := utils.Max(len(result.Value), index+1)
	finalData := make([]float64, l)
	copy(finalData, result.Value)
	finalData[index] = value
	if err == mongo.ErrNoDocuments {
		_, err = MongoResult.InsertOne(context.TODO(), bson.D{{"time", timeStr}, {"name", name}, {"value", finalData}})
		if err != nil {
			log.Print(err)
		}
	} else {
		_, err = MongoResult.UpdateOne(context.TODO(), bson.D{{"time", timeStr}, {"name", name}}, bson.D{{"$set", bson.D{{"value", finalData}}}})
		if err != nil {
			log.Print(err)
		}
	}
}

func MongoUpsertOne(name string, value interface{}) {
	opts := options.Update().SetUpsert(true)
	_, err = MongoResult.UpdateOne(context.TODO(), bson.D{{"name", name}}, bson.D{{"$set", bson.D{{"value", value}}}}, opts)
	if err != nil {
		log.Print(err)
	}
}
