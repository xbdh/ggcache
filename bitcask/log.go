package bitcask

import (
	"fmt"
	"ggcache/bitcask/data"
	"ggcache/bitcask/index"
	"ggcache/bitcask/segment"

	//"gihub.com/xbdh/ggcache/bitcask/data"
	//"gihub.com/xbdh/bitcask/index"
	//"gihub.com/xbdh/bitcask/segment"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Log struct {
	mu            sync.RWMutex
	ops           Options
	activeSegment *segment.Segment
	oldSegments   map[uint32]*segment.Segment
	index         index.Index
	fids          []uint32
}

func NewLog(ops Options) (*Log, error) {
	if ops.MaxStoreSize == 0 {
		ops.MaxStoreSize = 1024
	}
	idx, err := index.NewIndex()
	if err != nil {
		return nil, err
	}
	l := &Log{
		ops:         ops,
		oldSegments: make(map[uint32]*segment.Segment),
		index:       idx,
	}
	err = l.initSegment()
	err = l.initIndex()
	return l, err

}
func (l *Log) initIndex() error {
	for i, fid := range l.fids {
		var sg *segment.Segment
		if fid == l.activeSegment.FileId {
			sg = l.activeSegment
		} else {
			sg = l.oldSegments[fid]
		}
		var offset uint64
		for { //要读到最后才知道后面的WriteOffset
			record, size, err := sg.ReadRecord(offset)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			pos := &data.RecordPos{
				Fid:    fid,
				Offset: offset,
			}
			err = l.index.Put(record.Key, pos)

			if err != nil {
				return err
			}
			offset += size
		}

		if i == len(l.fids)-1 {
			l.activeSegment.WriteOffset = offset
		}
	}
	return nil
}
func (l *Log) initSegment() error {
	//fmt.Println(l.ops.Dir)
	entries, err := os.ReadDir(l.ops.Dir)
	if err != nil {
		return err
	}
	var fileIds []uint32

	for _, entry := range entries {
		nameStr := strings.TrimSuffix(entry.Name(), path.Ext(entry.Name()))
		fileId, err := strconv.ParseUint(nameStr, 10, 32)

		if err != nil {
			continue
		}
		fileIds = append(fileIds, uint32(fileId))

	}
	l.fids = fileIds

	sort.Slice(fileIds, func(i, j int) bool {
		return fileIds[i] < fileIds[j]
	})
	for i := 0; i < len(fileIds); i++ {
		if err := l.newSegment(l.ops.Dir, fileIds[i]); err != nil {
			return err
		}
	}

	if l.activeSegment == nil {
		newSegment, err := segment.NewSegment(l.ops.Dir, 0)
		if err != nil {
			return err
		}
		l.activeSegment = newSegment
	}
	//fmt.Println("l.activeSegment:=", l.activeSegment)
	return nil
}
func (l *Log) Read(key []byte) ([]byte, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if len(key) == 0 {
		return nil, ErrKeyIsEmpty
	}
	pos, err := l.index.Get(key)
	if err != nil {
		return nil, err
	}
	if pos == nil {
		return nil, ErrKeyNotFound
	}
	var sg *segment.Segment

	if pos.Fid == l.activeSegment.FileId {
		sg = l.activeSegment
	} else {
		sg = l.oldSegments[pos.Fid]
	}
	if sg == nil {
		return nil, ErrDataFileNotFound
	}

	record, _, err := sg.ReadRecord(pos.Offset)
	if err != nil {
		return nil, err
	}

	return record.Value, nil
}

func (l *Log) Append(key []byte, value []byte) error {
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}
	data := data.Record{
		Key:   key,
		Value: value,
		//Type:  data.RecordTypeNormal,
	}
	logRecordPos, err := l.appendRecord(&data)
	fmt.Println("logRecordPos:=", logRecordPos)

	if err != nil {
		return err
	}
	err = l.index.Put(key, logRecordPos)
	if err != nil {
		return ErrIndexUpdateFailed
	}

	return nil

}

func (l *Log) appendRecord(record *data.Record) (*data.RecordPos, error) {

	l.mu.Lock()
	defer l.mu.Unlock()

	offset, err := l.activeSegment.Append(record)
	if err != nil {
		return nil, err
	}
	recordPos := &data.RecordPos{
		Fid:    l.activeSegment.FileId,
		Offset: offset,
	}
	if l.activeSegment.IsMaxed(l.ops.MaxStoreSize) {
		l.newSegment(l.ops.Dir, l.activeSegment.FileId+1)
	}

	return recordPos, nil
}

func (l *Log) newSegment(dir string, fileId uint32) error {
	s, err := segment.NewSegment(dir, fileId)
	if err != nil {
		return err
	}
	if l.activeSegment != nil {
		l.oldSegments[l.activeSegment.FileId] = l.activeSegment
	}

	l.activeSegment = s
	return nil
}
