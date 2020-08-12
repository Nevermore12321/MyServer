package adminActions

import (
	"MyServer/app"
	"MyServer/database/adminDB"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// @Summary 登录
// @Description 验证用户登录是否成功
// @Tags 用户信息
// @Security csrf-token
// @param csrf-token header string true "X-CSRF-TOKEN"
// @Accept  json
// @Produce  json
// @Param id path int true "ID"    //url参数：（name；参数类型[query(?id=),path(/123)]；数据类型；required；参数描述）
// @Success 200 {string} json "{"errcode":"200","data":"[{"expId":"111","expName":"TensorFlow","expType":"TensorFlow","expTrade":"制造业","expScene":"零售","expRemark":"零售零售零售零售","expDeg":"2","expCreateUser":"tfg"}]","msg":""}"
// @Failure 400 {string} json "{"errcode":"400","data":"","msg":"error......"}"
// @Router /login [get]
func LoginAction(ctx *gin.Context) {
	app.Logger.Info("log success", zap.String("url", ctx.Request.URL.Host))
	app.Logger.Debug("debug success")
	app.Logger.Error("Error success")
	app.Logger.Warn("Warn success")

	session := sessions.Default(ctx)

	//var count int

	v := session.Get("count")
	fmt.Println(v)
	fmt.Printf("type : %T\n", v)
	fmt.Println(ctx.HandlerName())
	session.Set("count", "gsh")
	_ = session.Save()
	//
	//if v == nil {
	//	count = 0
	//} else {
	//	count = v.(int)
	//	count += 1
	//}
	//session.Set("count", "test")
	////session.Set("count", 0)
	//err := session.Save()
	//if err != nil {
	//	errMsg := fmt.Sprintf("Session Err: %v", err)
	//		app.Logger.Error(errMsg)
	//}
	//
	type User = adminDB.UserInfo

	var user *User = &User{Name: "gsh", Password: "123"}

	var loginInterface adminDB.AdminDBOperation
	//
	loginInterface = user
	//  增

	//var out []*adminDB.UserInfo
	//err = loginInterface.QueryAllByName(out)
	//fmt.Println(out)
	//
	err := loginInterface.Delete(false)
	if err != nil && err.Code == 102 {
		app.Logger.Error("Mysql Err: ", zap.Int("Code", err.Code), zap.String("errType", err.Type), zap.String("Msg", err.Msg))
		//panic(err)
	}

	//err := loginInterface.Update(map[string]interface{}{"name": "hello", "password": "1234"}, "name = ?", user.Name)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
