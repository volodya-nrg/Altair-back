package response

type CatTreeFull struct {
	*СatFull
	Childes []*CatTreeFull `json:"childes"`
}
