// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package books

import (
	"context"
	"gozerogorm/model/gorm"

	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserByIdLogic) GetUserById(req *types.BookQuestById) (resp *types.CommonResponse, err error) {
	var book gorm.Book
	err = l.svcCtx.DB.Where("id=?", req.Id).First(&book).Error

	if err != nil || req.Id == 0 {
		return &types.CommonResponse{
			Code:    0,
			Success: false,
			Message: err.Error(),
			Data:    nil,
		}, nil
	}
	resp = &types.CommonResponse{
		Code:    200,
		Success: true,
		Message: "Success",
		Data:    book,
	}
	return
}
