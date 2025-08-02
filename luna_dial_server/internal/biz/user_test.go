package biz

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepo 使用 testify/mock 框架生成的 mock
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, userID string) (*User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepo) GetUserByUserName(ctx context.Context, userName string) (*User, error) {
	args := m.Called(ctx, userName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

// 辅助函数：创建测试用的 UserUsecase 和 MockUserRepo
func setupTest() (*UserUsecase, *MockUserRepo) {
	mockRepo := new(MockUserRepo)
	// 将 *MockUserRepo 转换为 UserRepo 接口
	var userRepo UserRepo = mockRepo
	usecase := NewUserUsecase(userRepo)
	return usecase, mockRepo
}

// TestCreateUser 测试创建用户功能
func TestCreateUser(t *testing.T) {
	t.Run("成功创建用户", func(t *testing.T) {
		usecase, mockRepo := setupTest()

		param := CreateUserParam{
			UserName: "testuser",
			Name:     "测试用户",
			Email:    "test@example.com",
			Password: "securepassword123",
		}

		// 设置 mock 期望
		mockRepo.On("GetUserByUserName", mock.Anything, param.UserName).Return(nil, ErrUserNotFound)
		mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*biz.User")).Return(nil)

		// 执行测试
		user, err := usecase.CreateUser(context.Background(), param)

		// 业务逻辑实现后，err 应该为 nil，user 不为 nil
		assert.Nil(t, err)
		assert.NotNil(t, user)
		if user != nil {
			assert.Equal(t, param.UserName, user.Username)
			assert.Equal(t, param.Name, user.Name)
			assert.Equal(t, param.Email, user.Email)
			assert.NotEmpty(t, user.ID)
			assert.NotEmpty(t, user.Password) // 密码应该被哈希，不应是明文
			assert.NotEqual(t, param.Password, user.Password)
		}

		// 验证 mock 是否被正确调用
		mockRepo.AssertExpectations(t)
	})

	t.Run("参数验证-空用户名", func(t *testing.T) {
		usecase, _ := setupTest()

		param := CreateUserParam{
			UserName: "",
			Name:     "测试用户",
			Email:    "test@example.com",
			Password: "securepassword123",
		}

		user, err := usecase.CreateUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)

		// 业务逻辑实现后应该返回 ErrUserNameEmpty，现在会失败
		assert.Equal(t, ErrUserNameEmpty, err, "应该返回 ErrUserNameEmpty 错误")
	})

	t.Run("参数验证-无效邮箱格式", func(t *testing.T) {
		usecase, _ := setupTest()

		param := CreateUserParam{
			UserName: "testuser",
			Name:     "测试用户",
			Email:    "invalid-email",
			Password: "securepassword123",
		}

		user, err := usecase.CreateUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)

		// 业务逻辑实现后应该返回 ErrEmailInvalid，当前返回 ErrNoPermission 会导致测试失败
		assert.Equal(t, ErrEmailInvalid, err, "应该返回 ErrEmailInvalid 错误")
	})

	t.Run("参数验证-弱密码", func(t *testing.T) {
		usecase, _ := setupTest()

		param := CreateUserParam{
			UserName: "testuser",
			Name:     "测试用户",
			Email:    "test@example.com",
			Password: "123", // 弱密码
		}

		user, err := usecase.CreateUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)

		// 业务逻辑实现后应该返回 ErrPasswordTooWeak，当前返回 ErrNoPermission 会导致测试失败
		assert.Equal(t, ErrPasswordTooWeak, err, "应该返回 ErrPasswordTooWeak 错误")
	})

	t.Run("用户名重复", func(t *testing.T) {
		usecase, mockRepo := setupTest()

		param := CreateUserParam{
			UserName: "existinguser",
			Name:     "测试用户",
			Email:    "test@example.com",
			Password: "securepassword123",
		}

		// 模拟用户名已存在
		existingUser := &User{
			ID:       "existing-id",
			Username: "existinguser",
			Email:    "existing@example.com",
		}
		mockRepo.On("GetUserByUserName", mock.Anything, param.UserName).Return(existingUser, nil)

		user, err := usecase.CreateUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrUserNameDuplicate, err, "应该返回 ErrUserNameDuplicate 错误")
	})

	t.Run("参数验证-空姓名", func(t *testing.T) {
		usecase, _ := setupTest()

		param := CreateUserParam{
			UserName: "testuser",
			Name:     "", // 空姓名
			Email:    "test@example.com",
			Password: "securepassword123",
		}

		user, err := usecase.CreateUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)

		// 业务逻辑实现后应该返回 ErrNameEmpty，现在会失败
		assert.Equal(t, ErrNameEmpty, err, "应该返回 ErrNameEmpty 错误")
	})

	t.Run("参数验证-用户名过长", func(t *testing.T) {
		usecase, _ := setupTest()

		param := CreateUserParam{
			UserName: strings.Repeat("a", 101), // 假设限制100字符
			Name:     "测试用户",
			Email:    "test@example.com",
			Password: "securepassword123",
		}

		user, err := usecase.CreateUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)

		// TODO: 业务逻辑实现后应该返回 ErrUserNameInvalid
		assert.Equal(t, ErrUserNameInvalid, err, "应该返回 ErrUserNameInvalid 错误")
	})

	t.Run("参数验证-用户名格式无效", func(t *testing.T) {
		usecase, _ := setupTest()

		param := CreateUserParam{
			UserName: "user@#$%", // 包含特殊字符
			Name:     "测试用户",
			Email:    "test@example.com",
			Password: "securepassword123",
		}

		user, err := usecase.CreateUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)

		// TODO: 业务逻辑实现后应该返回 ErrUserNameInvalid，现在会失败
		assert.Equal(t, ErrUserNameInvalid, err, "应该返回 ErrUserNameInvalid 错误")
	})

	t.Run("参数验证-密码为空", func(t *testing.T) {
		usecase, _ := setupTest()

		param := CreateUserParam{
			UserName: "testuser",
			Name:     "测试用户",
			Email:    "test@example.com",
			Password: "", // 空密码
		}

		user, err := usecase.CreateUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)

		// TODO: 业务逻辑实现后应该返回 ErrPasswordEmpty，现在会失败
		assert.Equal(t, ErrPasswordEmpty, err, "应该返回 ErrPasswordEmpty 错误")
	})
}

