package index

import "ggcache/bitcask/data"

type Index interface {
	Put(key []byte, pos *data.RecordPos) error
	Get(key []byte) (*data.RecordPos, error)
	Delete(key []byte) error
}

func NewIndex() (Index, error) {
	return NewBTree(), nil
}
