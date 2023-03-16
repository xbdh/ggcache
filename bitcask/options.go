package bitcask

type Options struct {
	Dir          string
	MaxStoreSize uint64
	SyncWrite    bool
}
