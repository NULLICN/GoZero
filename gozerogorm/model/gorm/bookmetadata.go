package gorm

type Bookmetadata struct {
	BooksBookmetadataId int    `json:"books_bookmetadata_id"`
	Extrainfo       string `json:"extrainfo"`
	//Book      []Book `json:"book"`
}

func (Bookmetadata) TableName() string {
	return "bookmetadata"
}
