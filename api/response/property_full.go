package response

import "altair/storage"

type PropertyFull struct {
	*storage.Property
	KindPropertyName string                   `json:"kindPropertyName" gorm:"column:kind_property_name"`
	PropertyPos      uint64                   `json:"propertyPos" gorm:"column:property_pos"`
	IsRequire        bool                     `json:"propertyIsRequire" gorm:"column:property_is_require"`
	IsCanAsFilter    bool                     `json:"propertyIsCanAsFilter" gorm:"column:property_is_can_as_filter"`
	Comment          string                   `json:"propertyComment" gorm:"column:property_comment"`
	Values           []*storage.ValueProperty `json:"values"`
}
