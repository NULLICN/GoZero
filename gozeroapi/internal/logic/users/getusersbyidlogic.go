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

type GetUsersByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsersByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersByIdLogic {
	return &GetUsersByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsersByIdLogic) GetUsersById(req *types.UserQuestById) (resp *types.CommonResponse, err error) {
	l.Infof("开始查询用户，ID: %s", req.Id)

	var user types.User
	result := model.DB.Where("id = ?", req.Id).First(&user)

	// 检查是否查询到用户
	if result.Error != nil {
		l.Infof("用户不存在，ID: %s", req.Id)
		resp = &types.CommonResponse{
			Success: false,
			Code:    404,
			Message: "用户不存在",
		}
		return
	}

	// 返回成功响应
	l.Infof("查询用户成功，ID: %s", req.Id)
	resp = &types.CommonResponse{
		Success: true,
		Code:    200,
		Message: "查询用户成功",
		Data:    user,
	}

	return
}
