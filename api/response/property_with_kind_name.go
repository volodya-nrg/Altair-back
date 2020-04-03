package response

import "altair/storage"

type PropertyWithKindName struct {
	*storage.Property
	KindPropertyName string `json:"kindPropertyName" gorm:"column:kind_property_name"`
}
