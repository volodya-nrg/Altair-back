package response

import "altair/storage"

// AdFull - структура ответа, объявление (полное)
type AdFull struct {
	*storage.Ad
	Images     []*storage.Image `json:"images"`
	DetailsExt []*AdDetailExt   `json:"detailsExt"`
}
