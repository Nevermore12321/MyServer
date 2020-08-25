package adminDB

import (
	"MyServer/app"
	"MyServer/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var AdminDB *gorm.DB

func init() {
	mysqlInit()
}

func mysqlInit() {
	dbHost := config.GetStringFromConfig("mysql.host")
	dbUser := config.GetStringFromConfig("mysql.user")
	dbPassword := config.GetStringFromConfig("mysql.password")
	adminDatabase := config.GetStringFromConfig("mysql.admin.database")

	//  mysql 连接地址  user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	mysqlUrl := dbUser + ":" + dbPassword + "@(" + dbHost + ")/" + adminDatabase + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open("mysql", mysqlUrl)
	if err != nil {
		errMsg := fmt.Sprintf("Mysql Err: %v", err)
		app.Logger.Error(errMsg)
		panic(err)
	}

	//  创建 mysql 连接池
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	//  建表
	db.AutoMigrate(&UserInfo{})
	db.AutoMigrate(&OriginAuth{})

	AdminDB = db
}
