// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package users

import (
	"context"
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

	// 1. 构建 API 类型的 User 对象（单数）
	apiUser := &types.User{
		Username: req.Username,
	}

	// 2. 使用转换函数将 API User 转换为数据库 Users 类型（复数）
	// 这里演示了 types.User <-> mysql.Users 的适配器模式
	dbUser := apiUser.ToDBModel()

	// 3. 调用数据库model层的Insert方法进行数据持久化
	// 数据库会自动设置 id 和 add_time
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

	// 4. 获取新插入的用户ID
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		resp = &types.CommonResponse{
			Success: false,
			Code:    500,
			Message: "获取用户ID失败",
		}
		return
	}

	// 5. 查询完整的用户数据（包括数据库自动生成的 add_time）
	userId := int64(lastInsertId)
	dbUserFull, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		l.Errorf("查询用户信息失败: %v", err)
		resp = &types.CommonResponse{
			Success: false,
			Code:    500,
			Message: "查询用户信息失败",
		}
		return
	}

	// 6. 使用转换函数将数据库用户对象转回API类型
	apiUser = types.UserFromDBModel(dbUserFull)

	// 7. 返回成功响应
	l.Infof("创建用户成功，用户ID: %d，创建时间: %s", apiUser.Id, apiUser.AddTime)
	resp = &types.CommonResponse{
		Success: true,
		Code:    200,
		Message: "创建用户成功",
		Data:    apiUser,
	}

	return
}
