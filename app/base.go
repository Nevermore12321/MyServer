package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

//  小写，只供该 app 模块使用
//  app 类型结构体，存放 用于启动的 gin 引擎
type app struct {
	Router *gin.Engine
}

//  创建变量app， 其中 Router 初始化为 nil
var webAppInstance *app

//  Once.Do, 只执行一次
var once sync.Once

func Application() *app {
	//  初始化引擎
	once.Do(func() {
		// todo
		router := gin.Default()

		//   404 not found configure
		router.NoRoute(func(ctx *gin.Context) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Not Found",
			})
		})

		webAppInstance = &app{router}
	})

	return webAppInstance
}
