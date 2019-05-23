package test

import (
	"bytes"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

var (
	addr       string    = "127.0.0.1:6111"
	initOnce   sync.Once = sync.Once{}
	serverOnce sync.Once = sync.Once{}
)

type HttpConf struct {
	ServerAddr     string   `mapstructure:"server_addr"`
	ReadTimeout    int      `mapstructure:"read_timeout"`
	WriteTimeout   int      `mapstructure:"write_timeout"`
	MaxHeaderBytes int      `mapstructure:"max_header_bytes"`
	AllowHost      []string `mapstructure:"allow_host"`
}

//获取 程序运行环境 dev prod
func Test_GetConfEnv(t *testing.T) {
	SetUp()
	fmt.Println(lib.GetConfEnv())
	TearDown()
}

// 加载自定义配置文件
func Test_ParseLocalConfig(t *testing.T) {
	SetUp()
	httpProfile := &HttpConf{}
	err := lib.ParseLocalConfig("test.toml", httpProfile)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(httpProfile)
	TearDown()
}

//测试PostJson请求
func TestJson(t *testing.T) {
	InitTestServer()
	//首次scrollsId不传递
	jsonStr := "{\"source\":\"control\",\"cityId\":\"12\",\"trailNum\":10,\"dayTime\":\"2018-11-21 16:08:00\",\"limit\":2,\"andOperations\":{\"cityId\":\"eq\",\"trailNum\":\"gt\",\"dayTime\":\"eq\"}}"
	url := "http://" + addr + "/postjson"
	_, res, err := lib.HttpJSON(lib.NewTrace(), url, jsonStr, 1000, nil)
	fmt.Println(string(res))
	if err != nil {
		fmt.Println(err.Error())
	}
}

//测试Get请求
func TestGet(t *testing.T) {
	InitTestServer()
	a := url.Values{
		"city_id": {"12"},
	}
	url := "http://" + addr + "/get"
	_, res, err := lib.HttpGET(lib.NewTrace(), url, a, 1000, nil)
	fmt.Println("city_id=" + string(res))
	if err != nil {
		fmt.Println(err.Error())
	}
}

//测试Post请求
func TestPost(t *testing.T) {
	InitTestServer()
	a := url.Values{
		"city_id": {"12"},
	}
	url := "http://" + addr + "/post"
	_, res, err := lib.HttpPOST(lib.NewTrace(), url, a, 1000, nil, "")
	fmt.Println("city_id=" + string(res))
	if err != nil {
		fmt.Println(err.Error())
	}
}

//初始化测试用例
func SetUp() {
	initOnce.Do(func() {
		if err := lib.InitModule("../conf/dev/", []string{"base", "mysql", "redis",}); err != nil {
			log.Fatal(err)
		}
	})
}

//销毁测试用例
func TearDown() {
	//lib.Destroy()
}

//只运行一次服务器
func InitTestServer() {
	serverOnce.Do(func() {
		http.HandleFunc("/postjson", func(writer http.ResponseWriter, request *http.Request) {
			bodyBytes, _ := ioutil.ReadAll(request.Body)
			request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Write body back
			writer.Write([]byte(bodyBytes))
		})
		http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
			request.ParseForm()
			cityID := request.FormValue("city_id")
			writer.Write([]byte(cityID))
		})
		http.HandleFunc("/post", func(writer http.ResponseWriter, request *http.Request) {
			request.ParseForm()
			cityID := request.FormValue("city_id")
			writer.Write([]byte(cityID))
		})
		go func() {
			log.Println("ListenAndServe ", addr)
			err := http.ListenAndServe(addr, nil) //设置监听的端口
			if err != nil {
				log.Fatal("ListenAndServe: ", err)
			}
		}()
		time.Sleep(time.Second)
	})
}

//测试获取配置string
func TestGetStringConf(t *testing.T) {
	SetUp()
	got := lib.GetStringConf("base.log.log_level")
	if got!="trace"{
		t.Fatal("got result error")
	}
}
