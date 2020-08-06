package router

import (
	"MyServer/app"
	"MyServer/controller"
)

func init() {
	appInstance := app.Application()

	appInstance.Router.GET("/v1", controller.LoginAction)
}
