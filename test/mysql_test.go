package test

import (
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/e421083458/gorm"
	"testing"
	"time"
)

type Test1 struct {
	Id        int64     `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (f *Test1) Table() string {
	return "test1"
}

func (f *Test1) DB() *gorm.DB {
	return lib.GORMDefaultPool
}

var (
	createTableSQL = "CREATE TABLE `test1` (`id` int(12) unsigned NOT NULL AUTO_INCREMENT" +
		" COMMENT '自增id',`name` varchar(255) NOT NULL DEFAULT '' COMMENT '姓名'," +
		"`created_at` datetime NOT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB " +
		"DEFAULT CHARSET=utf8"
	insertSQL    = "INSERT INTO `test1` (`id`, `name`, `created_at`) VALUES (NULL, '111', '2018-08-29 11:01:43');"
	dropTableSQL = "DROP TABLE `test1`"
	beginSQL     = "start transaction;"
	commitSQL    = "commit;"
	rollbackSQL  = "rollback;"
)

func Test_DBPool(t *testing.T) {
	SetUp()

	//获取链接池
	dbpool, err := lib.GetDBPool("default")
	if err != nil {
		t.Fatal(err)
	}
	//开始事务
	trace := lib.NewTrace()
	if _, err := lib.DBPoolLogQuery(trace, dbpool, beginSQL); err != nil {
		t.Fatal(err)
	}

	//创建表
	if _, err := lib.DBPoolLogQuery(trace, dbpool, createTableSQL); err != nil {
		lib.DBPoolLogQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}

	//插入数据
	if _, err := lib.DBPoolLogQuery(trace, dbpool, insertSQL); err != nil {
		lib.DBPoolLogQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}

	//循环查询数据
	current_id := 0
	table_name := "test1"
	fmt.Println("begin read table ", table_name, "")
	fmt.Println("------------------------------------------------------------------------")
	fmt.Printf("%6s | %6s\n", "id", "created_at")
	for {
		rows, err := lib.DBPoolLogQuery(trace, dbpool, "SELECT id,created_at FROM test1 WHERE id>? order by id asc", current_id)
		defer rows.Close()
		row_len := 0
		if err != nil {
			lib.DBPoolLogQuery(trace, dbpool, "rollback;")
			t.Fatal(err)
		}
		for rows.Next() {
			var create_time string
			if err := rows.Scan(&current_id, &create_time); err != nil {
				lib.DBPoolLogQuery(trace, dbpool, "rollback;")
				t.Fatal(err)
			}
			fmt.Printf("%6d | %6s\n", current_id, create_time)
			row_len++
		}
		if row_len == 0 {
			break
		}
	}
	fmt.Println("------------------------------------------------------------------------")
	fmt.Println("finish read table ", table_name, "")

	//删除表
	if _, err := lib.DBPoolLogQuery(trace, dbpool, dropTableSQL); err != nil {
		lib.DBPoolLogQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}

	//提交事务
	lib.DBPoolLogQuery(trace, dbpool, commitSQL)
	TearDown()
}

func Test_GORM(t *testing.T) {
	SetUp()

	//获取链接池
	dbpool, err := lib.GetGormPool("default")
	if err != nil {
		t.Fatal(err)
	}
	db := dbpool.Begin()
	traceCtx := lib.NewTrace()

	//设置trace信息
	db = db.SetCtx(traceCtx)
	if err := db.Exec(createTableSQL).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	//插入数据
	t1 := &Test1{Name: "test_name", CreatedAt: time.Now()}
	if err := db.Save(t1).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	//查询数据
	list := []Test1{}
	if err := db.Where("name=?", "test_name").Find(&list).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	fmt.Println(list)

	//删除表数据
	if err := db.Exec(dropTableSQL).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	db.Commit()
	TearDown()
}
