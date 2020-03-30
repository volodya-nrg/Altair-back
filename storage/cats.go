package storage

type Cat struct {
	CatId      uint64 `json:"catId" gorm:"primary_key;column:cat_id"`
	Name       string `json:"name" gorm:"column:name"`
	Slug       string `json:"slug" gorm:"column:slug"`
	ParentId   uint64 `json:"parentId" gorm:"column:parent_id"`
	Pos        uint64 `json:"pos" gorm:"column:pos"`
	IsDisabled bool   `json:"isDisabled" gorm:"column:is_disabled"`
}
