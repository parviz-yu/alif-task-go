package store

import "errors"

var (
	ErrRecordNotFound            = errors.New("record not found")
	ErrLimitExceededIdentified   = errors.New("limit exceeded for identified wallet — 100.000 TJS")
	ErrLimitExceededUnidentified = errors.New("limit exceeded for unidentified wallet — 10.000 TJS")
)
