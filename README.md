
## 定位
配置 Golang 基础服务（mysql、redis、http.client、log）比较繁琐，如果想 **快速接入** 基础服务可以使用本类库。
没有多余复杂的功能，方便你拓展其他功能。
你可以 import 引入直接使用，也可以拷贝代码到自己项目中使用，也可以用于构建自己的基础类库。
项目地址：https://github.com/e421083458/golang_common

## 功能
 1. 多套配置环境设置，比如：dev、prod。
 2. mysql、redis 多套数据源配置。
 3. 支持默认和自定义日志实例，自动滚动日志。
 4. 支持 mysql(基于gorm的二次开发支持ctx功能，不影响gorm原功能使用)、redis(redigo)、http.client 请求链路日志输出。
## 安装及使用
 1. 需要确保已经安装了 Go 1.8+，然后执行以下命令
```
go get -v github.com/e421083458/golang_common
```
2. 将配置文件拷贝到你的项目中，配置文件请参考：https://github.com/e421083458/golang_common/tree/master/conf/dev
可以通过 InitModule("配置地址","模块数组") 方法按模块需加载配置。

- base：包含日志和系统时间配置等
- mysql：包含mysql操作方法和日志
- redis：包含redis操作方法和日志


3. log消息打印代码举例：
```
package main

import (
	"github.com/e421083458/golang_common/lib"
	"log"
	"time"
)

func main() {
	if err := lib.InitModule("./conf/dev/",[]string{"base","mysql","redis",}); err != nil {
		log.Fatal(err)
	}
	defer lib.Destroy()

	//todo sth
	lib.Log.TagInfo(lib.NewTrace(), lib.DLTagUndefind, map[string]interface{}{"message": "todo sth"})
	time.Sleep(time.Second)
}
```
4. 运行代码

```
go run main.go
```

输出：

```
2019/05/21 10:31:08 ------------------------------------------------------------------------
2019/05/21 10:31:08 [INFO]  config=./conf/dev/
2019/05/21 10:31:08 [INFO]  start loading resources.
2019/05/21 10:31:08 [INFO]  success loading resources.
2019/05/21 10:31:08 ------------------------------------------------------------------------
2019-05-21T10:31:08.667 [INFO] log.go:58 _undef||traceid=ac182a1c5ce362ecf9280618104dc7b0||cspanid=||spanid=f0fb48f0380704bb||message=todo sth
2019/05/21 10:31:09 ------------------------------------------------------------------------
2019/05/21 10:31:09 [INFO]  start destroy resources.
2019/05/21 10:31:09 [INFO]  success destroy resources.
```

## 其他功能举例

- 初始化当前运行环境

```
//初始化测试用例
func SetUp()  {
	initOnce.Do(func() {
		if err:=lib.Init("../conf/dev/");err!=nil{
			log.Fatal(err)
		}
	})
}
```

- 获取当前运行环境

```
//获取 程序运行环境 dev prod
func Test_GetConfEnv(t *testing.T) {
	SetUp()
	fmt.Println(lib.GetConfEnv())
	TearDown()
}
```


- 以子树的形式获取配置

类似方法有：
- GetConf(key string) : interface{}
- GetBoolConf(key string) : bool
- GetFloat64Conf(key string) : float64
- GetIntConf(key string) : int
- GetStringConf(key string) : string
- GetStringMapConf(key string) : map[string]interface{}
- GetStringMapStringConf(key string) : map[string]string
- GetStringSliceConf(key string) : []string
- GetTimeConf(key string) : time.Time
- GetDurationConf(key string) : time.Duration
- IsSetConf(key string) : bool

```
func TestGetStringConf(t *testing.T) {
	SetUp()
	got := lib.GetStringConf("base.log.log_level")
	if got!="trace"{
		t.Fatal("got result error")
	}
}
```


- 加载自定义配置文件

