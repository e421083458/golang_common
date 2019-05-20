package lib

import (
	"testing"
	"time"
)

//测试日志打点
func TestDefaultLog(t *testing.T) {
	InitTest()
	Log.TagInfo(NewTrace(), DLTagMySqlSuccess, map[string]interface{}{
		"sql": "sql",
	})
	time.Sleep(time.Second)
	DestroyTest()
}