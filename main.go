package main

import (
	"MyServer/app"
)

func main() {
	//fmt.Println("start ...")

	myApp := app.Application()

	myApp.Router.Run()

}
