package storage

// Prop - структура таблицы свойств
type Prop struct {
	PropID         uint64 `json:"propID" gorm:"primaryKey;column:prop_id"`
	Title          string `json:"title" gorm:"column:title"`
	KindPropID     uint64 `json:"kindPropID" gorm:"column:kind_prop_id"`
	Name           string `json:"name" gorm:"column:name"`
	Suffix         string `json:"suffix" gorm:"column:suffix"`
	Comment        string `json:"comment" gorm:"column:comment"`
	PrivateComment string `json:"privateComment" gorm:"column:private_comment"`
}
