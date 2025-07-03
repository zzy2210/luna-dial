package service

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"okr-web/ent"
	"okr-web/internal/repository"
	"okr-web/internal/types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userRepo       repository.UserRepository
	jwtSecret      string
	jwtExpiryHours int
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository, jwtSecret string, jwtExpiryHours int) UserService {
	return &UserServiceImpl{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

// Register 用户注册
func (s *UserServiceImpl) Register(ctx context.Context, req RegisterRequest) (*ent.User, error) {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, &types.AppError{
			Code:    400,
			Message: "用户名已存在",
			Type:    "USERNAME_EXISTS",
		}
	}

	// 检查邮箱是否已存在
	existingUser, err = s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, &types.AppError{
			Code:    400,
			Message: "邮箱已存在",
			Type:    "EMAIL_EXISTS",
		}
	}

	// 密码加密
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "密码加密失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 创建用户
	user, err := s.userRepo.Create(ctx, func(create *ent.UserCreate) *ent.UserCreate {
		return create.
			SetUsername(req.Username).
			SetEmail(req.Email).
			SetPassword(hashedPassword)
	})
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "用户创建失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return user, nil
}

// Login 用户登录
func (s *UserServiceImpl) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// 根据用户名查找用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, &types.AppError{
			Code:    401,
			Message: "用户名或密码错误",
			Type:    "INVALID_CREDENTIALS",
		}
	}

	// 验证密码
	valid, err := s.verifyPassword(req.Password, user.Password)
	if err != nil || !valid {
		return nil, &types.AppError{
			Code:    401,
			Message: "用户名或密码错误",
			Type:    "INVALID_CREDENTIALS",
		}
	}

	// 生成JWT Token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "Token生成失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*ent.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    404,
			Message: "用户不存在",
			Type:    "USER_NOT_FOUND",
		}
	}
	return user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserServiceImpl) GetUserByUsername(ctx context.Context, username string) (*ent.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, &types.AppError{
			Code:    404,
			Message: "用户不存在",
			Type:    "USER_NOT_FOUND",
		}
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (s *UserServiceImpl) UpdateUser(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) (*ent.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    404,
			Message: "用户不存在",
			Type:    "USER_NOT_FOUND",
		}
	}

	updateData := map[string]interface{}{}

	if req.Email != nil {
		// 检查邮箱是否已被其他用户使用
		existingUser, err := s.userRepo.GetByEmail(ctx, *req.Email)
		if err == nil && existingUser != nil && existingUser.ID != userID {
			return nil, &types.AppError{
				Code:    400,
				Message: "邮箱已被使用",
				Type:    "EMAIL_EXISTS",
			}
		}
		updateData["email"] = *req.Email
	}

	if req.Password != nil {
		hashedPassword, err := s.hashPassword(*req.Password)
		if err != nil {
			return nil, &types.AppError{
				Code:    500,
				Message: "密码加密失败",
				Type:    "INTERNAL_SERVER_ERROR",
			}
		}
		updateData["password"] = hashedPassword
	}

	if len(updateData) == 0 {
		return user, nil
	}

	updatedUser, err := s.userRepo.Update(ctx, userID, func(update *ent.UserUpdateOne) *ent.UserUpdateOne {
		if email, ok := updateData["email"]; ok {
			update = update.SetEmail(email.(string))
		}
		if password, ok := updateData["password"]; ok {
			update = update.SetPassword(password.(string))
		}
		return update
	})
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "用户更新失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return updatedUser, nil
}

// JWT相关方法
type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

func (s *UserServiceImpl) generateToken(user *ent.User) (string, error) {
	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.jwtExpiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *UserServiceImpl) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// 密码加密相关方法
type PasswordConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

var defaultPasswordConfig = PasswordConfig{
	Time:    1,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
}

func (s *UserServiceImpl) hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, defaultPasswordConfig.Time, defaultPasswordConfig.Memory, defaultPasswordConfig.Threads, defaultPasswordConfig.KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	return fmt.Sprintf(format, argon2.Version, defaultPasswordConfig.Memory, defaultPasswordConfig.Time, defaultPasswordConfig.Threads, b64Salt, b64Hash), nil
}

func (s *UserServiceImpl) verifyPassword(password, hash string) (bool, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, errors.New("incompatible argon2 version")
	}

	var memory, time uint32
	var threads uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	comparisonHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(decodedHash)))

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
