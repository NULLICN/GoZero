// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package users

import (
	"context"
	"gozeroapi/model"
	"time"

	"gozeroapi/internal/svc"
	"gozeroapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserLogic {
	return &AddUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddUserLogic) AddUser(req *types.UserAdd) (resp *types.CommonResponse, err error) {
	l.Infof("开始创建用户，Name: %s", req.User.Name)

	// 构建用户对象
	user := types.User{
		Id:      req.User.Id,
		Name:    req.User.Name,
		AddTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 执行数据库插入
	result := model.DB.Create(&user)

	if result.Error != nil {
		l.Errorf("创建用户失败: %v", result.Error)
		resp = &types.CommonResponse{
			Success: false,
			Code:    500,
			Message: "创建用户失败",
			Data:    result.Error.Error(),
		}
		return
	}

	// 返回成功响应
	l.Infof("创建用户成功，用户ID: %s", user.Id)
	resp = &types.CommonResponse{
		Success: true,
		Code:    200,
		Message: "创建用户成功",
		Data:    user,
	}

	return
}
