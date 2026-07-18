package data

import "errors"

var ErrRecordNotFound = errors.New("record not found")
var ErrUpdateConflict = errors.New("update conflict")