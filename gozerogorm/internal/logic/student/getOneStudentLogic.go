// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package student

import (
	"context"
	"fmt"

	"gozerogorm/internal/biz"
	"gozerogorm/internal/svc"
	"gozerogorm/internal/types"
	"gozerogorm/model/gorm"

	ORM "gorm.io/gorm"

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
	//err = l.svcCtx.DB.Where("id=?", req.StudentId).First(&student).Error
	err = l.svcCtx.DB.Preload("Lesson", "id!=1"/*表示不查询id为1的课程*/, 
	func(db *ORM.DB) *ORM.DB { 
		db.Order("id DESC")
		return db
	}).Find(&student, req.StudentId).Error
	if err != nil || student.Id == 0 {
		return biz.Error(biz.DataNotExist), nil
	}

	resp = biz.Success(student)

	/* resp = &types.CommonResponse{
		Code:    200,
		Success: true,
		Message: "success",
		Data:    student,
	} */
	return
}
