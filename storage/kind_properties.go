package storage

type KindProperty struct {
	KindPropertyId uint64 `json:"kindPropertyId" gorm:"primary_key;column:kind_property_id"`
	Name           string `json:"name" gorm:"column:name"`
}
