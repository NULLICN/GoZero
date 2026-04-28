package gorm

type Book struct {
	Id             int          `json:"id"`
	Bookname       string       `json:"bookname"`
	Price          float64      `json:"price"`
	BookmetadataId int          `json:"bookmetadata_id"`
	// 此表目前是主表，外键id在Bookmetadata表中，gorm无法自动识别，所以需要手动指定外键和关联字段。解释：指定了Bookmetadata表的id为外键，且参照book表的BookmetadataId字段。
	Bookmetadata   Bookmetadata `gorm:"foreignKey:Id;references:BookmetadataId" json:"bookmetadata"`
}

func (Book) TableName() string {
	return "books"
}
