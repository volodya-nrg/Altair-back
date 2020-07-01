package storage

// Cat - структура таблицы категорий
type Cat struct {
	CatID               uint64 `json:"catId" gorm:"primary_key;column:cat_id"`
	Name                string `json:"name" gorm:"column:name"`
	Slug                string `json:"slug" gorm:"column:slug"`
	ParentID            uint64 `json:"parentId" gorm:"column:parent_id"`
	Pos                 uint64 `json:"pos" gorm:"column:pos"`
	IsDisabled          bool   `json:"isDisabled" gorm:"column:is_disabled"`
	PriceAlias          string `json:"priceAlias" gorm:"column:price_alias"`
	PriceSuffix         string `json:"priceSuffix" gorm:"column:price_suffix"`
	TitleHelp           string `json:"titleHelp" gorm:"column:title_help"`
	TitleComment        string `json:"titleComment" gorm:"column:title_comment"`
	IsAutogenerateTitle bool   `json:"isAutogenerateTitle" gorm:"column:is_autogenerate_title"`
}
