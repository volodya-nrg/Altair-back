package storage

type CatProp struct {
	CatId         uint64 `json:"catId" gorm:"column:cat_id"`
	PropId        uint64 `json:"propId" gorm:"column:prop_id"`
	Pos           uint64 `json:"pos" gorm:"column:pos"`
	IsRequire     bool   `json:"isRequire" gorm:"column:is_require"`
	IsCanAsFilter bool   `json:"isCanAsFilter" gorm:"column:is_can_as_filter"`
	Comment       string `json:"comment" gorm:"column:comment"`
}
