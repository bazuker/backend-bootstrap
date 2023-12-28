package filestore

type FileStore interface {
	// GetObject retrieves the object from the file store.
	GetObject(key string) ([]byte, error)
	// PutObject creates or overwrites the object in the filestore.
	PutObject(object []byte, key string) error
	// DeleteObject deletes the object from the filestore.
	DeleteObject(key string) error
}
