package middleware

import (
	"MyServer/app"
	"fmt"
	"github.com/gin-gonic/gin"
)

func catchError() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//   如果 有异常 panic ， 通过 defer 捕获
		defer func() {
			if err := recover(); err != nil {
				errMessage := fmt.Sprintf("Catch Err: %v", err)
				app.Logger.Panic(errMessage)
			}
		}()

		ctx.Next()
	}

}
