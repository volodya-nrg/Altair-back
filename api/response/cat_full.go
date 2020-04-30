package response

import "altair/storage"

type СatFull struct {
	*storage.Cat
	PropsFull []*PropFull `json:"props"`
}
