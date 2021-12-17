package lib

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

func InitDBPool(path string) error {
	//普通的db方式
	DbConfMap := &MysqlMapConf{}
	err := ParseConfig(path, DbConfMap)
	if err != nil {
		return err
	}
	if len(DbConfMap.List) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(TimeFormat), " empty mysql config.")
	}

	DBMapPool = map[string]*sql.DB{}
	GORMMapPool = map[string]*gorm.DB{}
	for confName, DbConf := range DbConfMap.List {
		dbpool, err := sql.Open("mysql", DbConf.DataSourceName)
		if err != nil {
			return err
		}
		dbpool.SetMaxOpenConns(DbConf.MaxOpenConn)
		dbpool.SetMaxIdleConns(DbConf.MaxIdleConn)
		dbpool.SetConnMaxLifetime(time.Duration(DbConf.MaxConnLifeTime) * time.Second)
		err = dbpool.Ping()
		if err != nil {
			return err
		}

		//gorm连接方式
		dbGorm, err := gorm.Open(mysql.New(mysql.Config{Conn: dbpool}), &gorm.Config{
			Logger: &DefaultMysqlGormLogger,
		})
		if err != nil {
			return err
		}
		DBMapPool[confName] = dbpool
		GORMMapPool[confName] = dbGorm
	}

	//手动配置连接
	if dbpool, err := GetDBPool("default"); err == nil {
		DBDefaultPool = dbpool
	}
	if dbpool, err := GetGormPool("default"); err == nil {
		GORMDefaultPool = dbpool
	}
	return nil
}

func GetDBPool(name string) (*sql.DB, error) {
	if dbpool, ok := DBMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get pool error")
}

func GetGormPool(name string) (*gorm.DB, error) {
	if dbpool, ok := GORMMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get pool error")
}

func CloseDB() error {
	for _, dbpool := range DBMapPool {
		dbpool.Close()
	}
	DBMapPool = make(map[string]*sql.DB)
	GORMMapPool = make(map[string]*gorm.DB)
	return nil
}

func DBPoolLogQuery(trace *TraceContext, sqlDb *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	startExecTime := time.Now()
	rows, err := sqlDb.Query(query, args...)
	endExecTime := time.Now()
	if err != nil {
		Log.TagError(trace, "_com_mysql_success", map[string]interface{}{
			"sql":       query,
			"bind":      args,
			"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		Log.TagInfo(trace, "_com_mysql_success", map[string]interface{}{
			"sql":       query,
			"bind":      args,
			"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return rows, err
}

//mysql 日志打印类型
var DefaultMysqlGormLogger = MysqlGormLogger{
	LogLevel:logger.Info,
	SlowThreshold:200 * time.Millisecond,
}

type MysqlGormLogger struct {
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

func (mgl *MysqlGormLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	mgl.LogLevel = logLevel
	return mgl
}

func (mgl *MysqlGormLogger) Info(ctx context.Context, message string, values ...interface{}) {
	trace := GetTraceContext(ctx)
	params := make(map[string]interface{})
	params["message"] = message
	params["values"] = fmt.Sprint(values)
	Log.TagInfo(trace, "_com_mysql_Info", params)
}
func (mgl *MysqlGormLogger) Warn(ctx context.Context, message string, values ...interface{}) {
	trace := GetTraceContext(ctx)
	params := make(map[string]interface{})
	params["message"] = message
	params["values"] = fmt.Sprint(values)
	Log.TagInfo(trace, "_com_mysql_Warn", params)
}
func (mgl *MysqlGormLogger) Error(ctx context.Context, message string, values ...interface{}) {
	trace := GetTraceContext(ctx)
	params := make(map[string]interface{})
	params["message"] = message
	params["values"] = fmt.Sprint(values)
	Log.TagInfo(trace, "_com_mysql_Error", params)
}
func (mgl *MysqlGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	trace := GetTraceContext(ctx)

	if mgl.LogLevel <= logger.Silent {
		return
	}

	sqlStr, rows := fc()
	currentTime := begin.Format(TimeFormat)
	elapsed := time.Since(begin)
	msg := map[string]interface{}{
		"FileWithLineNum": utils.FileWithLineNum(),
		"sql":             sqlStr,
		"rows":            "-",
		"proc_time":       float64(elapsed.Milliseconds()),
		"current_time":    currentTime,
	}
	switch {
	case err != nil && mgl.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound)):
		msg["err"] = err
		if rows == -1 {
			Log.TagInfo(trace, "_com_mysql_failure", msg)
		} else {
			msg["rows"] = rows
			Log.TagInfo(trace, "_com_mysql_failure", msg)
		}
	case elapsed > mgl.SlowThreshold && mgl.SlowThreshold != 0 && mgl.LogLevel >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", mgl.SlowThreshold)
		msg["slowLog"] = slowLog
		if rows == -1 {
			Log.TagInfo(trace, "_com_mysql_success", msg)
		} else {
			msg["rows"] = rows
			Log.TagInfo(trace, "_com_mysql_success", msg)
		}
	case mgl.LogLevel == logger.Info:
		if rows == -1 {
			Log.TagInfo(trace, "_com_mysql_success", msg)
		} else {
			msg["rows"] = rows
			Log.TagInfo(trace, "_com_mysql_success", msg)
		}
	}
}
