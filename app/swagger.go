package app

import (
	"MyServer/docs"
	_ "MyServer/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func swaggerConfigure(engine *gin.Engine) {
	//  swagger 设置
	docs.SwaggerInfo.Title = "Golang Web API"
	docs.SwaggerInfo.Description = "This is a my backend server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
