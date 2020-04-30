package storage

type KindProp struct {
	KindPropId uint64 `json:"kindPropId" gorm:"primary_key;column:kind_prop_id"`
	Name       string `json:"name" gorm:"column:name"`
}
