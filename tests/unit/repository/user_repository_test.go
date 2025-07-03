package repository

import (
	"context"
	"testing"

	"okr-web/ent"
	"okr-web/ent/enttest"
	"okr-web/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func TestUserRepository_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := repository.NewUserRepository(client)

	ctx := context.Background()
	userCreate := client.User.Create().
		SetUsername("testuser").
		SetPassword("hashedpassword").
		SetEmail("test@example.com")

	user, err := repo.Create(ctx, userCreate)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestUserRepository_GetByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := repository.NewUserRepository(client)
	ctx := context.Background()

	// 创建测试用户
	createdUser := client.User.Create().
		SetUsername("testuser").
		SetPassword("hashedpassword").
		SetEmail("test@example.com").
		SaveX(ctx)

	// 测试获取用户
	user, err := repo.GetByID(ctx, createdUser.ID)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, createdUser.ID, user.ID)
	assert.Equal(t, "testuser", user.Username)

	// 测试获取不存在的用户
	_, err = repo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_GetByUsername(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := repository.NewUserRepository(client)
	ctx := context.Background()

	// 创建测试用户
	createdUser := client.User.Create().
		SetUsername("testuser").
		SetPassword("hashedpassword").
		SetEmail("test@example.com").
		SaveX(ctx)

	// 测试获取用户
	user, err := repo.GetByUsername(ctx, "testuser")
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, createdUser.ID, user.ID)
	assert.Equal(t, "testuser", user.Username)

	// 测试获取不存在的用户
	_, err = repo.GetByUsername(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := repository.NewUserRepository(client)
	ctx := context.Background()

	// 创建测试用户
	createdUser := client.User.Create().
		SetUsername("testuser").
		SetPassword("hashedpassword").
		SetEmail("test@example.com").
		SaveX(ctx)

	// 测试更新用户
	updatedUser, err := repo.Update(ctx, createdUser.ID, func(u *ent.UserUpdateOne) *ent.UserUpdateOne {
		return u.SetEmail("newemail@example.com")
	})
	require.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, "newemail@example.com", updatedUser.Email)
	assert.Equal(t, "testuser", updatedUser.Username) // username 保持不变

	// 测试更新不存在的用户
	_, err = repo.Update(ctx, uuid.New(), func(u *ent.UserUpdateOne) *ent.UserUpdateOne {
		return u.SetEmail("test@example.com")
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := repository.NewUserRepository(client)
	ctx := context.Background()

	// 创建测试用户
	createdUser := client.User.Create().
		SetUsername("testuser").
		SetPassword("hashedpassword").
		SetEmail("test@example.com").
		SaveX(ctx)

	// 测试删除用户
	err := repo.Delete(ctx, createdUser.ID)
	require.NoError(t, err)

	// 验证用户已删除
	_, err = repo.GetByID(ctx, createdUser.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	// 测试删除不存在的用户
	err = repo.Delete(ctx, uuid.New())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}
