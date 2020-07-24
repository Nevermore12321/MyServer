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

	docs.SwaggerInfo.Title = config.ConfigViper.GetString("swagger.title")
	docs.SwaggerInfo.Description = config.ConfigViper.GetString("swagger.description")
	docs.SwaggerInfo.Version = config.ConfigViper.GetString("swagger.version")
	docs.SwaggerInfo.BasePath = config.ConfigViper.GetString("swagger.basePath")
	docs.SwaggerInfo.Schemes = config.ConfigViper.GetStringSlice("swagger.schemes")

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	//  打印日志 ， swagger 初始化成功
	Logger.Info("swagger 配置完成",
		zap.String("Swagger Uri", "/swagger/index.html"))
}
