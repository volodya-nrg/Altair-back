package request

type PutAd struct {
	PostAd
	AdId            uint64   `form:"adId" binding:"required"`
	IsDisabled      bool     `form:"isDisabled"`
	FilesAlreadyHas []string `form:"filesAlreadyHas[]"`
}
