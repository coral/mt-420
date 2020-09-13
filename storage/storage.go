package storage

import "os"

type Events struct {
}

type Storage interface {
	Init()
	ListFiles() []os.FileInfo
	LoadFile(f os.FileInfo) ([]byte, error)
}
