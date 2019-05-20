package lib

import (
	"database/sql"
	dlog "github.com/e421083458/golang_common/log"
	"github.com/e421083458/gorm"
)

type BaseConf struct {
	DebugMode    string    `toml:"debug_mode"`
	TimeLocation string    `toml:"time_location"`
	Log          LogConfig `toml:"log"`
}

type LogConfFileWriter struct {
	On              bool   `toml:"on"`
	LogPath         string `toml:"log_pathLogPath"`
	RotateLogPath   string `toml:"rotate_log_path"`
	WfLogPath       string `toml:"wf_log_path"`
	RotateWfLogPath string `toml:"rotate_wf_path"`
}

type LogConfConsoleWriter struct {
	On    bool `toml:"on"`
	Color bool `toml:"color"`
}

type LogConfig struct {
	Level string               `toml:"log_level"`
	FW    LogConfFileWriter    `toml:"file_writer"`
	CW    LogConfConsoleWriter `toml:"console_writer"`
}

type MysqlMapConf struct {
	List map[string]*MySQLConf `toml:"list"`
}

type MySQLConf struct {
	DriverName      string `toml:"driver_name"`
	DataSourceName  string `toml:"data_source_name"`
	MaxOpenConn     int    `toml:"max_open_conn"`
	MaxIdleConn     int    `toml:"max_idle_conn"`
	MaxConnLifeTime int    `toml:"max_conn_life_time"`
}

type RedisMapConf struct {
	List map[string]*RedisConf `toml:"list"`
}

type RedisConf struct {
	ProxyList []string `toml:"proxy_list"`
	MaxActive int      `toml:"max_active"`
	MaxIdle   int      `toml:"max_idle"`
	Downgrade bool     `toml:"down_grade"`
}

//全局变量
var ConfBase *BaseConf
var DBMapPool map[string]*sql.DB
var GORMMapPool map[string]*gorm.DB
var DBDefaultPool *sql.DB
var GORMDefaultPool *gorm.DB
var ConfRedis *RedisConf
var ConfRedisMap *RedisMapConf

//获取基本配置信息
func GetBaseConf() *BaseConf {
	return ConfBase
}

func InitBaseConf(path string) error {
	ConfBase = &BaseConf{}
	err := ParseConfig(path, ConfBase)
	if err != nil {
		return err
	}
	//配置日志
	logConf := dlog.LogConfig{
		Level: ConfBase.Log.Level,
		FW: dlog.ConfFileWriter{
			On:              ConfBase.Log.FW.On,
			LogPath:         ConfBase.Log.FW.LogPath,
			RotateLogPath:   ConfBase.Log.FW.RotateLogPath,
			WfLogPath:       ConfBase.Log.FW.WfLogPath,
			RotateWfLogPath: ConfBase.Log.FW.RotateWfLogPath,
		},
		CW: dlog.ConfConsoleWriter{
			On:    ConfBase.Log.CW.On,
			Color: ConfBase.Log.CW.Color,
		},
	}
	if err := dlog.SetupDefaultLogWithConf(logConf); err != nil {
		panic(err)
	}
	dlog.SetLayout("2006-01-02T15:04:05.000")
	return nil
}

func InitLogger(path string) error {
	if err := dlog.SetupDefaultLogWithFile(path); err != nil {
		panic(err)
	}
	dlog.SetLayout("2006-01-02T15:04:05.000")
	return nil
}

func InitRedisConf(path string) error {
	ConfRedis := &RedisMapConf{}
	err := ParseConfig(path, ConfRedis)
	if err != nil {
		return err
	}
	ConfRedisMap = ConfRedis
	return nil
}
