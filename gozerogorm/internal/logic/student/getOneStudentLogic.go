// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package student

import (
	"context"
	"fmt"

	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"
	"gozerogorm/model/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOneStudentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOneStudentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOneStudentLogic {
	return &GetOneStudentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOneStudentLogic) GetOneStudent(req *types.StudentIdReq) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line
	var student gorm.Student
	fmt.Printf("================req.studentId: %v\n", req.StudentId)
	err = l.svcCtx.DB.Where("id=?", req.StudentId).First(&student).Error
	if err != nil {
		return &types.CommonResponse{
			Code:    0,
			Success: false,
			Message: "查询学生详情失败",
		}, nil
	}

	resp = &types.CommonResponse{
		Code:    200,
		Success: true,
		Message: "success",
		Data:    student,
	}
	return
}
