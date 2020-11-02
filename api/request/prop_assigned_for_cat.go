package request

// PropAssignedForCat - структура запроса прикрепленных к категории
type PropAssignedForCat struct {
	PropID        uint64 `form:"propID" binding:"required"`
	Comment       string `form:"comment"`
	Pos           uint64 `form:"pos"`
	IsRequire     bool   `json:"isRequire"`
	IsCanAsFilter bool   `json:"isCanAsFilter"`
}
