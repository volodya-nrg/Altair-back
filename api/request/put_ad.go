package request

type PutAd struct {
	AdId            uint64   `form:"adId" binding:"required"`
	Title           string   `form:"title" binding:"required"`
	CatId           uint64   `form:"catId" binding:"required"`
	UserId          uint64   `form:"userId"`
	Text            string   `form:"text"`
	Price           uint64   `form:"price"`
	IsDisabled      bool     `form:"isDisabled"`
	FilesAlreadyHas []string `form:"filesAlreadyHas[]"`
}
