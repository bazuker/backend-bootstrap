package local

import (
	"os"
	"path/filepath"
)

type FileStore struct {
	cfg Config
}

type Config struct {
	Directory string
}

func New(cfg Config) *FileStore {
	return &FileStore{
		cfg: cfg,
	}
}

func (f *FileStore) PutObject(object []byte, key string) error {
	path := filepath.Join(f.cfg.Directory, key)
	return os.WriteFile(path, object, os.ModePerm)
}

func (f *FileStore) GetObject(key string) ([]byte, error) {
	path := filepath.Join(f.cfg.Directory, key)
	return os.ReadFile(path)
}

func (f *FileStore) DeleteObject(key string) error {
	path := filepath.Join(f.cfg.Directory, key)
	return os.Remove(path)
}
