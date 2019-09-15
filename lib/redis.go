package lib

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"time"
)

func RedisConnFactory(name string) (redis.Conn, error) {
	if ConfRedisMap != nil && ConfRedisMap.List != nil {
		for confName, cfg := range ConfRedisMap.List {
			if name == confName {
				randHost := cfg.ProxyList[rand.Intn(len(cfg.ProxyList))]
				if cfg.ConnTimeout == 0 {
					cfg.ConnTimeout = 50
				}
				if cfg.ReadTimeout == 0 {
					cfg.ReadTimeout = 100
				}
				if cfg.WriteTimeout == 0 {
					cfg.WriteTimeout = 100
				}
				return redis.Dial(
					"tcp",
					randHost,
					redis.DialConnectTimeout(time.Duration(cfg.ConnTimeout)*time.Millisecond),
					redis.DialReadTimeout(time.Duration(cfg.ReadTimeout)*time.Millisecond),
					redis.DialWriteTimeout(time.Duration(cfg.WriteTimeout)*time.Millisecond))
			}
		}
	}
	return nil, errors.New("create redis conn fail")
}

func RedisLogDo(trace *TraceContext, c redis.Conn, commandName string, args ...interface{}) (interface{}, error) {
	startExecTime := time.Now()
	reply, err := c.Do(commandName, args...)
	endExecTime := time.Now()
	if err != nil {
		Log.TagError(trace, "_com_redis_failure", map[string]interface{}{
			"method":    commandName,
			"err":       err,
			"bind":      args,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		replyStr, _ := redis.String(reply, nil)
		Log.TagInfo(trace, "_com_redis_success", map[string]interface{}{
			"method":    commandName,
			"bind":      args,
			"reply":     replyStr,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return reply, err
}

//通过配置 执行redis
func RedisConfDo(trace *TraceContext, name string, commandName string, args ...interface{}) (interface{}, error) {
	c, err := RedisConnFactory(name)
	if err != nil {
		Log.TagError(trace, "_com_redis_failure", map[string]interface{}{
			"method": commandName,
			"err":    errors.New("RedisConnFactory_error:" + name),
			"bind":   args,
		})
		return nil, err
	}
	defer c.Close()

	startExecTime := time.Now()
	reply, err := c.Do(commandName, args...)
	endExecTime := time.Now()
	if err != nil {
		Log.TagError(trace, "_com_redis_failure", map[string]interface{}{
			"method":    commandName,
			"err":       err,
			"bind":      args,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		replyStr, _ := redis.String(reply, nil)
		Log.TagInfo(trace, "_com_redis_success", map[string]interface{}{
			"method":    commandName,
			"bind":      args,
			"reply":     replyStr,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return reply, err
}
