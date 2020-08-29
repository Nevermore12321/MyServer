package adminActions

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 刷新 jwt token
// @Description 更新login后获取的jwt token
// @Tags 校验
// @Security csrf-token
// @Security authorization
// @param X-Csrf-Token header string true "X-CSRF-TOKEN"
// @param Authorization header string true "Bearer xxx"
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"status":200,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJSb2xlIjoiYWEiLCJVc2VySUQiOiJnc2hhYWEiLCJVc2VybmFtZSI6ImdzaCIsImV4cCI6MTU5ODQ1MzAyMywib3JpZ19pYXQiOjE1OTg0NDk0MjMsInRlc3RLZXkiOiIifQ.lY7eeTtIyO9eexpmEAWh8s196MGiFpJR-xjiFgdlRLA","expire":"2020-08-26T22:43:43+08:00"}"
// @Failure 401 {string} json "{"status":401,"data":{},"message":"token is expired"}"
// @Router /admin/refresh_token [get]
func RefreshTokenHandler(ctx *gin.Context) {
	//  更新 操作，返回 token string， 过期时间，错误
	tokenStr, expire, err := JwtAuth.RefreshToken(ctx)
	if err != nil {
		JwtAuth.Unauthorized(ctx, http.StatusUnauthorized, JwtAuth.HTTPStatusMessageFunc(err, ctx))
		return
	}

	JwtAuth.RefreshResponse(ctx, http.StatusOK, tokenStr, expire)
}
