package main

import (
	"MyServer/app"
	_ "MyServer/router"
)

func main() {
	//fmt.Println("start ...")

	myApp := app.Application()

	myApp.Router.Run()

}