// TestUpdateUser 测试更新用户功能
func TestUpdateUser(t *testing.T) {
	t.Run("成功更新用户姓名", func(t *testing.T) {
		usecase, mockRepo := setupTest()

		newName := "更新后的姓名"
		param := UpdateUserParam{
			UserID: "550e8400e29b41d4a716446655440000",
			Name:   &newName,
		}

		// 设置 mock 期望
		existingUser := &User{
			ID:        "550e8400e29b41d4a716446655440000",
			Username:  "testuser",
			Name:      "原姓名",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		}
		mockRepo.On("GetUserByID", mock.Anything, param.UserID).Return(existingUser, nil)
		mockRepo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*biz.User")).Return(nil)

		user, err := usecase.UpdateUser(context.Background(), param)

		// 业务逻辑实现后，err 应该为 nil
		assert.Nil(t, err)
		assert.NotNil(t, user)
		if user != nil {
			assert.Equal(t, *param.Name, user.Name)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		usecase, mockRepo := setupTest()

		newName := "新姓名"
		param := UpdateUserParam{
			UserID: "6ba7b8109dad11d180b400c04fd430c8",
			Name:   &newName,
		}

		mockRepo.On("GetUserByID", mock.Anything, param.UserID).Return(nil, ErrUserNotFound)

		user, err := usecase.UpdateUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrUserNotFound, err, "应该返回 ErrUserNotFound 错误")
	})

	t.Run("参数验证-无效用户ID格式", func(t *testing.T) {
		usecase, _ := setupTest()

		newName := "新姓名"
		param := UpdateUserParam{
			UserID: "invalid-uuid", // 无效的UUID格式
			Name:   &newName,
		}

		user, err := usecase.UpdateUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrUserIDInvalid, err, "应该返回 ErrUserIDInvalid 错误")
	})
}

