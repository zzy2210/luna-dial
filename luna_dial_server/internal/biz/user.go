package biz

import "time"

type User struct {
	ID        string    `json:"id"`
	UserName  string    `json:"user_name"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 创建用户参数
type CreateUserParam struct {
	UserName string
	Name     string
	Email    string
	Password string
}

// 编辑用户参数
type UpdateUserParam struct {
	UserID   string
	UserName *string
	Name     *string
	Email    *string
	Password *string
}

// 删除用户参数
type DeleteUserParam struct {
	UserID string
}

// 获取用户参数
type GetUserParam struct {
	UserID string
}

// 用户登录参数
type UserLoginParam struct {
	UserName string
	Password string
}

type UserUsecase struct {
	UserRepo *UserRepo
}

func NewUserUsecase(userRepo *UserRepo) *UserUsecase {
	return &UserUsecase{
		UserRepo: userRepo,
	}
}

// 创建用户
func (uc *UserUsecase) CreateUser(param CreateUserParam) (*User, error) {
	return nil, ErrNoPermission // TODO: 实现
}

// 编辑用户
func (uc *UserUsecase) UpdateUser(param UpdateUserParam) (*User, error) {
	return nil, ErrNoPermission // TODO: 实现
}

// 删除用户
func (uc *UserUsecase) DeleteUser(param DeleteUserParam) error {
	return ErrNoPermission // TODO: 实现
}

// 获取用户
func (uc *UserUsecase) GetUser(param GetUserParam) (*User, error) {
	return nil, ErrNoPermission // TODO: 实现
}

// 用户登录
func (uc *UserUsecase) UserLogin(param UserLoginParam) (*User, error) {
	return nil, ErrNoPermission // TODO: 实现
}
