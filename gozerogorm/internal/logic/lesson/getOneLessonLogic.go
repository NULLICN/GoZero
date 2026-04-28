// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package lesson

import (
	"context"
	"fmt"

	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"
	"gozerogorm/model/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOneLessonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOneLessonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOneLessonLogic {
	return &GetOneLessonLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOneLessonLogic) GetOneLesson(req *types.LessonIdReq) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line
	var lesson gorm.Lesson
	fmt.Printf("================req.id: %v\n", req.LessonId)
	//err = l.svcCtx.DB.Where("id=?", req.LessonId).First(&lesson).Error
	err = l.svcCtx.DB.Preload("Student", "").Find(&lesson, req.LessonId).Error
	if err != nil {
		return &types.CommonResponse{
			Code:    0,
			Success: false,
			Message: "查询课程及学生详情失败",
		}, nil
	}

	resp = &types.CommonResponse{
		Code:    200,
		Success: true,
		Message: "success",
		Data:    lesson,
	}
	return
}
