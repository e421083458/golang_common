package test

import (
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/garyburd/redigo/redis"
	"testing"
)

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
