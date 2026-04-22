// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"gozeroapi/internal/svc"
	"gozeroapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFocusByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFocusByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFocusByIdLogic {
	return &GetFocusByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFocusByIdLogic) GetFocusById(req *types.FocusRequest) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
