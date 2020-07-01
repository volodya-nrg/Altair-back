package request

// PostCat - структура запроса на добавление категории
type PostCat struct {
	Name                string               `form:"name" binding:"required"`
	ParentID            string               `form:"parentId"`
	Pos                 uint64               `form:"pos"`
	PriceAlias          string               `form:"priceAlias"`
	PriceSuffix         string               `form:"priceSuffix"`
	TitleHelp           string               `form:"titleHelp"`
	TitleComment        string               `form:"titleComment"`
	IsAutogenerateTitle bool                 `form:"isAutogenerateTitle"`
	PropsAssignedForCat []PropAssignedForCat `form:"propsAssignedForCat"`
}
