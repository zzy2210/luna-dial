package biz

import (
	"context"
	"net/mail"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/tjfoc/gmsm/sm3"
)

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
	UserRepo UserRepo
}

func NewUserUsecase(userRepo UserRepo) *UserUsecase {
	return &UserUsecase{
		UserRepo: userRepo,
	}
}

// 创建用户
// email 非必填
func (uc *UserUsecase) CreateUser(ctx context.Context, param CreateUserParam) (*User, error) {
	if param.UserName == "" {
		return nil, ErrUserNameEmpty
	}
	if !isValidUserName(param.UserName) {
		return nil, ErrUserNameInvalid // 用户名格式不合法
	}
	if param.Password == "" {
		return nil, ErrPasswordEmpty
	}
	if !isStrongPassword(param.Password) {
		return nil, ErrPasswordTooWeak
	}
	// 验证 email
	if param.Email != "" {
		if !validEmail(param.Email) {
			return nil, ErrEmailInvalid
		}
	}

	if param.Name == "" {
		return nil, ErrNameEmpty // 姓名不能为空
	}
	if !(isValidName(param.Name)) {
		return nil, ErrUserNameInvalid // 用户名格式不合法
	}

	// 查重 username
	existingUser, err := uc.UserRepo.GetUserByUserName(ctx, param.UserName)
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserNameDuplicate
	}

	id := uuid.NewString()
	user := &User{
		ID:        id,
		UserName:  param.UserName,
		Name:      param.Name,
		Email:     param.Email,
		Password:  string(hashPassword(param.Password)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.UserRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// 编辑用户
// 禁止修改 userName,ID
func (uc *UserUsecase) UpdateUser(ctx context.Context, param UpdateUserParam) (*User, error) {
	if param.UserID == "" {
		return nil, ErrUserIDEmpty
	}
	if !isValidUUID(param.UserID) {
		return nil, ErrUserIDInvalid
	}
	// 根据 用户ID 查用户
	user, err := uc.UserRepo.GetUserByID(ctx, param.UserID)
	if err != nil {
		return nil, err
	}
	if param.Email != nil && *param.Email != "" {
		if !validEmail(*param.Email) {
			return nil, ErrEmailInvalid // 邮箱格式不合法
		}
		user.Email = *param.Email
	}
	if param.Name != nil && *param.Name != "" {
		user.Name = *param.Name
	}
	if param.Password != nil && *param.Password != "" {
		if !isStrongPassword(*param.Password) {
			return nil, ErrPasswordTooWeak // 密码强度不足
		}
		// 使用国密进行密码hash
		user.Password = string(hashPassword(*param.Password))
	}
	if !(isValidName(*param.Name)) {
		return nil, ErrUserNameInvalid // 用户名格式不合法
	}

	user.UpdatedAt = time.Now()

	err = uc.UserRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// 删除用户
func (uc *UserUsecase) DeleteUser(ctx context.Context, param DeleteUserParam) error {
	if param.UserID == "" {
		return ErrUserIDEmpty
	}
	// 使用uuid库校验uuid
	if !isValidUUID(param.UserID) {
		return ErrUserIDInvalid
	}
	// 根据 用户ID 查用户
	user, err := uc.UserRepo.GetUserByID(ctx, param.UserID)
	if err != nil {
		if err == ErrUserNotFound {
			return ErrUserNotFound // 用户不存在
		}
		return ErrNoPermission
	}

	return uc.UserRepo.DeleteUser(ctx, user.ID)
}

// 获取用户
// 不会做脱敏处理，业务层需要脱敏
// 但是获取的密码已经被sm3 hash
func (uc *UserUsecase) GetUser(ctx context.Context, param GetUserParam) (*User, error) {

	if param.UserID == "" {
		return nil, ErrUserIDEmpty
	}
	if !isValidUUID(param.UserID) {
		return nil, ErrUserIDInvalid
	}
	// 根据 用户ID 查用户
	user, err := uc.UserRepo.GetUserByID(ctx, param.UserID)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrNoPermission
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// 用户登录
func (uc *UserUsecase) UserLogin(ctx context.Context, param UserLoginParam) (*User, error) {
	if param.UserName == "" || param.Password == "" {
		return nil, ErrInvalidInput // 参数不合法
	}

	user, err := uc.UserRepo.GetUserByUserName(ctx, param.UserName)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, ErrUserNotFound // 用户不存在
		}
		return nil, err // 其他错误
	}
	if user.Password != string(hashPassword(param.Password)) {
		return nil, ErrPasswordIncorrect // 密码错误
	}

	return user, nil
}

func hashPassword(password string) []byte {
	return sm3.Sm3Sum([]byte(password)) // 使用国密SM3算法进行密码哈希
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false // 密码长度至少8位
	}
	// 可以添加更多复杂度检查，比如包含数字、特殊字符等
	return true
}

// 校验用户姓名
func isValidName(name string) bool {
	if len(name) < 2 || len(name) > 20 {
		return false // 用户名长度必须在2到20个字符之间
	}
	for _, r := range name {
		if !(r >= 'a' && r <= 'z') &&
			!(r >= 'A' && r <= 'Z') &&
			!(r >= '0' && r <= '9') &&
			r != '_' && r != '-' &&
			!unicode.Is(unicode.Han, r) { // 允许中文
			return false // 只允许字母、数字、下划线、连字符和中文
		}
	}
	return true
}

// 校验用户名
func isValidUserName(userName string) bool {
	if len(userName) < 2 || len(userName) > 20 {
		return false // 用户名长度必须在2到20个字符之间
	}
	for _, r := range userName {
		if !(r >= 'a' && r <= 'z') &&
			!(r >= 'A' && r <= 'Z') &&
			!(r >= '0' && r <= '9') &&
			r != '_' && r != '-' {
			return false // 只允许字母、数字、下划线和连字符
		}
	}
	return true
}

func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
