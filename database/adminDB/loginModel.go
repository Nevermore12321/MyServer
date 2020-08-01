package adminDB

import (
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
func (user *UserInfo) Insert(users []*UserInfo) error {
	if len(users) == 0 {
		return errors.New("No users input.")
	}

	for _, user := range users {
		if AdminDB.NewRecord(user) {
			result := AdminDB.Create(&user)
			if result.Error != nil {
				return result.Error
			}
		} else {
			return errors.New("UserInfo data is already existed.")
		}
	}
	return nil
}

func (user *UserInfo) Delete(users []*UserInfo) error {
	return nil
}

func (user *UserInfo) Update(users []*UserInfo) error {
	return nil
}

func (user *UserInfo) Query(users []*UserInfo) error {
	return nil
}
