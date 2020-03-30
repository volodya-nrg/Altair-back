package storage

type AdDetail struct {
	AdId       uint64 `json:"adId" gorm:"column:ad_id"`
	PropertyId uint64 `json:"propertyId" gorm:"column:property_id"`
	Value      string `json:"value" gorm:"column:value"`
}
