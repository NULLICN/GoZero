// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"firstdemo/internal/svc"
	"firstdemo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FirstdemoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFirstdemoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FirstdemoLogic {
	return &FirstdemoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 处理函数被返回
func (l *FirstdemoLogic) Firstdemo(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	resp = &types.Response{
		Message: "Hello " + req.Name,
	}
	return
}
