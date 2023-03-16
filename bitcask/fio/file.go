package fio

import (
	"os"
)

type File struct {
	fd *os.File
}

func (f *File) Read(b []byte, offset uint64) (int, error) {
	return f.fd.ReadAt(b, int64(offset))

}

func (f *File) Append(b []byte) (int, error) {
	return f.fd.Write(b)
}

func (f *File) Sync() error {
	return f.fd.Sync()
}

func (f *File) Close() error {
	return f.fd.Close()
}

func (f *File) Size() (uint64, error) {
	stat, err := f.fd.Stat()
	return uint64(stat.Size()), err
}

func NewFile(filename string) (*File, error) {

	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &File{fd: fd}, nil
}
