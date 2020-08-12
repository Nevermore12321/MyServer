package errorslib

import "encoding/json"

//  自定义  数据库 的 Error 结构体
type DBError struct {
	Code int
	Type string
	Msg  string
}

//  结构体类型 实现 Error 方法，以满足 error 接口
func (e *DBError) Error() string {
	errStr, parseErr := json.Marshal(e)
	if parseErr != nil {
		panic(parseErr)
	}
	return string(errStr)
}

//  直接传入 三个参数 生成新的 error， 自定义 error
func New(code int, errType string, msg string) *DBError {
	return &DBError{
		Code: code,
		Type: errType,
		Msg:  msg,
	}
}
