package bitcask

import "errors"

var (
	ErrKeyNotFound       = errors.New("key not found")
	ErrKeyIsEmpty        = errors.New("key is empty")
	ErrIndexUpdateFailed = errors.New("index update failed")
	ErrDataFileNotFound  = errors.New("data file not found")
)
