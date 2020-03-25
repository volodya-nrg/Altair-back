package response

import "altair/storage"

type AdFull struct {
	*storage.Ad
	Images []*storage.Image `json:"images"`
}
