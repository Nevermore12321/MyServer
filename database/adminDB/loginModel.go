package adminDB

import (
	"MyServer/app"
	"MyServer/errorslib"
	"github.com/jinzhu/gorm"
)

//  用户信息表  表名：  字段有： ID （主键） ，name， password
type UserInfo struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Password string `gorm:"unique"`
}

//  对 UserInfo 表，实现 AdminDBOperation 接口的 Insert 方法
func (user *UserInfo) Insert() *errorslib.DBError {
	if user == nil {
		//  100 : ReceiverIsNil
		return errorslib.NewDbErrorInt(100)
	}
	var queryUser UserInfo
	queryErr := AdminDB.Where("name = ?", user.Name).First(&queryUser).Error

	if queryErr == gorm.ErrRecordNotFound {
		app.Logger.Info("come in")
		result := AdminDB.Create(&user)
		if result.Error != nil {
			panic(result.Error)
		}
	} else if queryErr != nil {
		panic(queryErr)
	} else {
		return errorslib.NewDbErrorInt(101)
	}

	return nil
}

func (user *UserInfo) Delete(hardDel bool) *errorslib.DBError {
	if user == nil {
		return errorslib.NewDbErrorInt(100)
	}

	var res UserInfo
	//  先查找数据库是否存在， 如果不存在，跳过
	queryErr := AdminDB.Where("name = ?", user.Name).Find(&res).Error
	//  如果 没查到数据，表示数据库中没有，返回 102 错误
	if queryErr == gorm.ErrRecordNotFound {
		return errorslib.NewDbErrorInt(102)
		//  查找时遇到别的错误，直接panic
	} else if queryErr != nil {
		panic(queryErr)
		// 找到后，删除
	} else {
		//  如果找到，就删掉, 分为 软删除和 硬删除
		var deleteErr error
		//  硬删除，直接将 整个字段从磁盘删除
		if hardDel {
			deleteErr = AdminDB.Unscoped().Where("name = ?", user.Name).Delete(UserInfo{}).Error
			//  如果是软删除， 仅仅会对 delete_at 字段做修改，并且为了保证 name 的唯一性，需要修改 name 的值 添加后缀 _DELETE
		} else {
			deleteErr = AdminDB.Where("name = ?", user.Name).Delete(UserInfo{}).Error
			AdminDB.Model(user).Unscoped().Where("name = ?", user.Name).Update("name", user.Name+"_DELETE")
		}
		if deleteErr != nil {
			panic(deleteErr)
		}
	}

	return nil
}

//  这个方法，根据 接收者的 name 作为查询，找到就更新
func (user *UserInfo) Update(modify map[string]interface{}) *errorslib.DBError {
	if user == nil {
		return errorslib.NewDbErrorInt(100)
	}

	var myUser UserInfo
	//  先查找数据库是否存在， 如果不存在，跳过
	AdminDB.Where("name = ?", user.Name).Find(&myUser)
	if myUser == (UserInfo{}) {
		return errorslib.NewDbErrorInt(103)
	}

	updateErr := AdminDB.Model(user).Where("name = ?", user.Name).Update(modify).Error
	if updateErr != nil {
		panic(updateErr)
	}

	return nil
}

//  例如： Update(map[string]interface{}{"name": "hello"}, "name = ?", "gsh")
//  这个方法 跟接收者 没有关系
func (user *UserInfo) UpdateByWhere(modify map[string]interface{}, queryString interface{}, keyList ...interface{}) *errorslib.DBError {
	var myUser UserInfo
	//  先查找数据库是否存在， 如果不存在，跳过
	AdminDB.Where(queryString, keyList...).Find(&myUser)
	if myUser == (UserInfo{}) {
		return errorslib.NewDbErrorInt(103)
	}

	updateErr := AdminDB.Model(user).Where(queryString, keyList...).Updates(modify).Error
	if updateErr != nil {
		panic(updateErr)
	}

	return nil
}

//  这个方法 跟接收者 UserInfo 没有关系，根据参数查询
func (user *UserInfo) QueryAll(out *[]UserInfo, where string, args ...interface{}) *errorslib.DBError {
	if where == "" {
		return errorslib.NewDbErrorInt(200)
	}

	// 获取第一个匹配的记录
	queryErr := AdminDB.Where(where, args...).Find(out).Error
	if queryErr != nil {
		panic(queryErr)
	}
	return nil
}

func (user *UserInfo) QueryAllByName(out *[]UserInfo) *errorslib.DBError {
	// 获取第一个匹配的记录
	queryErr := AdminDB.Where("name = ?", user.Name).Find(out).Error
	if queryErr != nil {
		panic(queryErr)
	}
	return nil
}

//  这个方法 与接收者 也没有关系，通过 not 语句查找
func (user *UserInfo) QueryNot(out *[]UserInfo, not string, args ...interface{}) *errorslib.DBError {
	if not == "" {
		return errorslib.NewDbErrorInt(200)
	}

	queryErr := AdminDB.Not(not, args...).Find(out).Error
	if queryErr != nil {
		panic(queryErr)
	}
	return nil
}
