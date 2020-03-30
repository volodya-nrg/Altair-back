package service

import "errors"

var (
	errNotCreateNewAd             = errors.New("not create new ad")
	errOnNewRecordNewAd           = errors.New("err on NewRecord new ad")
	errEmptyListCatIds            = errors.New("err empty list on cat_ids")
	errNotCreateNewCat            = errors.New("not create new cat")
	errNotCreateNewCatProperty    = errors.New("not create new cat property")
	errNotCreateNewImage          = errors.New("not create new image")
	errOnNewRecordNewKindProperty = errors.New("err on NewRecord new kindProperty")
	errOnNewRecordNewProperty     = errors.New("err on NewRecord new property")
	errNotCorrectEmail            = errors.New("not correct email")
	errPasswordIsShort            = errors.New("password is short")
	errNotCreateNewUser           = errors.New("not create new user")
	errOnNewRecordNewAdDetail     = errors.New("err on NewRecord new adDetail")
	errNotCreateNewAdDetail       = errors.New("not create new adDetail")
	minLenPassword                = 6
)
