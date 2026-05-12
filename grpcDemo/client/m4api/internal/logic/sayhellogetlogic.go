package logic

import (
	"context"

	"client/m4api/internal/svc"
	"client/m4api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SayHelloGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSayHelloGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SayHelloGetLogic {
	return &SayHelloGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SayHelloGetLogic) SayHelloGet(req *types.SayHelloReq) (*types.SayHelloRes, error) {
	l.Infof("[gateway] incoming GET name=%q", req.Name)
	return callRpc(l.ctx, l.svcCtx.GreeterRpc, req.Name)
}
