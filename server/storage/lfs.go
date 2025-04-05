package storage

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/artifact-space/ArtiSpace/lib"
	"github.com/artifact-space/ArtiSpace/log"
)

type localFileStorage struct {
	storageDir string
	keyLock    *lib.KeyLock
}

func NewLFS(props map[string]string) *localFileStorage {
	storagePath, _ := props["fs.storage.path"]
	return &localFileStorage{
		storageDir: storagePath,
		keyLock:    lib.NewKeyLock(),
	}
}

func (lfs *localFileStorage) Init() error {
	fileInfo, err := os.Stat(lfs.storageDir)
	if err == nil && !fileInfo.IsDir() {
		log.Logger().Error().Msgf("file: %s is not a directory. please remove the file and start the server", lfs.storageDir)
		return fmt.Errorf("file: %s is not a directory", lfs.storageDir)
	}

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Logger().Error().Err(err).Msg("unable to initialize local file system storage")
		return err
	}

	if err != nil && errors.Is(err, fs.ErrNotExist) {
		log.Logger().Debug().Msgf("storage directory: %s already exists", lfs.storageDir)
		return nil
	}

	err = os.MkdirAll(lfs.storageDir, os.FileMode(0666))
	if err != nil {
		log.Logger().Error().Err(err).Msgf("unable to intialize local file system storage : %s", lfs.storageDir)
		return err
	}

	return nil
}

func (lfs *localFileStorage) GetFile(filePath string) ([]byte, error) {
	targetFilePath := filepath.Join(lfs.storageDir, filePath)

	file, err := os.Open(targetFilePath)

	if err != nil && os.IsNotExist(err) {
		log.Logger().Error().Msgf("file: %s does not exist.", targetFilePath)
		return nil, err
	} else if err != nil {
		log.Logger().Error().Err(err).Msgf("unexpected error occured while opening file: %s", targetFilePath)
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Logger().Error().Err(err).Msgf("File: %s exists but cannot be read", targetFilePath)
		return nil, err
	}

	return data, nil
}

func (lfs *localFileStorage) PutFile(filePath string, data []byte) error {
	targetFilePath := filepath.Join(lfs.storageDir, filePath)

	// create necessary directories
	err := os.MkdirAll(filepath.Dir(targetFilePath), 0666)

	if err != nil {
		log.Logger().Error().Err(err).Msgf("unable to create directories: %s", filepath.Dir(targetFilePath))
		return err
	}

	file, err := os.OpenFile(targetFilePath, os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Logger().Error().Err(err).Msgf("unable to create file: %s", targetFilePath)
		return err
	}
	defer file.Close()

	n, err := file.Write(data)
	if err != nil || n != len(data) {
		log.Logger().Error().Err(err).Msg("file will be removed as the file write is not complete")

		err = os.Remove(targetFilePath)
		log.Logger().Error().Err(err).Msgf("removing file %s failed", targetFilePath)
		return err
	}

	return nil
}

func (lfs *localFileStorage) ListFiles(filePath string) ([]string, error) {
	var files []string

	targetFilepath := filepath.Join(lfs.storageDir, filePath)

	f, err := os.Open(targetFilepath)
	if err != nil {
		log.Logger().Error().Err(err).Msgf("unable to open file: %s", targetFilepath)
		return files, err
	}
	defer f.Close()

	fileInfos, err := f.ReadDir(-1)
	if err != nil {
		log.Logger().Err(err).Msgf("unable to read the directory: %s", targetFilepath)
		return files, err
	}

	for _, fi := range fileInfos {
		if !fi.IsDir() {
			files = append(files, fi.Name())
		}
	}

	return files, nil
}

func (lfs *localFileStorage) RenameFile(oldPath, newPath string) error {
	oldFullPath := filepath.Join(lfs.storageDir, oldPath)
	newFullPath := filepath.Join(lfs.storageDir, newPath)

	if err := os.MkdirAll(filepath.Dir(newFullPath), 0666); err != nil {
		log.Logger().Error().Err(err).Msgf("unable to create directories: %s", filepath.Dir(newFullPath))
		return err
	}

	if err := os.Rename(oldFullPath, newFullPath); err != nil {
		log.Logger().Error().Err(err).Msgf("unable to rename file: %s to %s", oldFullPath, newFullPath)
		return err
	}

	return nil
}

func (lfs *localFileStorage) DeleteFile(filePath string) error {
	targetFilePath := filepath.Join(lfs.storageDir, filePath)

	err := os.Remove(targetFilePath)
	if err != nil {
		log.Logger().Error().Err(err).Msgf("error occured when removing file: %s", targetFilePath)
	}
	return err
}

func (lfs *localFileStorage) Size(filePath string) (int64, error) {
	targetFilePath := filepath.Join(lfs.storageDir, filePath)

	fileInfo, err := os.Stat(targetFilePath)
	if err != nil {
		log.Logger().Error().Err(err).Msgf("unable to retrive file info: %s", targetFilePath)
		return -1, err
	}

	return fileInfo.Size(), nil
}

func (lfs *localFileStorage) PutFileChunk(filePath string, chunk []byte, offset int64) error {
	targetFilePath := filepath.Join(lfs.storageDir, filePath)

	lfs.keyLock.Lock(targetFilePath)
	defer lfs.keyLock.Unlock(targetFilePath)

	var file *os.File

	_, err := os.Stat(targetFilePath)
	if err == nil {
		file, err = os.OpenFile(targetFilePath, os.O_APPEND, 0666)

		if err != nil {
			log.Logger().Error().Err(err).Msgf("unable to open file to append chunk: %s", targetFilePath)
			return err
		}

	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(targetFilePath), 0666); err != nil {
			log.Logger().Error().Err(err).Msgf("unable to create directories for path: %s", filepath.Dir(targetFilePath))
			return err
		}

		file, err = os.OpenFile(targetFilePath, os.O_CREATE|os.O_WRONLY, 0666)

		if err != nil {
			log.Logger().Error().Err(err).Msgf("unable to open file to append chunk: %s", targetFilePath)
			return err
		}

	} else {
		log.Logger().Error().Err(err).Msgf("unexpected error occured when checking file existence: %s", targetFilePath)
		return err
	}

	defer file.Close()

	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		log.Logger().Error().Err(err).Msgf("error occured when seeking file cursor, path: %s, offset: %d", targetFilePath, offset)
		return err
	}

	n, err := file.Write(chunk)
	if err != nil {
		log.Logger().Error().Err(err).Msgf("chunk write failed")
		return err
	}
	if len(chunk) != n {
		err = fmt.Errorf("only %d bytes were written out of %d", n, len(chunk))
		log.Logger().Error().Err(err)
		return err
	}
	return nil
}