```
type HttpConf struct {
	ServerAddr     string   `mapstructure:"server_addr"`
	ReadTimeout    int      `mapstructure:"read_timeout"`
	WriteTimeout   int      `mapstructure:"write_timeout"`
	MaxHeaderBytes int      `mapstructure:"max_header_bytes"`
	AllowHost      []string `mapstructure:"allow_host"`
}
// 加载自定义配置文件
func Test_ParseLocalConfig(t *testing.T) {
	SetUp()
	httpProfile := &HttpConf{}
	err:=lib.ParseLocalConfig("http.toml",httpProfile)
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(httpProfile)
	TearDown()
}
```

- 测试PostJson请求

```
//测试PostJson请求
func TestJson(t *testing.T) {
	SetUpServer()
	//首次scrollsId不传递
	jsonStr := "{\"source\":\"control\",\"cityId\":\"12\",\"trailNum\":10,\"dayTime\":\"2018-11-21 16:08:00\",\"limit\":2,\"andOperations\":{\"cityId\":\"eq\",\"trailNum\":\"gt\",\"dayTime\":\"eq\"}}"
	url := "http://"+addr+"/postjson"
	_, res, err := lib.HttpJSON(lib.NewTrace(), url, jsonStr, 1000, nil)
	fmt.Println(string(res))
	if err != nil {
		fmt.Println(err.Error())
	}
}
```

- 测试Get请求

```
//测试Get请求
func TestGet(t *testing.T) {
	SetUpServer()
	a := url.Values{
		"city_id": {"12"},
	}
	url := "http://"+addr+"/get"
	_, res, err := lib.HttpGET(lib.NewTrace(), url, a, 1000, nil)
	fmt.Println("city_id="+string(res))
	if err != nil {
		fmt.Println(err.Error())
	}
}
```

- 测试Post请求

```
//测试Post请求
func TestPost(t *testing.T) {
	SetUpServer()
	a := url.Values{
		"city_id": {"12"},
	}
	url := "http://"+addr+"/post"
	_, res, err := lib.HttpPOST(lib.NewTrace(), url, a, 1000, nil, "")
	fmt.Println("city_id="+string(res))
	if err != nil {
		fmt.Println(err.Error())
	}
}
```

- 测试默认日志打点

```
//测试日志打点
func TestDefaultLog(t *testing.T) {
	SetUp()
	lib.Log.TagInfo(lib.NewTrace(), lib.DLTagMySqlSuccess, map[string]interface{}{
		"sql": "sql",
	})
	time.Sleep(time.Second)
	TearDown()
}
```
- 测试自定义日志实例打点

```
//测试日志实例打点
func TestLogInstance(t *testing.T) {
	nlog:= log.NewLogger()
	logConf:= log.LogConfig{
		Level:"trace",
		FW: log.ConfFileWriter{
			On:true,
			LogPath:"./log_test.log",
			RotateLogPath:"./log_test.log",
			WfLogPath:"./log_test.wf.log",
			RotateWfLogPath:"./log_test.wf.log",
		},
		CW: log.ConfConsoleWriter{
			On:true,
			Color:true,
		},
	}
	log.SetupLogInstanceWithConf(logConf,nlog)
	nlog.Info("test message")
	nlog.Close()
	time.Sleep(time.Second)
}
```

- 测试mysql普通sql

