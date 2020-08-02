package adminDB

import (
	"MyServer/app"
	"errors"
	"github.com/jinzhu/gorm"
)

//  用户信息表  表名：  字段有： ID （主键） ，name， password
type UserInfo struct {
	gorm.Model
	Name     string `gorm: type:varchar(100);UNIQUE;NOT NULL`
	Password string `gorm: type:varchar(100);UNIQUE;NOT NULL`
}

//  对 UserInfo 表，实现 AdminDBOperation 接口的 Insert 方法
func (user *UserInfo) Insert() error {
	if user == nil {
		return errors.New("no users input")
	}

	if AdminDB.NewRecord(user) {
		result := AdminDB.Create(&user)
		if result.Error != nil {
			return result.Error
		}
	} else {
		return errors.New("userInfo data is already existed")
	}

	return nil
}

func (user *UserInfo) Delete(hardDel bool) error {
	if user == nil {
		return errors.New("uo users input")
	}

	var res UserInfo
	//  先查找数据库是否存在， 如果不存在，跳过
	queryErr := AdminDB.Where("name = ?", user.Name).Find(&res).Error
	if queryErr != nil {
		return queryErr
	}
	//  如果 没查到数据，表示数据库中没有，跳过
	if res != (UserInfo{}) {
		//  如果找到，就删掉, 分为 软删除和 硬删除
		var deleteErr error
		if hardDel {
			deleteErr = AdminDB.Unscoped().Where("name = ?", user.Name).Delete(&user).Error
		} else {
			deleteErr = AdminDB.Where("name = ?", user.Name).Delete(&user).Error
		}
		if deleteErr != nil {
			return deleteErr
		}

	}

	return nil
}

//  例如： Update(map[string]interface{}{"name": "hello"}, "name = ?", "gsh")
func (user *UserInfo) Update(modify map[string]interface{}, queryString interface{}, keyList ...interface{}) error {
	if user == nil {
		return errors.New("uo users input")
	}
	var myUser UserInfo
	//  先查找数据库是否存在， 如果不存在，跳过
	AdminDB.Where(queryString, keyList...).Find(&myUser)
	app.Logger.Info(myUser.Name)

	updateErr := AdminDB.Model(&user).Where(queryString, keyList...).Updates(modify).Error
	if updateErr != nil {
		return updateErr
	}

	return nil
}
