// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package users

import (
	"context"
	"gozeroapi/model"

	"gozeroapi/internal/svc"
	"gozeroapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersLogic {
	return &GetUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsersLogic) GetUsers() (resp *types.CommonResponse, err error) {
	l.Infof("开始查询所有用户")

	// 查询所有用户
	var users []types.User
	result := model.DB.Find(&users)

	// 检查查询是否出错
	if result.Error != nil {
		l.Errorf("查询用户列表失败: %v", result.Error)
		resp = &types.CommonResponse{
			Success: false,
			Code:    500,
			Message: "查询用户列表失败",
			Data:    result.Error.Error(),
		}
		return
	}

	// 返回成功响应
	l.Infof("查询用户列表成功，共 %d 条记录", len(users))
	resp = &types.CommonResponse{
		Success: true,
		Code:    200,
		Message: "查询用户列表成功",
		Data:    users,
	}

	return
}
