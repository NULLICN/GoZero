// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package focus

import (
	"context"

	"gozeroapi/internal/svc"
	"gozeroapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFocusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFocusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFocusLogic {
	return &GetFocusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFocusLogic) GetFocus() (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line
	resp = &types.CommonResponse{
		Success: true,
		Code:    200,
		Message: "two focuses",
		Data: []types.Focus{
			{
				Id:    "1",
				Name:  "f1",
				Title: "a tree",
				Link:  "https://www.zeromicro.com/",
			},
			{
				Id:    "2",
				Name:  "f2",
				Title: "a mount",
				Link:  "https://www.zeromicro.com/",
			},
		},
	}

	return
}
