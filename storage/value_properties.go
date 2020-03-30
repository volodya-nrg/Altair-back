package storage

type ValueProperty struct {
	ValueId    uint64 `json:"valueId" gorm:"primary_key;column:value_id"`
	Title      string `json:"title" gorm:"column:title"`
	Pos        uint64 `json:"pos" gorm:"column:pos"`
	PropertyId uint64 `json:"propertyId" gorm:"column:property_id"`
}
