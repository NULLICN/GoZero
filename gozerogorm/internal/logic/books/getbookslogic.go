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

type GetbooksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetbooksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetbooksLogic {
	return &GetbooksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetbooksLogic) Getbooks() (resp *types.CommonResponse, err error) {
	var books []gorm.Book
	err = l.svcCtx.DB.Model(&books).Find(&books).Error
	if err != nil {
		return &types.CommonResponse{
			Code:    0,
			Success: false,
			Message: err.Error(),
			Data:    nil,
		}, nil
	}
	resp = &types.CommonResponse{
		Code:    0,
		Message: "success",
		Data:    books,
		Success: true,
	}

	return
}
