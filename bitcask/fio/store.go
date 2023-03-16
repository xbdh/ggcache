package fio

type Store interface {
	Read([]byte, uint64) (int, error)
	Append([]byte) (int, error)
	Sync() error
	Close() error
	Size() (uint64, error)
}

func NewStore(filename string) (Store, error) {
	return NewFile(filename)
}
