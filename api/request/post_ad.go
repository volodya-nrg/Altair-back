package request

type PostAd struct {
	Title       string `form:"title"` // может быть и не обязательным. Нужно помнить про Slug
	CatId       uint64 `form:"catId" binding:"required"`
	Description string `form:"description" binding:"required"`
	Price       uint64 `form:"price"` // может быть ноль
	UserId      uint64 `form:"userId"`
	Youtube     string `form:"youtube"`
}
