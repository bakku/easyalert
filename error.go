package easyalert

import "errors"

// ErrRecordDoesNotExist is a generic error in case a record which is searched does not exist
var ErrRecordDoesNotExist = errors.New("record does not exist")
