package main

import (
	"MyServer/app"
	"MyServer/config"
	_ "MyServer/middleware"
	_ "MyServer/router"
)

func main() {
	//fmt.Println("start ...

	myApp := app.Application()
	//  进程结束 同步 log
	defer app.Logger.Sync()

	//  读取配置文件中的 IP 和 端口
	host := config.GetStringFromConfig("server.host")
	port := config.GetStringFromConfig("server.port")
	listenAddr := host + ":" + port
	_ = myApp.Router.Run(listenAddr)
}
