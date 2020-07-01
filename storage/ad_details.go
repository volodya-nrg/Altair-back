package storage

// AdDetail - структура таблицы деталей объявлений
type AdDetail struct {
	AdID   uint64 `json:"adId" gorm:"column:ad_id"`
	PropID uint64 `json:"propId" gorm:"column:prop_id"`
	Value  string `json:"value" gorm:"column:value"`
}
