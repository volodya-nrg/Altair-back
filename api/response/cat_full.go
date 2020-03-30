package response

import "altair/storage"

type СatFull struct {
	*storage.Cat
	PropertiesFull []*PropertyFull `json:"properties"`
}