```
var (
	createTableSQL = "CREATE TABLE `test1` (`id` int(12) unsigned NOT NULL AUTO_INCREMENT" +
		" COMMENT '自增id',`name` varchar(255) NOT NULL DEFAULT '' COMMENT '姓名'," +
		"`created_at` datetime NOT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB " +
		"DEFAULT CHARSET=utf8"
	insertSQL    = "INSERT INTO `test1` (`id`, `name`, `created_at`) VALUES (NULL, '111', '2018-08-29 11:01:43');"
	dropTableSQL = "DROP TABLE `test1`"
	beginSQL     = "start transaction;"
	commitSQL    = "commit;"
	rollbackSQL  = "rollback;"
)

func Test_DBPool(t *testing.T) {
	SetUp()

	//获取链接池
	dbpool, err := lib.GetDBPool("default")
	if err != nil {
		t.Fatal(err)
	}
	//开始事务
	trace := lib.NewTrace()
	if _, err := lib.DBPoolLogQuery(trace, dbpool, beginSQL); err != nil {
		t.Fatal(err)
	}

	//创建表
	if _, err := lib.DBPoolLogQuery(trace, dbpool, createTableSQL); err != nil {
		lib.DBPoolLogQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}

	//插入数据
	if _, err := lib.DBPoolLogQuery(trace, dbpool, insertSQL); err != nil {
		lib.DBPoolLogQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}

	//循环查询数据
	current_id := 0
	table_name := "test1"
	fmt.Println("begin read table ", table_name, "")
	fmt.Println("------------------------------------------------------------------------")
	fmt.Printf("%6s | %6s\n", "id", "created_at")
	for {
		rows, err := lib.DBPoolLogQuery(trace, dbpool, "SELECT id,created_at FROM test1 WHERE id>? order by id asc", current_id)
		defer rows.Close()
		row_len := 0
		if err != nil {
			lib.DBPoolLogQuery(trace, dbpool, "rollback;")
			t.Fatal(err)
		}
		for rows.Next() {
			var create_time string
			if err := rows.Scan(&current_id, &create_time); err != nil {
				lib.DBPoolLogQuery(trace, dbpool, "rollback;")
				t.Fatal(err)
			}
			fmt.Printf("%6d | %6s\n", current_id, create_time)
			row_len++
		}
		if row_len == 0 {
			break
		}
	}
	fmt.Println("------------------------------------------------------------------------")
	fmt.Println("finish read table ", table_name, "")

	//删除表
	if _, err := lib.DBPoolLogQuery(trace, dbpool, dropTableSQL); err != nil {
		lib.DBPoolLogQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}

	//提交事务
	lib.DBPoolLogQuery(trace, dbpool, commitSQL)
	TearDown()
}
```

- 测试Gorm

```
func Test_GORM(t *testing.T) {
	SetUp()

	//获取链接池
	dbpool, err := lib.GetGormPool("default")
	if err != nil {
		t.Fatal(err)
	}
	db := dbpool.Begin()
	traceCtx := lib.NewTrace()

	//设置trace信息
	db = db.SetCtx(traceCtx)
	if err := db.Exec(createTableSQL).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	//插入数据
	t1 := &Test1{Name: "test_name", CreatedAt: time.Now()}
	if err := db.Save(t1).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	//查询数据
	list := []Test1{}
	if err := db.Where("name=?", "test_name").Find(&list).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	fmt.Println(list)

	//删除表数据
	if err := db.Exec(dropTableSQL).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	db.Commit()
	TearDown()
}
```

- 测试redis查询

```
func Test_Redis(t *testing.T) {
	SetUp()
	
	c, err := lib.RedisConnFactory("default")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// 调用SET
	trace := lib.NewTrace()
	redisKey := "test_key1"
	lib.RedisLogDo(trace, c, "SET", redisKey, "test_dpool")
	lib.RedisLogDo(trace, c, "expire", "test_key1", 10)

	// 调用GET
	v, err := redis.String(lib.RedisLogDo(trace, c, "GET", redisKey))
	fmt.Println(v)
	if v != "test_dpool" || err != nil {
		t.Fatal("test redis get fatal!")
	}

	TearDown()
}
```

- 销毁当前运行环境

```
//销毁测试用例
func TearDown()  {
	Destroy()
}
```
## 思考
**怎么样才能让类库做的更通用？**
有人说简洁才是golang特性，不需要做通用类库。如果你想加一个功能直接引一个包用即可。比如：日志那就引日志包、redis就引redis包。

我赞同此观点，但是引入包后要面临的问题是功能改造和代码适配，并且每做一个项目都要搞上一套。这意味着需要耗费一定的时间在重复的工作上。

我感觉类库应该有一下特点：

 1. 轻量级，没有太多依赖，否则容易有类库依赖冲突。
 2. 只封装重复使用率高的功能。
 3. 拓展性强。

你的 star ，我的动力。有任何关于类库的问题，请提交issue，谢谢。