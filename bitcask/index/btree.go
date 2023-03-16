package index

import (
	"errors"
	"ggcache/bitcask/data"
	"github.com/google/btree"
	"sync"
)

type BTree struct {
	tree *btree.BTree
	mu   sync.RWMutex
}

func NewBTree() *BTree {
	return &BTree{
		tree: btree.New(32),
	}
}
func (b *BTree) Put(key []byte, pos *data.RecordPos) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tree.ReplaceOrInsert(&item{key, pos})
	return nil
}
func (b *BTree) Get(key []byte) (*data.RecordPos, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	it := b.tree.Get(&item{key, nil})
	if it == nil {
		return nil, errors.New("not found")
	}
	return it.(*item).pos, nil

}
func (b *BTree) Delete(key []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tree.Delete(&item{key, nil})
	return nil
}

type item struct {
	key []byte
	pos *data.RecordPos
}

func (i *item) Less(than btree.Item) bool {
	return string(i.key) < string(than.(*item).key)
}
