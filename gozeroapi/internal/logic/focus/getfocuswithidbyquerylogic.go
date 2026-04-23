// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package focus

import (
	"context"

	"gozeroapi/internal/svc"
	"gozeroapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFocusWithIdByQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFocusWithIdByQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFocusWithIdByQueryLogic {
	return &GetFocusWithIdByQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFocusWithIdByQueryLogic) GetFocusWithIdByQuery(req *types.FocusRequestByQuery) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line
	logx.Info("GetFocusById Query", req.Id)
	resp = &types.CommonResponse{
		Success: true,
		Code:    200,
		Message: "a focus",
		Data: []types.Focus{
			{
				Id:    req.Id,
				Name:  "f1",
				Title: "a glass of cup",
				Link:  "https://www.zeromicro.com/",
			},
		},
	}
	return
}