// TestDeleteUser 测试删除用户功能
func TestDeleteUser(t *testing.T) {
	t.Run("成功删除用户", func(t *testing.T) {
		usecase, mockRepo := setupTest()

		param := DeleteUserParam{
			UserID: "550e8400e29b41d4a716446655440000",
		}

		// 设置 mock 期望
		existingUser := &User{
			ID:       "550e8400e29b41d4a716446655440000",
			Username: "testuser",
			Email:    "test@example.com",
		}
		mockRepo.On("GetUserByID", mock.Anything, param.UserID).Return(existingUser, nil)
		mockRepo.On("DeleteUser", mock.Anything, param.UserID).Return(nil)

		err := usecase.DeleteUser(context.Background(), param)

		// 业务逻辑实现后，err 应该为 nil
		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		usecase, mockRepo := setupTest()

		param := DeleteUserParam{
			UserID: "6ba7b8109dad11d180b400c04fd430c8",
		}

		mockRepo.On("GetUserByID", mock.Anything, param.UserID).Return(nil, ErrUserNotFound)

		err := usecase.DeleteUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Equal(t, ErrUserNotFound, err, "应该返回 ErrUserNotFound 错误")
	})

	t.Run("删除空用户ID", func(t *testing.T) {
		usecase, _ := setupTest()

		param := DeleteUserParam{
			UserID: "",
		}

		err := usecase.DeleteUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Equal(t, ErrUserIDEmpty, err, "应该返回 ErrUserIDEmpty 错误")
	})

	t.Run("参数验证-无效用户ID格式", func(t *testing.T) {
		usecase, _ := setupTest()

		param := DeleteUserParam{
			UserID: "not-a-uuid", // 无效的UUID格式
		}

		err := usecase.DeleteUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Equal(t, ErrUserIDInvalid, err, "应该返回 ErrUserIDInvalid 错误")
	})
}

// TestGetUser 测试获取用户功能
func TestGetUser(t *testing.T) {
	t.Run("成功获取用户", func(t *testing.T) {
		usecase, mockRepo := setupTest()

		param := GetUserParam{
			UserID: "550e8400e29b41d4a716446655440000",
		}

		// 设置 mock 期望
		expectedUser := &User{
			ID:        "550e8400e29b41d4a716446655440000",
			Username:  "testuser",
			Name:      "测试用户",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.On("GetUserByID", mock.Anything, param.UserID).Return(expectedUser, nil)

		user, err := usecase.GetUser(context.Background(), param)

		// 业务逻辑实现后，err 应该为 nil
		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		usecase, mockRepo := setupTest()

		param := GetUserParam{
			UserID: "6ba7b8109dad11d180b400c04fd430c8",
		}

		mockRepo.On("GetUserByID", mock.Anything, param.UserID).Return(nil, ErrUserNotFound)

		user, err := usecase.GetUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrUserNotFound, err, "应该返回 ErrUserNotFound 错误")
	})

	t.Run("空用户ID", func(t *testing.T) {
		usecase, _ := setupTest()

		param := GetUserParam{
			UserID: "",
		}

		user, err := usecase.GetUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrUserIDEmpty, err, "应该返回 ErrUserIDEmpty 错误")
	})

	t.Run("参数验证-无效用户ID格式", func(t *testing.T) {
		usecase, _ := setupTest()

		param := GetUserParam{
			UserID: "invalid-format", // 无效的UUID格式
		}

		user, err := usecase.GetUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrUserIDInvalid, err, "应该返回 ErrUserIDInvalid 错误")
	})
}

