package response

import "altair/storage"

type CatTree struct {
	*storage.Cat
	Childes []*CatTree `json:"childes"`
}
