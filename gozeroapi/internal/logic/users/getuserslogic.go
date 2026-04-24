// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package users

import (
	"context"

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

	// 使用 sqlx 查询所有用户（mysql.Users 类型）
	dbUsers, err := l.svcCtx.UserModel.FindAll(l.ctx)
	if err != nil {
		l.Errorf("查询用户列表失败: %v", err)
		resp = &types.CommonResponse{
			Success: false,
			Code:    500,
			Message: "查询用户列表失败",
		}
		return
	}

	// 转换为 API 响应类型 (types.User)
	apiUsers := make([]*types.User, 0, len(dbUsers))
	for _, dbUser := range dbUsers {
		// 把dbUsers里的每一个User都进行处理，并append到切片变量中
		apiUser := types.UserFromDBModel(dbUser)
		apiUsers = append(apiUsers, apiUser)
	}

	// 返回成功响应
	l.Infof("查询用户列表成功，共 %d 条记录", len(apiUsers))
	resp = &types.CommonResponse{
		Success: true,
		Code:    200,
		Message: "查询用户列表成功",
		Data:    apiUsers, // 最后返回处理好的User切片
	}

	return
}
