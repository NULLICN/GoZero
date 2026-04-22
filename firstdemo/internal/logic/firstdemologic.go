// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"fmt"

	"firstdemo/internal/svc"
	"firstdemo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FirstdemoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFirstdemoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FirstdemoLogic {
	return &FirstdemoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 处理函数被返回
func (l *FirstdemoLogic) Firstdemo(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	fmt.Println("DataSource in logic:", l.svcCtx.Config.Mysql.DataSource) // 来自于更远的注入的配置
	fmt.Println("DataSource in logic:", l.svcCtx.DataSource)              // 来自于更近的注入的配置
	resp = &types.Response{
		Message: "Hello " + req.Name,
	}
	return
}
