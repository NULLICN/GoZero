package biz

const Ok = 200

var (
	DBError                  = NewError(10000, "数据库错误")
	ParamError               = NewError(10001, "参数错误")
	DataNotExist             = NewError(10002, "数据不存在")
	ServerError              = NewError(10003, "服务器错误")
	UserNotExist             = NewError(10004, "用户不存在")
	UserNotLogin             = NewError(10005, "用户未登录")
	UserNotExistOrPasswordError = NewError(10006, "用户不存在或密码错误")
	UserExist				 = NewError(10007, "用户已存在")
)