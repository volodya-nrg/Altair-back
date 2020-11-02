package request

// PutCat - структура запроса изменение категории
type PutCat struct {
	CatID               uint64               `form:"catID" binding:"required"`
	Name                string               `form:"name" binding:"required"`
	ParentID            uint64               `form:"parentID"`
	Pos                 uint64               `form:"pos"`
	IsDisabled          bool                 `form:"isDisabled"`
	PriceAlias          string               `form:"priceAlias"`
	PriceSuffix         string               `form:"priceSuffix"`
	TitleHelp           string               `form:"titleHelp"`
	TitleComment        string               `form:"titleComment"`
	IsAutogenerateTitle bool                 `form:"isAutogenerateTitle"`
	PropsAssignedForCat []PropAssignedForCat `form:"propsAssignedForCat"`
}
