package middleware

import (
	"MyServer/app"
	"MyServer/config"
)

//  middleware init 用来注册 公共的 中间件，而 gin-jwt 这种中间件，需要根据具体的 路由来配置，因此在 route 中配置

func init() {
	//  将所有的中间件 在这里注册进入 gin framework， 注意，在main函数中 要 import middleware
	myApp := app.Application()

	//  注册 跨域中间件
	myApp.Router.Use(corssDomain())
	app.Logger.Debug("Corss Domain middleware success configure.")

	//  注册 session redis 中间件
	myApp.Router.Use(sessionByRedis())
	app.Logger.Debug("Session store to Redis middleware success configure.")

	//  注册 捕获异常 中间件 catchErr
	myApp.Router.Use(catchError())
	app.Logger.Debug("Catch error middleware success configure.")

	//  注册 csrf token 验证 中间件
	myApp.Router.Use(csrfTokenValidate(CsrfOptions{
		Secret:    config.GetStringFromConfig("csrf.csrfSecret"),
		ErrorFunc: AddErrorFunc,
	}))
	app.Logger.Debug("CSRF token validator middleware success configure.")
}
