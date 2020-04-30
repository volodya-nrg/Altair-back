package storage

type AdDetail struct {
	AdId   uint64 `json:"adId" gorm:"column:ad_id"`
	PropId uint64 `json:"propId" gorm:"column:prop_id"`
	Value  string `json:"value" gorm:"column:value"`
}
