package gorm

type Lesson struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}

func (Lesson) TableName() string {
	return "lesson"
}