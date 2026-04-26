package gorm

type Book struct {
	Id       int     `gorm:"primaryKey" json:"id"`
	Bookname string  `json:"bookname"`
	Price    float64 `json:"price"`
}

func (Book) TableName() string {
	return "books"
}
