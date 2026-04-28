package biz

//定义一个ErrorUtil的结构体
type ErrorUtil struct {
	Code int
	Msg  string
}

func NewError(code int, msg string) *ErrorUtil {
	return &ErrorUtil{Code: code, Msg: msg}
}