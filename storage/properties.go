package storage

type Property struct {
	PropertyId     uint64 `json:"propertyId" gorm:"primary_key;column:property_id"`
	Title          string `json:"title" gorm:"column:title"`
	KindPropertyId uint64 `json:"kindPropertyId" gorm:"column:kind_property_id"`
	Name           string `json:"name" gorm:"column:name"`
	Suffix         string `json:"suffix" gorm:"column:suffix"`
	Comment        string `json:"comment" gorm:"column:comment"`
	PrivateComment string `json:"privateComment" gorm:"column:private_comment"`
}
