package gorm

type Book struct {
	Id             int            `json:"id"`
	Bookname       string         `json:"bookname"`
	Price          float64        `json:"price"`
	BookmetadataId int            `json:"bookmetadata_id"`
	Bookmetadata   []Bookmetadata `gorm:"foreignKey:BooksBookmetadataId;references:BookmetadataId" json:"bookmetadata"`
}

func (Book) TableName() string {
	return "books"
}
