package middleware

import (
	"MyServer/app"
)



func init() {
	//  将所有的中间件 在这里注册进入 gin framework， 注意，在main函数中 要 import middleware
	myApp := app.Application()

	//  注册 跨域中间件
	myApp.Router.Use(corssDomain())
}