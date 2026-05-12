package logic

import (
	"context"

	"client/m4api/internal/svc"
	"client/m4api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SayHelloPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSayHelloPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SayHelloPostLogic {
	return &SayHelloPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SayHelloPostLogic) SayHelloPost(req *types.SayHelloReq) (*types.SayHelloRes, error) {
	l.Infof("[gateway] incoming POST name=%q", req.Name)
	return callRpc(l.ctx, l.svcCtx.GreeterRpc, req.Name)
}
