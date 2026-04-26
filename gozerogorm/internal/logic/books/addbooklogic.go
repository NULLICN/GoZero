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

type AddBookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBookLogic {
	return &AddBookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddBookLogic) AddBook(req *types.BookAdd) (resp *types.CommonResponse, err error) {
	var book = gorm.Book{
		Bookname: req.Bookname,
		Price:    req.Price,
	}

	err = l.svcCtx.DB.Create(&book).Error
	if err != nil {
		return nil, err
	}

	resp = &types.CommonResponse{
		Code:    200,
		Message: "success",
		Success: true,
		Data:    book,
	}
	return
}
