package lib

import (
	"database/sql"
	dlog "github.com/e421083458/golang_common/xlog4go"
	"github.com/jinzhu/gorm"
)

type BaseConf struct {
	DebugMode      string   `toml:"debug_mode"`
	TimeLocation   string   `toml:"time_location"`
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
	return nil
}

func InitLogger(path string) error {
	if err := dlog.SetupLogWithConf(path); err != nil {
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
