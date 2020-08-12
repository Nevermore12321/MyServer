package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 设置跨域 的 gin 中间件
func corssDomain() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//  获取 请求的 method 和 Headers 中的 Origin
		method := ctx.Request.Method
		requestOrigin := ctx.Request.Header.Get("Origin")
		//  判断 请求的 来源 ， 请求 Headers 中的 Origin 表示请求来源
		//  如果有来源，则回response， 否则视为 不安全请求，不给予 跨域认证
		if requestOrigin != "" {
			//  取到 客户端 request 的Origin，也就是 客户端的地址，并且设置为 Access-Control-Allow-Origin 头部
			//  ctx.Header() ctx.Writer.Header().Set() 作用一一样
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", requestOrigin)

			//  Access-Control-Allow-Method 表示跨域请求允许的方法
			ctx.Header("Access-Control-Allow-Method", "POST, GET, PUT, DELETE, UPDATE, OPTIONS")

			//  Access-Control-Allow-Headers 表示跨域请求的Header中允许带的字段
			ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-Csrf-Token, X-Xsrf-Token,Token, Session")

			//  Access-Control-Expose-Headers  表示客户端（浏览器） 可以解析出来的头部Header
			ctx.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Access-Control-Allow-Origin, Content-Length, X-Csrf-Token")

			//  Access-Control-Allow-Credentials   表示允许客户端传递验证信息（COokies）
			ctx.Header("Access-Control-Allow-Credentials", "true")

			//  Access-Control-Max-Age  表示 设置缓存时间,单位秒 (24小时)
			ctx.Header("Access-Control-Max-Age", "86400")
		}

		//  需要允许 OPTIONS 的预检 请求, 如果是 OPTIONS请求， 直接返回 200 OK
		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, "OK")
		}

		//   如果 有异常 panic ， 通过 defer 捕获
		//defer func() {
		//	if err := recover(); err != nil {
		//		errMessage := fmt.Sprintf("Corss Domain Panic Info : %v", err)
		//		app.Logger.Panic(errMessage)
		//	}
		//}()
		ctx.Next()
	}
}
