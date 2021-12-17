package main

import (
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"log"
	"time"
)

func main() {
	if err := lib.InitModule("./conf/dev/",[]string{"base","mysql","redis",}); err != nil {
		log.Fatal(err)
	}
	defer lib.Destroy()

	type IdentifyInfo struct {
		Id               int       `json:"id" gorm:"primary_key" description:"自增主键"`
	}
	idf := &IdentifyInfo{}
	//var total int64
	err := lib.GORMDefaultPool.Table("agv_identity_info").Find(idf).Error
	if err != nil {
		//fmt.Println(err)
	}
	fmt.Println(idf)
	//todo sth
	lib.Log.TagInfo(lib.NewTrace(), lib.DLTagUndefind, map[string]interface{}{"message": "todo sth"})
	time.Sleep(time.Second)
}
