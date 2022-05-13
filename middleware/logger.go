package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	retalog "github.com/lestrrat-go/file-rotatelogs" // 日志分隔
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus" // 记录日志文件包
	"math"
	"os"
	"time"
)

func Logger() gin.HandlerFunc {
	// 写入log文件中
	filePath := "log/log"
	//linkName := "latestLog.log"

	src, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("err : ", err)
	}
	logger := logrus.New()
	// 输出到日志
	logger.Out = src
	// 下面将日志按时间分隔
	logger.SetLevel(logrus.DebugLevel)
	logWriter, _ := retalog.New(
		filePath+"%Y%m%d.log",                  // log名字加上年月日
		retalog.WithMaxAge(7*24*time.Hour),     // 最大保存时间一周
		retalog.WithRotationTime(24*time.Hour), // 24小时分隔一次
		//retalog.WithLinkName(linkName),         // 软连接，生成的log指向最新的log，永远保存最新的log（Windows下需要管理员权限）
	)

	writeMap := lfshook.WriterMap{ // 把不同日志级别写到logWriter中
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // 时间格式化
	})
	logger.AddHook(Hook) // 新增hook方法

	// 中间件函数
	return func(c *gin.Context) {
		startTime := time.Now() // 记录请求的开始时间
		c.Next()
		stopTime := time.Since(startTime)                                                            //时间段
		spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds())/1000000.0))) // 以毫秒为单位的持续时间
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}
		statusCode := c.Writer.Status()
		clientIp := c.ClientIP()
		userAgent := c.Request.UserAgent() // 记录连接方式，比如谷歌浏览器，或其他客户端的信息
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method
		path := c.Request.RequestURI

		entry := logger.WithFields(logrus.Fields{
			"HostName":  hostName,
			"Status":    statusCode,
			"SpendTime": spendTime,
			"Ip":        clientIp,
			"Method":    method,
			"Path":      path,
			"DataSize":  dataSize,
			"Agent":     userAgent,
		})
		if len(c.Errors) > 0 { // 系统内部有错误
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String()) // 记录系统内部错误
		}
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}

	}
}
