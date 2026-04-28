package gorm

type Student struct {
	Id   int    `json:"id"`
	Number int `json:"number"`
	Password string `json:"password"`
	ClassId int `json:"class_id"`
	Name string `json:"name"`
	Lesson []Lesson `json:"lesson" gorm:"many2many:user_languages;"`
}

func (Student) TableName() string {
	return "student"
}