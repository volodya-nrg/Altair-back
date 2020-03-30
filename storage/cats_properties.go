package storage

type CatProperty struct {
	CatId      uint64 `json:"catId" gorm:"column:cat_id"`
	PropertyId uint64 `json:"propertyId" gorm:"column:property_id"`
	Pos        uint64 `json:"pos" gorm:"column:pos"`
	IsRequire  bool   `json:"isRequire" gorm:"column:is_require"`
}
