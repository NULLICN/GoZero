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

type SQLBookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSQLBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SQLBookLogic {
	return &SQLBookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SQLBookLogic) SQLBook(req *types.BookSQL) (resp *types.CommonResponse, err error) {
	if req.Id == 0 {
		return &types.CommonResponse{
			Code:    0,
			Success: false,
			Message: "id is required",
			Data:    nil,
		}, nil
	}

	var book gorm.Book

	// 查询使用raw方法，执行原生SQL
	err = l.svcCtx.DB.Raw("SELECT * FROM books WHERE id=?", req.Id).Scan(&book).Error

	if err != nil {
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
