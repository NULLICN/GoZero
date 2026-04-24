// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package users

import (
	"context"
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
	l.Infof("开始创建用户，Name: %s", req.Username)

	// 获取当前时间（字符串格式）
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// 1. 构建 API 类型的 User 对象（单数）
	apiUser := &types.User{
		Id:       0,
		Username: req.Username,
		AddTime:  currentTime,
	}

	// 2. 使用转换函数将 API User 转换为数据库 Users 类型（复数）
	// 这里演示了 types.User <-> mysql.Users 的适配器模式
	dbUser := apiUser.ToDBModel()

	// 3. 调用数据库model层的Insert方法进行数据持久化
	result, err := l.svcCtx.UserModel.Insert(l.ctx, dbUser)
	if err != nil {
		l.Errorf("创建用户失败: %v", err)
		resp = &types.CommonResponse{
			Success: false,
			Code:    500,
			Message: "创建用户失败",
		}
		return
	}

	// 4. 获取新插入的用户ID（可选，取决于是否需要返回）
	lastInsertId, err := result.LastInsertId()
	if err == nil {
		apiUser.Id = int(lastInsertId)
	}

	// 5. 返回成功响应
	l.Infof("创建用户成功，用户ID: %d", apiUser.Id)
	resp = &types.CommonResponse{
		Success: true,
		Code:    200,
		Message: "创建用户成功",
		Data:    apiUser,
	}

	return
}
