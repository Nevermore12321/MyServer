package main

import (
	"MyServer/app"
	_ "MyServer/router"
)

func main() {
	//fmt.Println("start ...")

	myApp := app.Application()
	//  进程结束 同步 log
	defer app.Logger.Sync()
	myApp.Router.Run()

}
