// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package bookmetadata

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
	// var bookmetadata []gorm.Bookmetadata
	//err = l.svcCtx.DB.Preload("Book").Find(&bookmetadata).Error
	//err = l.svcCtx.DB.Preload("Book").Find(&bookmetadata).Error
	var book []gorm.Bookmetadata
	err = l.svcCtx.DB.Preload("Book").Find(&book).Error
	if err != nil {
		return &types.CommonResponse{
			Code:    0,
			Success: false,
			Message: "联表查询失败",
		}, nil
	}

	resp = &types.CommonResponse{
		Code:    200,
		Success: true,
		Message: "success",
		//Data:    bookmetadata,
		Data: book,
	}
	return
}
