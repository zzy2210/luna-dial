package biz

import "context"

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userID string) error
	GetUserByID(ctx context.Context, userID string) (*User, error)
	GetUserByUserName(ctx context.Context, userName string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
