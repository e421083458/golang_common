package main

import (
	"github.com/e421083458/golang_common/lib"
	"log"
	"time"
)

func main(){
	if err:=lib.Init("./conf/dev/");err!=nil{
		log.Fatal(err)
	}
	defer lib.Destroy()

	//todo sth
	lib.Log.TagInfo(lib.NewTrace(), lib.DLTagUndefind, map[string]interface{}{"message": "todo sth"})
	time.Sleep(time.Second)
}