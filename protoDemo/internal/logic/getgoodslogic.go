package logic

import (
	"context"

	"proto_demo/goodsService"
	"proto_demo/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsLogic {
	return &GetGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGoodsLogic) GetGoods(in *goodsService.GetGoodsReq) (*goodsService.GetGoodsRes, error) {
	// todo: add your logic here and delete this line

	return &goodsService.GetGoodsRes{}, nil
}
