package storage

type Property struct {
	PropertyId     uint64 `json:"propertyId" gorm:"primary_key;column:property_id"`
	Title          string `json:"title" gorm:"column:title"`
	KindPropertyId uint64 `json:"kindPropertyId" gorm:"column:kind_property_id"`
	Name           string `json:"name" gorm:"column:name"`
	IsCanAsFilter  bool   `json:"isCanAsFilter" gorm:"column:is_can_as_filter"`
	MaxInt         uint64 `json:"maxInt" gorm:"column:max_int"`
}
