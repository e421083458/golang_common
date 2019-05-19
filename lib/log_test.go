package lib

import (
	"testing"
	"time"
)

//测试日志打点
func TestInitLog(t *testing.T) {
	InitTest()
	Log.TagInfo(NewTrace(), DLTagMySqlSuccess, map[string]interface{}{
		"sql": "dltag",
	})
	time.Sleep(time.Second)
	DestroyTest()
}