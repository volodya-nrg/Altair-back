package storage

// ValueProp - структура таблицы значений для свойств
type ValueProp struct {
	ValueID uint64 `json:"valueID" gorm:"primaryKey;column:value_id"`
	Title   string `json:"title" gorm:"column:title"`
	Pos     uint64 `json:"pos" gorm:"column:pos"`
	PropID  uint64 `json:"propID" gorm:"column:prop_id"`
}
