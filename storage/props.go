package storage

type Prop struct {
	PropId         uint64 `json:"propId" gorm:"primary_key;column:prop_id"`
	Title          string `json:"title" gorm:"column:title"`
	KindPropId     uint64 `json:"kindPropId" gorm:"column:kind_prop_id"`
	Name           string `json:"name" gorm:"column:name"`
	Suffix         string `json:"suffix" gorm:"column:suffix"`
	Comment        string `json:"comment" gorm:"column:comment"`
	PrivateComment string `json:"privateComment" gorm:"column:private_comment"`
}
