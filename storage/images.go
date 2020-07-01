package storage

import "time"

// Image - структура таблицы изображений
type Image struct {
	ImgID      uint64     `json:"imgId" gorm:"primary_key;column:img_id"`
	Filepath   string     `json:"filepath" gorm:"column:filepath"`
	ElID       uint64     `json:"elId" gorm:"column:el_id"`
	IsDisabled bool       `json:"isDisabled" gorm:"column:is_disabled"`
	Opt        string     `json:"opt" gorm:"column:opt"`
	CreatedAt  *time.Time `json:"createdAt" gorm:"column:created_at"`
}
