package logic

import (
	"context"

	"proto_demo/goodsService"
	"proto_demo/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGoodsLogic {
	return &AddGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddGoodsLogic) AddGoods(in *goodsService.AddGoodsReq) (*goodsService.AddGoodsRes, error) {
	// todo: add your logic here and delete this line

	return &goodsService.AddGoodsRes{}, nil
}
