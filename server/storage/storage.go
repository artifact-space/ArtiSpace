package storage

import "github.com/artifact-space/ArtiSpace/config"

type BinaryStorage interface {
	Init() error

	GetFile(filePath string) ([]byte, error)

	PutFile(filePath string, data []byte) error

	ListFiles(filePath string) ([]string, error)

	RenameFile(oldPath, newPath string) error

	PutFileChunk(filePath string, chunk []byte, offset int64) error

	DeleteFile(filePath string) error

	Size(filePath string) (int64, error)
}

var Storage BinaryStorage

func initLFS(config *config.StorageConfig) {
	Storage = NewLFS(config.Properties)
}

func InitStorage(config *config.StorageConfig) error {
	//for now, We'll initalize LFS, later we have to initalize based on the configurations
	initLFS(config)
	return nil
}
