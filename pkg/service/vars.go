package service

import "errors"

var (
	errNotCreateNewAd         = errors.New("not create new ad")
	errOnNewRecordNewAd       = errors.New("err on NewRecord new ad")
	errEmptyListCatIds        = errors.New("err empty list on cat_ids")
	errNotCreateNewCat        = errors.New("not create new cat")
	errNotCreateNewCatProp    = errors.New("not create new cat prop")
	errNotCreateNewImage      = errors.New("not create new image")
	errOnNewRecordNewKindProp = errors.New("err on NewRecord new kindProp")
	errOnNewRecordNewProp     = errors.New("err on NewRecord new prop")
	errNotCorrectEmail        = errors.New("not correct email")
	errPasswordIsShort        = errors.New("password is short")
	errNotCreateNewUser       = errors.New("not create new user")
	errOnNewRecordNewAdDetail = errors.New("err on NewRecord new adDetail")
	errNotCreateNewAdDetail   = errors.New("not create new adDetail")
	minLenPassword            = 6
)
