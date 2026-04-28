// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package lesson

import (
	"context"

	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"
	"gozerogorm/model/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLessonListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLessonListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLessonListLogic {
	return &GetLessonListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLessonListLogic) GetLessonList() (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line
	var lessons []gorm.Lesson
	err = l.svcCtx.DB.Find(&lessons).Error
	if err != nil {
		return &types.CommonResponse{
			Code:    0,
			Success: false,
			Message: "查询课程列表失败",
		}, nil
	}

	resp = &types.CommonResponse{
		Code:    200,
		Success: true,
		Message: "success",
		Data:    lessons,
	}
	return
}
