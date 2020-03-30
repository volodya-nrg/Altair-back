package response

import "altair/storage"

type AdDetailExt struct {
	storage.AdDetail
	PropertyName     string `json:"propertyName" gorm:"column:property_name"`
	KindPropertyName string `json:"kindPropertyName" gorm:"column:kind_property_name"`
}
