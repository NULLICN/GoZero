// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"context"

	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddressLogic {
	return &AddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressLogic) Address() (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
