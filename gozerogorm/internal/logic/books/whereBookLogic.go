// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package books

import (
	"context"
	"strings"

	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"
	"gozerogorm/model/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type WhereBookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWhereBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WhereBookLogic {
	return &WhereBookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 处理特殊字符，防止SQL注入
func escapeLike(str string) string {
	str = strings.ReplaceAll(str, "%", "\\%")
	str = strings.ReplaceAll(str, "_", "\\_")
	return str
}

func (l *WhereBookLogic) WhereBook(req *types.BookWhere) (resp *types.CommonResponse, err error) {
	var books []gorm.Book

	// 使用转义后的值
	bookname := escapeLike(req.Bookname)
	err = l.svcCtx.DB.Where("bookname LIKE ? ESCAPE '\\\\'", "%"+bookname+"%").Find(&books).Error
	// err = l.svcCtx.DB.Where("bookname like ?", "%"+req.Bookname+"%").Find(&books).Error
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
		Data:    books,
	}
	return
}
