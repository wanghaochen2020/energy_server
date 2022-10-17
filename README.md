### 文件目录

```go
energy                 
├─ api                 // api接口实现
│  ├─ analysis         
│  ├─ system           
│  └─ login.go         
├─ config              // 保存项目初始化所需数据
│  └─ config.ini       
├─ dataReceive         // 存放接受设备数据的接口
├─ log                 // 日志打印文件夹
│  ├─ log              
│  └─ log20220513.log  
├─ middleware          // 中间件
│  ├─ cors.go          // 跨域中间件
│  ├─ jwt.go           // jwt身份鉴权中间件
│  └─ logger.go        // 打印日志中间件
├─ model               // 存放子功能所需结构体和一些方法
│  ├─ analysis         
│  ├─ system           
│  ├─ db.go            // 连接数据库
│  └─ User.go          
├─ routes              // 存放各个子功能路由文件 
│  └─ login.go         
├─ utils               // 工具包
│  ├─ errmsg           
│  │  └─ errmsg.go     // 存放错误代码
│  └─ setting.go       // 读取config.ini中的数据并加载
├─ go.mod              
├─ go.sum              
├─ main.go             // main函数文件
├─ README.md           
└─ router.go           // 路由文件
```

### mongo 相关数据文档

#### calculation_result：存储计算结果
```json
{
    "time": "2022/05/01",
    "name": "boiler_efficiency_day",//name见redis的table_name表格
    "value": [0.9, 0.8, 0, 0.88]
}
```

### redis文档
使用db2，存储表格所需数据，命名为"2022/05/01 08 table_name"，前面表示本地时间的年月日时，table_name为表格名字

| 表格         | table_name            | time精确到         | 值格式                                                |
| ------------ | --------------------- | ------------------ | ----------------------------------------------------- |
| 电锅炉热效率 | boiler_efficiency_day | 日（"2022/05/01"） | 长度为0到23的数组（根据记录值变化），如[90, 50, 0...] |
|              |                       |                    |                                                       |
|              |                       |                    |                                                       |

