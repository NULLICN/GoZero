// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package users

import (
	"context"
	"gozeroapi/internal/svc"
	"gozeroapi/internal/types"
	"strconv"

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

	// 将字符串ID转换为int64
	aUserId, parseErr := strconv.ParseInt(req.Id, 10, 64)
	if parseErr != nil {
		l.Errorf("ID格式错误: %s", req.Id)
		resp = &types.CommonResponse{
			Success: false,
			Code:    400,
			Message: "ID格式错误",
		}
		return
	}

	// 从数据库查询 mysql.Users 类型的数据
	dbUser, queryErr := l.svcCtx.UserModel.FindOne(l.ctx, aUserId)
	if queryErr != nil {
		l.Errorf("查询用户失败，ID: %s, 错误: %v", req.Id, queryErr)
		resp = &types.CommonResponse{
			Success: false,
			Code:    500,
			Message: "查询用户失败",
		}
		return
	}

	// 使用转换函数将 mysql.Users 转换为 types.User
	// 这是从数据库类型到API类型的适配器
	apiUser := types.UserFromDBModel(dbUser)

	resp = &types.CommonResponse{
		Success: true,
		Code:    200,
		Message: "查询用户成功",
		Data:    apiUser,
	}
	return
}
