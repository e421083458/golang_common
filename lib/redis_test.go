package lib

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"testing"
)

func Test_Redis(t *testing.T) {
	InitTest()

	c,err:=RedisConnFactory("default")
	if err!=nil{
		t.Fatal(err)
	}
	defer c.Close()

	// 调用SET
	trace:=NewTrace()
	redisKey:="test_key1"
	RedisLogDo(trace, c,"SET", redisKey, "test_dpool")
	RedisLogDo(trace, c,"expire", "test_key1", 10)

	// 调用GET
	v, err := redis.String(RedisLogDo(trace, c,"GET", redisKey))
	fmt.Println(v)
	if v!="test_dpool" || err!=nil{
		t.Fatal("test redis get fatal!")
	}

	DestroyTest()
}