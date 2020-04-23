package request

type PostCat struct {
	Name                string `form:"name" binding:"required"`
	ParentId            uint64 `form:"parentId"`
	Pos                 uint64 `form:"pos"`
	PriceAlias          string `form:"priceAlias"`
	PriceSuffix         string `form:"priceSuffix"`
	TitleHelp           string `form:"titleHelp"`
	TitleComment        string `form:"titleComment"`
	IsAutogenerateTitle bool   `form:"isAutogenerateTitle"`
}
