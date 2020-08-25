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
			admin.POST("login", adminActions.LoginAction)
			auth := admin.Group("auth")
			{
				auth.GET("/refresh_token", adminActions.JwtAuth.RefreshTokenHandler)
			}
			auth.Use(adminActions.JwtAuth.MiddlewareUseFunc())
			{

			}
		}
	}
}
