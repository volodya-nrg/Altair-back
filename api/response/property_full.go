package response

import "altair/storage"

type PropertyFull struct {
	*storage.Property
	KindPropertyName string                   `json:"kindPropertyName" gorm:"column:kind_property_name"`
	PropertyPos      uint64                   `json:"propertyPos,omitempty" gorm:"column:property_pos"`
	IsRequire        bool                     `json:"propertyIsRequire,omitempty" gorm:"column:property_is_require"`
	Values           []*storage.ValueProperty `json:"values"`
}
