package router

import (
	"MyServer/app"
	"MyServer/controller/adminActions"
)

func init() {
	appInstance := app.Application()
	v1 := appInstance.Router.Group("v1")
	{
		v1.GET("getCSRF", adminActions.GetCsrfTokenAction)
		admin := v1.Group("admin")
		{
			admin.GET("login", adminActions.LoginAction)
		}
	}
}
