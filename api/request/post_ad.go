package request

type PostAd struct {
	Title  string `form:"title" binding:"required"`
	CatId  uint64 `form:"catId" binding:"required"`
	UserId uint64 `form:"userId"`
	Text   string `form:"text"`
	Price  uint64 `form:"price"`
}
