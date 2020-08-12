package adminDB

import (
	"MyServer/errorslib"
)

type AdminDBOperation interface {
	Insert() *errorslib.DBError
	Delete(hardDel bool) *errorslib.DBError
	Update(modify map[string]interface{}) *errorslib.DBError
	UpdateByWhere(modify map[string]interface{}, queryString interface{}, keyList ...interface{}) *errorslib.DBError
	QueryAll(out *[]UserInfo, where string, args ...interface{}) *errorslib.DBError
	QueryAllByName(out *[]UserInfo) *errorslib.DBError
	QueryNot(out *[]UserInfo, not string, args ...interface{}) *errorslib.DBError
}
