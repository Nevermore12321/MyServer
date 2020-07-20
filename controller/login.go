package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 登录
// @Description 验证用户登录是否成功
// @Accept  json
// @Produce  json
// @Param   pageNum     path    int     true        "pageNum"
// @Param expName body string false "expName"
// @Param expType body string false "expType"
// @Param expTrade body string false "expTrade"
// @Param expScene body string false "expScene"
// @Param expRemark body string false "expRemark"
// @Param expDeg body int false "expDeg"
// @Param expCreateUser body string false "expCreateUser"
// @Success 200 {string} json "{"errcode":"200","data":"[{"expId":"111","expName":"TensorFlow","expType":"TensorFlow","expTrade":"制造业","expScene":"零售","expRemark":"零售零售零售零售","expDeg":"2","expCreateUser":"tfg"}]","msg":""}"
// @Failure 400 {string} json "{"errcode":"400","data":"","msg":"error......"}"
// @Router / [get]
func LoginController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
