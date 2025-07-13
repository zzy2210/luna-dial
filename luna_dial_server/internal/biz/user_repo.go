package biz

type UserRepo interface {
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(userID string) error
	GetUserByID(userID string) (*User, error)
	GetUserByUserName(userName string) (*User, error)
	GetUserByEmail(email string) (*User, error)
}
