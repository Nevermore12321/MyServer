package main

import (
	"MyServer/app"
	"MyServer/docs"
	"fmt"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func main() {
	fmt.Println("start ...")
	// programatically set swagger info
	docs.SwaggerInfo.Title = "Golang Web API"
	docs.SwaggerInfo.Description = "This is a my backend server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := gin.Default()
	url := ginSwagger.URL("http://localhost:1234/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	myApp := app.Application()

	myApp.Router.Run()

}
