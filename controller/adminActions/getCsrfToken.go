package adminActions

import (
	middleware "MyServer/middleware"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// @Summary 获取 CSRF Token
// @Description 用户第一次访问后端服务器时，需要获取 CSRF Token
// @Tags 校验
// @Success 200 {string} json "{"status":200,"message":"CSRF token is in response header"}"
// @Failure 400 {string} json "{"status":400,"message":"CSRF token mismatch"}"
// @Router /getCSRF [get]
func GetCsrfTokenAction(ctx *gin.Context) {
	token := middleware.GetToken(ctx)
	ctx.Writer.Header().Set("X-Csrf-Token", token)
	session := sessions.Default(ctx)
	fmt.Println(session.Get("csrfSalt"))
	ctx.JSON(200, gin.H{
		"Message": "CSRF token is in response header",
	})
}
