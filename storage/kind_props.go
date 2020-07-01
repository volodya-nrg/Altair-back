package storage

// KindProp - структура таблицы видов свойств
type KindProp struct {
	KindPropID uint64 `json:"kindPropId" gorm:"primary_key;column:kind_prop_id"`
	Name       string `json:"name" gorm:"column:name"`
}
