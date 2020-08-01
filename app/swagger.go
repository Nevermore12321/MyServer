package app

import (
	"MyServer/config"
	"MyServer/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/zap"
)

func swaggerConfigure(engine *gin.Engine) {
	//  swagger 设置
	//  读取配置文件

	docs.SwaggerInfo.Title = config.GetStringFromConfig("swagger.title")
	docs.SwaggerInfo.Description = config.GetStringFromConfig("swagger.description")
	docs.SwaggerInfo.Version = config.GetStringFromConfig("swagger.version")
	docs.SwaggerInfo.BasePath = config.GetStringFromConfig("swagger.basePath")
	docs.SwaggerInfo.Schemes = config.GetStringSliceFromConfig("swagger.schemes")

	swaggerPath := "http://" +
		config.ConfigViper.GetString("server.host") + ":" +
		config.ConfigViper.GetString("server.port") + "/swagger/doc.json"

	url := ginSwagger.URL(swaggerPath)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	//  打印日志 ， swagger 初始化成功
	Logger.Info("swagger 配置完成",
		zap.String("Swagger Uri", "/swagger/index.html"))
}
