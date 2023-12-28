package local

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/bazuker/backend-bootstrap/pkg/db"
)

// DB represents a filesystem database suitable for local tests.
type DB struct {
	storage *localStorage
	mx      sync.Mutex
	cfg     Config
}

type localStorage struct {
	Users []db.User
}

type Config struct {
	Filename string
}

func New(cfg Config) (*DB, error) {
	database := &DB{
		storage: &localStorage{
			Users: []db.User{},
		},
		cfg: cfg,
		mx:  sync.Mutex{},
	}

	if _, err := os.Stat(cfg.Filename); err == nil {
		if err := database.loadStorage(); err != nil {
			return nil, err
		}
	} else {
		if err := database.saveStorage(); err != nil {
			return nil, err
		}
	}

	return database, nil
}

func (d *DB) CreateUser(user *db.User) error {
	d.mx.Lock()
	defer d.mx.Unlock()

	d.storage.Users = append(d.storage.Users, *user)

	return d.saveStorage()
}

func (d *DB) UpdateUserPhotoURL(userID, photoURL string) error {
	d.mx.Lock()
	defer d.mx.Unlock()

	userIndex := -1
	for i := range d.storage.Users {
		if d.storage.Users[i].ID == userID {
			userIndex = i
			break
		}
	}
	if userIndex < 0 {
		return db.ErrNotFound
	}
	d.storage.Users[userIndex].PhotoURL = photoURL

	return d.saveStorage()
}

func (d *DB) GetUserByID(id string) (db.User, error) {
	if id == "" {
		return db.User{}, errors.New("missing id")
	}

	d.mx.Lock()
	defer d.mx.Unlock()

	for i := range d.storage.Users {
		if d.storage.Users[i].ID == id {
			return d.storage.Users[i], nil
		}
	}

	return db.User{}, db.ErrNotFound
}

func (d *DB) GetUserByEmail(email string) (db.User, error) {
	if email == "" {
		return db.User{}, errors.New("missing email")
	}

	d.mx.Lock()
	defer d.mx.Unlock()

	for i := range d.storage.Users {
		if d.storage.Users[i].Email == email {
			return d.storage.Users[i], nil
		}
	}

	return db.User{}, db.ErrNotFound
}

func (d *DB) saveStorage() error {
	data, _ := json.Marshal(d.storage)
	return os.WriteFile(d.cfg.Filename, data, os.ModePerm)
}

func (d *DB) loadStorage() error {
	data, err := os.ReadFile(d.cfg.Filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, d.storage)
}
