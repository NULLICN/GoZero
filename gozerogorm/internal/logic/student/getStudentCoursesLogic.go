// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package student

import (
	"context"

	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStudentCoursesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStudentCoursesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStudentCoursesLogic {
	return &GetStudentCoursesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStudentCoursesLogic) GetStudentCourses(req *types.StudentIdReq) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
