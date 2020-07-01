package response

import "altair/storage"

// PropFull - структура ответа, свойство (полное)
type PropFull struct {
	*storage.Prop
	KindPropName  string               `json:"kindPropName" gorm:"column:kind_prop_name"`
	PropPos       uint64               `json:"propPos" gorm:"column:prop_pos"`
	IsRequire     bool                 `json:"propIsRequire" gorm:"column:prop_is_require"`
	IsCanAsFilter bool                 `json:"propIsCanAsFilter" gorm:"column:prop_is_can_as_filter"`
	Comment       string               `json:"propComment" gorm:"column:prop_comment"`
	Values        []*storage.ValueProp `json:"values"`
}
