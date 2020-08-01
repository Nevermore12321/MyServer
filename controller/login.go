package controller

import (
	"MyServer/app"
	"MyServer/database/adminDB"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// @Summary 登录
// @Description 验证用户登录是否成功
// @Tags 用户信息
// @Accept  json
// @Produce  json
// @Param id path int true "ID"    //url参数：（name；参数类型[query(?id=),path(/123)]；数据类型；required；参数描述）
// @Success 200 {string} json "{"errcode":"200","data":"[{"expId":"111","expName":"TensorFlow","expType":"TensorFlow","expTrade":"制造业","expScene":"零售","expRemark":"零售零售零售零售","expDeg":"2","expCreateUser":"tfg"}]","msg":""}"
// @Failure 400 {string} json "{"errcode":"400","data":"","msg":"error......"}"
// @Router / [get]
func LoginController(ctx *gin.Context) {
	app.Logger.Info("log success", zap.String("url", ctx.Request.URL.Host))
	app.Logger.Debug("debug success")
	app.Logger.Error("Error success")
	app.Logger.Warn("Warn success")
	var user *adminDB.UserInfo = &adminDB.UserInfo{Name: "zy", Password: "123"}
	insertUsers := []*adminDB.UserInfo{}
	insertUsers = append(insertUsers, user)
	var loginInterface adminDB.AdminDBOperation

	loginInterface = user
	err := loginInterface.Insert(insertUsers)
	if err != nil {
		errMsg := fmt.Sprintf("Mysql Err: %v", err)
		app.Logger.Error(errMsg)
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
