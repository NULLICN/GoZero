package gorm

type Lesson struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Student []Student `json:"student" gorm:"many2many:lesson_student;"`
}

func (Lesson) TableName() string {
	return "lesson"
}