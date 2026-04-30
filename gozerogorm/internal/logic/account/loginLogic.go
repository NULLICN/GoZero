// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package account

import (
	"context"
	"time"

	"gozerogorm/internal/biz"
	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line
	token, _ := biz.GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, req.Username)

    return biz.Success(token), nil
	
}
