
## 定位
配置 Golang 基础服务（mysql、redis、http.client、log）比较繁琐，如果想 **快速接入** 基础服务可以使用本类库。
没有多余复杂的功能，方便你拓展其他功能。
你可以 import 引用使用，也可以拷贝代码到自己项目中使用。
项目地址：https://github.com/e421083458/golang_common

## 功能
 1. 多套配置环境设置，比如：dev、prod。
 2. mysql、redis 多套数据源配置。
 3. 支持多个自定义日志实例，自动滚动切分日志。
 4. 支持 mysql(gorm)、redis(redigo)、http.client 请求链路日志输出。
 5. 支持多个自定义日志实例。
## 安装及使用
 1. 需要确保已经安装了 Go 1.8+，然后执行以下命令
```
go get -v github.com/e421083458/golang_common
```
2. 将配置文件拷贝到你的项目中，配置文件请参考：https://github.com/e421083458/golang_common/tree/master/conf/dev

3. 引入到你的项目：
```
package main

import (
	"github.com/e421083458/golang_common/lib"
	"log"
	"time"
)

func main(){
	if err:=lib.Init("./conf/dev/");err!=nil{
		log.Fatal(err)
	}

	//todo sth
	lib.Log.TagInfo(lib.NewTrace(), lib.DLTagUndefind, map[string]interface{}{"message": "todo sth"})
	time.Sleep(time.Second)

	lib.Destroy()
}
```
4. 运行代码

```
go run main.go
```

输出：

```
2019/05/19 18:55:17 [INFO]  config=./conf/dev/
2019/05/19 18:55:17 [INFO]  start loading resources.
2019/05/19 18:55:17 [INFO]  success loading resources.
------------------------------------------------------------------------
2019-05-19T18:55:17.783 [INFO] log.go:58 _undef||message=todo sth||traceid=c0a8fe315ce13615b6c0fc6e104dc7b0||cspanid=||spanid=9c49c824380704bb
2019/05/19 18:55:18 [INFO]  start destroy resources.
------------------------------------------------------------------------
2019/05/19 18:55:18 [INFO]  success destroy resources.
```

## 其他功能举例

- 初始化当前运行环境

```
//初始化测试用例
func InitTest()  {
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
	InitTest()
	fmt.Println(lib.GetConfEnv())
	DestroyTest()
}
```

- 加载自定义配置文件

```
type HttpConf struct {
	ServerAddr     string   `toml:"server_addr"`
	ReadTimeout    int      `toml:"read_timeout"`
	WriteTimeout   int      `toml:"write_timeout"`
	MaxHeaderBytes int      `toml:"max_header_bytes"`
	AllowHost      []string `toml:"allow_host"`
}
// 加载自定义配置文件
func Test_ParseLocalConfig(t *testing.T) {
	InitTest()
	httpProfile := &HttpConf{}
	err:=lib.ParseLocalConfig("http.toml",httpProfile)
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(httpProfile)
	DestroyTest()
}
```

- 测试PostJson请求

```
//测试PostJson请求
func TestJson(t *testing.T) {
	InitTestServer()
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
	InitTestServer()
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
	InitTestServer()
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
	InitTest()
	lib.Log.TagInfo(lib.NewTrace(), lib.DLTagMySqlSuccess, map[string]interface{}{
		"sql": "sql",
	})
	time.Sleep(time.Second)
	DestroyTest()
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
	InitTest()

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
	DestroyTest()
}
```

- 测试Gorm

```
func Test_GORM(t *testing.T) {
	InitTest()

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
	DestroyTest()
}
```

- 测试redis查询

```
func Test_Redis(t *testing.T) {
	InitTest()
	
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

	DestroyTest()
}
```

- 销毁当前运行环境

```
//销毁测试用例
func DestroyTest()  {
	Destroy()
}
```
你的 star ，我的动力。有任何关于类库的问题，请提交issue，谢谢。