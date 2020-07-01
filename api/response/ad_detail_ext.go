package response

import "altair/storage"

// AdDetailExt - структура ответа, детали объявления (расширенные)
type AdDetailExt struct {
	*storage.AdDetail
	PropName     string `json:"propName" gorm:"column:prop_name"`
	KindPropName string `json:"kindPropName" gorm:"column:kind_prop_name"`
	ValueName    string `json:"valueName" gorm:"column:value_name"`
}
