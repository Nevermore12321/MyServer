package app

import (
	"MyServer/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

//   全局可以使用的  logger
var Logger *zap.Logger

func Application() *app {
	//  初始化引擎
	once.Do(func() {

		//  初始化 Logger，并将初始化好的 赋值给全局变量 Logger
		Logger = loggerConfigure()

		//  设置 gin 是 调试模式 还是 生产模式
		ginModeConfig := config.ConfigViper.GetString("server.mode")
		gin.SetMode(ginModeConfig)
		//gin.SetMode(gin.DebugMode)

		router := gin.Default()

		//  配置 swagger
		swaggerConfigure(router)
		//   404 not found configure
		router.NoRoute(func(ctx *gin.Context) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Not Found",
			})
		})																																																																																																																																																																																																																																																																																																														
																																	
		//  将配置好的  gin 的 实例 复制给 webAppInstance
		webAppInstance = &app{router}

	})

	return webAppInstance
}
