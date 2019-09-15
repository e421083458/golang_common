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
	vint, verr := redis.Int64(lib.RedisLogDo(trace, c, "INCR", "test_incr"))
	fmt.Println(vint)
	fmt.Println(verr)

	// 调用GET
	v, err := redis.String(lib.RedisLogDo(trace, c, "GET", redisKey))
	fmt.Println(v)
	fmt.Println(err)
	if v != "test_dpool" || err != nil {
		t.Fatal("test redis get fatal!")
	}

	// 使用RedisConfDo调用GET
	v2, err := redis.String(lib.RedisConfDo(trace, "default", "GET", redisKey))
	fmt.Println(v2)
	fmt.Println(err)
	if v != "test_dpool" || err != nil {
		t.Fatal("test redis get fatal!")
	}
	TearDown()
}
