// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package lesson

import (
	"context"

	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLessonStudentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLessonStudentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLessonStudentsLogic {
	return &GetLessonStudentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLessonStudentsLogic) GetLessonStudents(req *types.LessonIdReq) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
