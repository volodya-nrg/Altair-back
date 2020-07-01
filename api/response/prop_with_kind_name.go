package response

import "altair/storage"

// PropWithKindName - структура ответа, свойство с именем (разновидность)
type PropWithKindName struct {
	*storage.Prop
	KindPropName string `json:"kindPropName" gorm:"column:kind_prop_name"`
}
