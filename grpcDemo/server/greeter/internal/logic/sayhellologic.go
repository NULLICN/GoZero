package logic

import (
	"context"

	"greeter/greeter"
	"greeter/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SayHelloLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSayHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SayHelloLogic {
	return &SayHelloLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SayHelloLogic) SayHello(in *greeter.HelloReq) (*greeter.HelloRes, error) {
	// todo: add your logic here and delete this line

	return &greeter.HelloRes{
		Message: "你好" + in.Name,
	}, nil
}
