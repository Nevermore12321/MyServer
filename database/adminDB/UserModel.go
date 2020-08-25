package adminDB

import (
	"MyServer/errorslib"
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

//  创建 用户信息表
//  userid， 用户的唯一 id
//  usernmae， 用户名
//
type UserInfo struct {
	gorm.Model
	UserID   string    `gorm: "column:userid;type:varchar(100);unique;not null;index:UserID"`
	Username string    `gorm: "column:username;type:varchar(250);unique;not null;"`
	Birthday time.Time `gorm: "column:birthday;type:DATE;"`
	Email    string    `gorm: "column:email;type:varchar(50);unique;not null;"`
	Avatar   string    `gorm: "column:avatar;type:varchar(200);	"`
	Role     string    `gorm: "column:role;type:ENUM("admin", "user");not null;"`
	AuthType string    `gorm: "column:auth_type;type:ENUM("origin", "qq", "weibo");not null;"`
}

//  创建 用户对应的 登录种类，可以有 oauth 登录
type OriginAuth struct {
	gorm.Model
	UserID   string `gorm: "column:userid;unique;not null;"`
	Password string `gorm: "column:password;unique;not null"`
}

//  检查用户名和密码是否配对
//  如果 用户名不存在，返回错误
//  如果 用户名和密码 不匹配，返回错误
//  如果 用户名密码 匹配， 返回 true，nil
func CheckAuth(username, password string) (bool, *errorslib.DBError) {
	var user UserInfo

	//  先查找 userinfo 表中，如果 查到了 user ，那么取出 userid，在查找 password
	queryErr := AdminDB.Table("user_infos").Where("username = ?", username).First(&user).Error
	//  如果 没有找到 该用户的信息，返回用户 不存在
	if queryErr == gorm.ErrRecordNotFound {
		return false, errorslib.ErrUsernmaeNotFound
	} else if queryErr != nil {
		panic(queryErr)
	}

	//  查找出 userinfo 后，根据 验证的类型，是 传统的密码验证还是 oauth
	switch user.AuthType {
	case "origin":
		var userAuth OriginAuth
		queryErr = AdminDB.Table("origin_auths").Where("user_id = ?", user.UserID).First(&userAuth).Error
		if queryErr != nil {
			panic(queryErr)
		}

		//  查出 password 后，进行比对
		if userAuth.Password != password {
			return false, errorslib.ErrIncorrectPasswrod
		}
		return true, nil
	case "qq":
		return true, nil
	case "weibo":
		return true, nil
	default:
		panic(errors.New("auth type is not defined"))
	}

}

//  根据 username  获取 对应的字段
//  filedName  可以是  UserID， Username， Birthday, Email, Avatar, Role, AuthType
func GetUserInfo(username string) (*UserInfo, *errorslib.DBError) {
	var user UserInfo

	//  先查找 userinfo 表中，如果 查到了 user ，那么取出 userid，在查找 password
	queryErr := AdminDB.Table("user_infos").Where("username = ?", username).First(&user).Error
	//  如果 没有找到 该用户的信息，返回用户 不存在
	if queryErr == gorm.ErrRecordNotFound {
		return nil, errorslib.ErrUsernmaeNotFound
	} else if queryErr != nil {
		panic(queryErr)
	}

	return &user, nil
}
