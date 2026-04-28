package gorm

type Bookmetadata struct {
	Id int    `json:"id"`
	Extrainfo       string `json:"extrainfo"`
	//Book      []Book `json:"book" gorm:"foreignKey:BookmetadataId;references:Id"`
}

func (Bookmetadata) TableName() string {
	return "bookmetadata"
}
