package errorslib

//  在此 定义 error
/*
 *  1. 接收者为空， Error info : Code: 100, Type: ReceiverIsNil, Msg: The receiver of this method is nik
 *  2. 插入的数据已经存在， Error info : Code: 101, Type: InsertDataExistErr, Msg: The data to be inserted already exists
 *  3. 删除的数据不存在： Error info: Code: 102, Type: DeleteDataNotExistErr, Msg: The data to be deleted does not exist
 *	4. 更新的数据不存在： Error info: Code: 103, Type: UpdateDataNotExistErr, Msg: The data to be updated does not exist
 *	5. 输入的过滤条件为空： Error info: Code: 200, Type: NoSqlFilterErr, Msg: Filter SQL condition does not exist
 */

//  根据 err code 生成新的 error
func NewDbErrorInt(code int) *DBError {
	var errType string
	var msg string

	switch code {
	case 100:
		errType = "ReceiverIsNilErr"
		msg = "The receiver of this method is nil"
	case 101:
		errType = "InsertDataExistErr"
		msg = "The data to be inserted already exists"
	case 102:
		errType = "DeleteDataNotExistErr"
		msg = "The data to be deleted does not exist"
	case 103:
		errType = "UpdateDataNotExistErr"
		msg = "The data to be updated does not exist"
	case 200:
		errType = "NoSqlFilterErr"
		msg = "Filter SQL condition does not exist"
	default:
		code = 0
		errType = "NoExistErr"
		msg = "The error is not existed"
	}

	return &DBError{
		Code: code,
		Type: errType,
		Msg:  msg,
	}
}

//  根据 err type 来生成新的 error
func NewDbErrorType(errType string) *DBError {
	var code int
	var msg string

	switch errType {
	case "ReceiverIsNilErr":
		code = 100
		msg = "The receiver of this method is nil"
	case "InsertDataExistErr":
		code = 101
		msg = "The data to be inserted already exists"
	case "DeleteDataNotExistErr":
		code = 102
		msg = "The data to be deleted does not exist"
	case "UpdateDataNotExistErr":
		code = 103
		msg = "The data to be updated does not exist"
	case "NoSqlFilterErr":
		code = 200
		msg = "Filter SQL condition does not exist"
	default:
		code = 0
		errType = "NoExistErr"
		msg = "The error is not existed"
	}

	return &DBError{
		Code: code,
		Type: errType,
		Msg:  msg,
	}
}
