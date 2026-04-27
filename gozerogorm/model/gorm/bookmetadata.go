package gorm

type Bookmetadata struct {
	Id        int    `gorm:"primaryKey" json:"id"`
	Extrainfo string `json:"extrainfo"`
	Book      []Book `json:"book"`
}

func (Bookmetadata) TableName() string {
	return "bookmetadata"
}
