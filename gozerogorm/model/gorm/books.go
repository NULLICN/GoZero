package gorm

type Book struct {
	Id             int     `gorm:"primaryKey" json:"id"`
	Bookname       string  `json:"bookname"`
	Price          float64 `json:"price"`
	BookmetadataId int     `json:"bookmetadata_id"`
}

func (Book) TableName() string {
	return "books"
}
