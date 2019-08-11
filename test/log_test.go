package test

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/e421083458/golang_common/log"
	"testing"
	"time"
)

//测试日志打点
func TestDefaultLog(t *testing.T) {
	SetUp()
	lib.Log.TagInfo(lib.NewTrace(), lib.DLTagMySqlSuccess, map[string]interface{}{
		"sql": "sql",
	})
	time.Sleep(time.Second)
	TearDown()
}

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