// TestUserLogin 测试用户登录功能
func TestUserLogin(t *testing.T) {
	t.Run("成功登录", func(t *testing.T) {
		usecase, mockRepo := setupTest()

		param := UserLoginParam{
			Username: "testuser",
			Password: "correctpassword",
		}

		// 设置 mock 期望

		hashedPassword := string(hashPassword(param.Password))
		expectedUser := &User{
			ID:        "550e8400e29b41d4a716446655440000",
			Username:  "testuser",
			Name:      "测试用户",
			Email:     "test@example.com",
			Password:  hashedPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.On("GetUserByUserName", mock.Anything, param.Username).Return(expectedUser, nil)

		user, err := usecase.UserLogin(context.Background(), param)

		// 业务逻辑实现后，err 应该为 nil
		assert.Nil(t, err)
		assert.NotNil(t, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户名不存在", func(t *testing.T) {
		usecase, mockRepo := setupTest()

		param := UserLoginParam{
			Username: "nonexistentuser",
			Password: "anypassword",
		}

		mockRepo.On("GetUserByUserName", mock.Anything, param.Username).Return(nil, ErrUserNotFound)

		user, err := usecase.UserLogin(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrUserNotFound, err, "应该返回 ErrUserNotFound 错误")
	})

	t.Run("密码错误", func(t *testing.T) {
		usecase, mockRepo := setupTest()

		param := UserLoginParam{
			Username: "testuser",
			Password: "wrongpassword",
		}

		existingUser := &User{
			ID:       "550e8400e29b41d4a716446655440000",
			Username: "testuser",
			Password: "hashed_correctpassword",
		}
		mockRepo.On("GetUserByUserName", mock.Anything, param.Username).Return(existingUser, nil)

		user, err := usecase.UserLogin(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrPasswordIncorrect, err, "应该返回 ErrPasswordIncorrect 错误")
	})

	t.Run("空用户名", func(t *testing.T) {
		usecase, _ := setupTest()

		param := UserLoginParam{
			Username: "",
			Password: "password",
		}

		user, err := usecase.UserLogin(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrInvalidInput, err, "应该返回 ErrInvalidInput 错误")
	})

	t.Run("空密码", func(t *testing.T) {
		usecase, _ := setupTest()

		param := UserLoginParam{
			Username: "testuser",
			Password: "",
		}

		user, err := usecase.UserLogin(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrInvalidInput, err, "应该返回 ErrInvalidInput 错误")
	})

	t.Run("使用邮箱登录", func(t *testing.T) {
		usecase, mockRepo := setupTest()

		param := UserLoginParam{
			Username: "test@example.com",
			Password: "correctpassword",
		}

		// 可以尝试通过邮箱查找用户
		expectedUser := &User{
			ID:       "550e8400e29b41d4a716446655440000",
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashed_correctpassword",
		}
		mockRepo.On("GetUserByUserName", mock.Anything, param.Username).Return(nil, ErrUserNotFound)
		mockRepo.On("GetUserByEmail", mock.Anything, param.Username).Return(expectedUser, nil)

		user, err := usecase.UserLogin(context.Background(), param)

		// 假设不支持邮箱登录时，应返回用户不存在
		assert.Equal(t, ErrUserNotFound, err, "应该返回 ErrUserNotFound 错误")
		assert.Nil(t, user)

		t.Log("待实现: 考虑是否支持邮箱登录功能")
	})
}

// TestUser_Fields 测试 User 结构体字段
func TestUser_Fields(t *testing.T) {
	now := time.Now()
	user := User{
		ID:        "550e8400e29b41d4a716446655440000",
		Username:  "testuser",
		Name:      "测试用户",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, "550e8400e29b41d4a716446655440000", user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "测试用户", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.Equal(t, now, user.CreatedAt)
	assert.Equal(t, now, user.UpdatedAt)
}

// TestCreateUserParam_Fields 测试创建用户参数结构体
func TestCreateUserParam_Fields(t *testing.T) {
	param := CreateUserParam{
		UserName: "newuser",
		Name:     "新用户",
		Email:    "new@example.com",
		Password: "securepassword",
	}

	assert.Equal(t, "newuser", param.UserName)
	assert.Equal(t, "新用户", param.Name)
	assert.Equal(t, "new@example.com", param.Email)
	assert.Equal(t, "securepassword", param.Password)
}

// TestUpdateUserParam_Fields 测试更新用户参数结构体
func TestUpdateUserParam_Fields(t *testing.T) {
	newUserName := "updateduser"
	newName := "更新用户"
	newEmail := "updated@example.com"
	newPassword := "newpassword"

	param := UpdateUserParam{
		UserID:   "550e8400e29b41d4a716446655440000",
		UserName: &newUserName,
		Name:     &newName,
		Email:    &newEmail,
		Password: &newPassword,
	}

	assert.Equal(t, "550e8400e29b41d4a716446655440000", param.UserID)
	assert.NotNil(t, param.UserName)
	assert.Equal(t, newUserName, *param.UserName)
	assert.NotNil(t, param.Name)
	assert.Equal(t, newName, *param.Name)
	assert.NotNil(t, param.Email)
	assert.Equal(t, newEmail, *param.Email)
	assert.NotNil(t, param.Password)
	assert.Equal(t, newPassword, *param.Password)
}

// TestUserLoginParam_Fields 测试登录参数结构体
func TestUserLoginParam_Fields(t *testing.T) {
	param := UserLoginParam{
		Username: "loginuser",
		Password: "loginpassword",
	}

	assert.Equal(t, "loginuser", param.Username)
	assert.Equal(t, "loginpassword", param.Password)
}

// TestUser_EdgeCases 边界和安全测试
func TestUser_EdgeCases(t *testing.T) {
	t.Run("极长用户名", func(t *testing.T) {
		usecase, _ := setupTest()

		longUserName := strings.Repeat("u", 1000)
		param := CreateUserParam{
			UserName: longUserName,
			Name:     "测试用户",
			Email:    "test@example.com",
			Password: "password123",
		}
		user, err := usecase.CreateUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrUserNameInvalid, err, "应该返回 ErrUserNameTooLong 错误")

		t.Log("待实现: 用户名长度限制验证")
	})

	t.Run("特殊字符用户名", func(t *testing.T) {
		usecase, _ := setupTest()

		param := CreateUserParam{
			UserName: "user<script>alert('xss')</script>",
			Name:     "测试用户",
			Email:    "test@example.com",
			Password: "password123",
		}

		user, err := usecase.CreateUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrUserNameInvalid, err, "应该返回 ErrUserNameInvalid 错误")

		t.Log("待实现: 用户名格式验证，防止XSS攻击")
	})

	t.Run("SQL注入防护测试", func(t *testing.T) {
		usecase, _ := setupTest()

		param := CreateUserParam{
			UserName: "user'; DROP TABLE users; --",
			Name:     "测试用户",
			Email:    "test@example.com",
			Password: "password123",
		}

		user, err := usecase.CreateUser(context.Background(), param)

		assert.NotNil(t, err)
		assert.Nil(t, user)
		assert.Equal(t, ErrUserNameInvalid, err, "应该返回 ErrUserNameInvalid 错误")

		t.Log("待实现: SQL注入防护")
	})

	t.Run("密码强度测试", func(t *testing.T) {
		usecase, _ := setupTest()

		weakPasswords := []string{
			"123",
			"abc",
			"123456",
			"qwerty",
			"111111",
			"",
		}

		for _, password := range weakPasswords {
			param := CreateUserParam{
				UserName: "testuser",
				Name:     "测试用户",
				Email:    "test@example.com",
				Password: password,
			}

			user, err := usecase.CreateUser(context.Background(), param)

			assert.NotNil(t, err)
			assert.Nil(t, user)

			// 根据密码内容判断期望的错误类型
			if password == "" {
				assert.Equal(t, ErrPasswordEmpty, err, "空密码应该返回 ErrPasswordEmpty")
			} else {
				assert.Equal(t, ErrPasswordTooWeak, err, "弱密码应该返回 ErrPasswordTooShort")
			}

			t.Logf("弱密码 '%s': 期望返回对应的密码验证错误", password)
		}
	})
}

// 测试辅助函数
func TestHelperFunctions(t *testing.T) {
	t.Run("生成ID函数", func(t *testing.T) {
		id1 := generateID()
		id2 := generateID()

		assert.NotEmpty(t, id1)
		assert.NotEmpty(t, id2)
		assert.NotEqual(t, id1, id2)  // 每次生成的ID应该不同
		assert.Equal(t, 32, len(id1)) // 16字节的十六进制表示应该是32个字符
	})

	t.Run("邮箱验证函数", func(t *testing.T) {
		testCases := []struct {
			email string
			valid bool
		}{
			{"test@example.com", true},
			{"user+tag@example.com", true},
			{"user.name@example.co.uk", true},
			{"invalid-email", false},
			{"@example.com", false},
			{"test@", false},
			{"test.example.com", false},
			{"", false},
		}

		for _, tc := range testCases {
			result := validEmail(tc.email)
			assert.Equal(t, tc.valid, result, "邮箱 %s 的验证结果不正确", tc.email)
		}
	})

	t.Run("密码强度验证函数", func(t *testing.T) {
		testCases := []struct {
			password string
			strong   bool
		}{
			{"password123", true},
			{"securepassword", true},
			{"12345678", true}, // 长度足够但简单
			{"123", false},
			{"abc", false},
			{"", false},
		}

		for _, tc := range testCases {
			result := isStrongPassword(tc.password)
			assert.Equal(t, tc.strong, result, "密码 '%s' 的强度验证结果不正确", tc.password)
		}
	})

	t.Run("UUID验证函数", func(t *testing.T) {
		testCases := []struct {
			uuid  string
			valid bool
		}{
			{"550e8400-e29b-41d4-a716-446655440000", true},  // 标准UUID v4
			{"6ba7b810-9dad-11d1-80b4-00c04fd430c8", true},  // 标准UUID v1
			{"6ba7b811-9dad-11d1-80b4-00c04fd430c8", true},  // 另一个有效UUID
			{"550e8400-e29b-41d4-a716-44665544000", false},  // 缺少一位
			{"550e8400-e29b-41d4-a716-44665544000g", false}, // 包含非十六进制字符
			{"550e8400e29b41d4a716446655440000", true},      // 缺少连字符
			{"550e8400-e29b-41d4-a716", false},              // 格式不完整
			{"", false},                                     // 空字符串
			{"not-a-uuid", false},                           // 完全不是UUID格式
		}

		for _, tc := range testCases {
			result := isValidUUID(tc.uuid)
			assert.Equal(t, tc.valid, result, "UUID %s 的验证结果不正确", tc.uuid)
		}
	})
}
