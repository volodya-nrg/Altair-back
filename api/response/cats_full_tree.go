package response

type CatFullTree struct {
	*СatFull
	Childes []*CatFullTree `json:"childes"`
}
