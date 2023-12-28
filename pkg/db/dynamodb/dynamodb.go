package dynamodb

import "github.com/bazuker/backend-bootstrap/pkg/db"

type DB struct {
	cfg Config
}

type Config struct {
	UsersTableName string
}

func New(cfg Config) *DB {
	return &DB{cfg: cfg}
}

func (D DB) CreateUser(user *db.User) error {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (D DB) UpdateUser(user *db.User) error {
	//TODO implement me
	panic("implement me")
}

func (D DB) GetUserByID(ID string) (*db.User, error) {
	//TODO implement me
	return &db.User{
		ID:        ID,
		FirstName: "Bob",
		LastName:  "Builder",
	}, nil
}

func (D DB) GetUserByEmail(email string) (*db.User, error) {
	//TODO implement me
	//panic("implement me")
	return nil, db.ErrNotFound
}
