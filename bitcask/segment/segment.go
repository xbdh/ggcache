package segment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"ggcache/bitcask/data"
	"ggcache/bitcask/fio"
)

var (
	keyLenSize = 8
	valLenSize = 8
)

type Segment struct {
	FileId      uint32
	WriteOffset uint64 //写数据的启示位置,activeSegment
	fio         fio.Store
}

func NewSegment(dir string, fileId uint32) (*Segment, error) {
	filename := fmt.Sprintf("%s/%d.store", dir, fileId)
	fio, err := fio.NewStore(filename)
	if err != nil {
		return nil, err
	}
	return &Segment{
		FileId: fileId,
		fio:    fio,
	}, nil
}
func (s *Segment) Sync() error {
	s.fio.Sync()
	return nil
}
func (s *Segment) Append(record *data.Record) (uint64, error) {

	//str := fmt.Sprintf("%s%s", record.Key, record.Value)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint64(len(record.Key)))
	binary.Write(buf, binary.BigEndian, uint64(len(record.Value)))
	_, err := buf.Write(record.Key)

	_, err = buf.Write(record.Value)
	if err != nil {
		return 0, err
	}

	recordoffset := s.WriteOffset
	n, err := s.fio.Append(buf.Bytes())
	s.WriteOffset += uint64(n)

	return recordoffset, err
}
func (s *Segment) Read(d []byte, offset uint64) (uint64, error) {
	return 0, nil
}

func (s *Segment) ReadRecord(offset uint64) (*data.Record, uint64, error) {
	keyLenBuf := make([]byte, keyLenSize)
	ValueLenBuf := make([]byte, valLenSize)
	_, err := s.fio.Read(keyLenBuf, offset)
	if err != nil {
		return nil, 0, err
	}
	_, err = s.fio.Read(ValueLenBuf, offset+uint64(keyLenSize))
	if err != nil {
		return nil, 0, err
	}
	keyLen := binary.BigEndian.Uint64(keyLenBuf)
	valLen := binary.BigEndian.Uint64(ValueLenBuf)
	//fmt.Println("keyLen:=", keyLen, "valLen:=", valLen)
	key := make([]byte, keyLen)
	val := make([]byte, valLen)
	_, err = s.fio.Read(key, offset+uint64(keyLenSize)+uint64(valLenSize))
	if err != nil {
		return nil, 0, err
	}
	_, err = s.fio.Read(val, offset+uint64(keyLenSize)+uint64(valLenSize)+keyLen)
	if err != nil {
		return nil, 0, err
	}
	return &data.Record{
		Key:   key,
		Value: val,
	}, keyLen + valLen + uint64(keyLenSize) + uint64(valLenSize), nil

}

func (s *Segment) IsMaxed(maxSize uint64) bool {
	size, _ := s.fio.Size()

	return size >= maxSize
}
