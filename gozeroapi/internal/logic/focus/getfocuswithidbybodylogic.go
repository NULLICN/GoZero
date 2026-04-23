// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package focus

import (
	"context"

	"gozeroapi/internal/svc"
	"gozeroapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFocusWithIdByBodyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFocusWithIdByBodyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFocusWithIdByBodyLogic {
	return &GetFocusWithIdByBodyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFocusWithIdByBodyLogic) GetFocusWithIdByBody(req *types.FocusRequestByBody) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line
	logx.Info("GetFocusById Path", req.Id)
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
