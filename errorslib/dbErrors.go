package errorslib

//  在此 定义 error
/*
 *  1. 接收者为空， Error info : Code: 100, Type: ReceiverIsNil, Msg: The receiver of this method is nik
 *  2. 插入的数据已经存在， Error info : Code: 101, Type: InsertDataExistErr, Msg: The data to be inserted already exists
 *  3. 删除的数据不存在： Error info: Code: 102, Type: DeleteDataNotExistErr, Msg: The data to be deleted does not exist
 *	4. 更新的数据不存在： Error info: Code: 103, Type: UpdateDataNotExistErr, Msg: The data to be updated does not exist
 *	5. 输入的过滤条件为空： Error info: Code: 200, Type: NoSqlFilterErr, Msg: Filter SQL condition does not exist
 */

var (
	ErrUsernmaeNotFound  = New(100, "UsernameNotFound", "the searched user does not exist")
	ErrIncorrectPasswrod = New(101, "IncorrectPassword", "the password is incorrect")
)
