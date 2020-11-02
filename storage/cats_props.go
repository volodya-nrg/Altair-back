package storage

// CatProp - структура таблицы свойств категорий
type CatProp struct {
	CatID         uint64 `json:"catID" gorm:"column:cat_id"`
	PropID        uint64 `json:"propID" gorm:"column:prop_id"`
	Pos           uint64 `json:"pos" gorm:"column:pos"`
	IsRequire     bool   `json:"isRequire" gorm:"column:is_require"`
	IsCanAsFilter bool   `json:"isCanAsFilter" gorm:"column:is_can_as_filter"`
	Comment       string `json:"comment" gorm:"column:comment"`
}
