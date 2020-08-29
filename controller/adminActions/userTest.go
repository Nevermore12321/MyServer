package adminActions

import "github.com/gin-gonic/gin"

func UserGet(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "aaaaa",
	})
}
