package data

type LogRecordType = byte

const (
	RecordTypeNormal LogRecordType = iota
	RecordTypeDelete
)

type Record struct {
	Key   []byte
	Value []byte
	//Type  LogRecordType
}

func (r Record) Size() uint64 {
	return 0
}

type RecordPos struct {
	Fid    uint32
	Offset uint64
}

func EncodeRecord(data *Record) ([]byte, uint64) {
	return nil, 0
}
