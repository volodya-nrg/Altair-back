package response

import "altair/storage"

type PropWithKindName struct {
	*storage.Prop
	KindPropName string `json:"kindPropName" gorm:"column:kind_prop_name"`
}
