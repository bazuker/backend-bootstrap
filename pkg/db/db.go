package db

type Adapter interface {
	// CreateUser creates a new user.
	CreateUser(user *User) error
	// UpdateUserPhotoURL updates user's photo URL.
	UpdateUserPhotoURL(userID, photoURL string) error
	// GetUserByID finds a user by user ID.
	GetUserByID(ID string) (User, error)
	// GetUserByEmail finds a user by user email.
	GetUserByEmail(email string) (User, error)
}
