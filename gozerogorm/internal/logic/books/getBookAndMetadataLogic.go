// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package books

import (
	"context"

	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"
	"gozerogorm/model/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBookAndMetadataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBookAndMetadataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBookAndMetadataLogic {
	return &GetBookAndMetadataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBookAndMetadataLogic) GetBookAndMetadata() (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line
	var book []gorm.Book
	err = l.svcCtx.DB.Preload("Metadata").Find(&book).Error
	if err != nil {
		return &types.CommonResponse{
			Code:    0,
			Success: false,
			Message: "查询失败",
		}, nil
	}

	resp = &types.CommonResponse{
		Code:    200,
		Success: true,
		Message: "success",
		Data:    book,
	}
	return
}
