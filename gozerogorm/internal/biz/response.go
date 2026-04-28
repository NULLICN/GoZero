package biz

import "gozerogorm/internal/types"

//成功的函数
func Success(data interface{}) *types.CommonResponse {
	return &types.CommonResponse{
		Code:    Ok,
		Message: "success",
		Success: true,
		Data:    data,
	}
}

//失败的函数
func Error(err *ErrorUtil) *types.CommonResponse {
	return &types.CommonResponse{
		Code:    err.Code,
		Message: err.Msg,
		Success: false,
	}
}