package storage

// AdDetail - структура таблицы деталей объявлений
type AdDetail struct {
	AdID   uint64 `json:"adID" gorm:"column:ad_id"`
	PropID uint64 `json:"propID" gorm:"column:prop_id"`
	Value  string `json:"value" gorm:"column:value"`
}
