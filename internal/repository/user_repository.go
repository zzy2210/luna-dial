package repository

import (
	"context"
	"fmt"

	"okr-web/ent"
	"okr-web/ent/user"

	"github.com/google/uuid"
)

// userRepository 用户Repository实现
type userRepository struct {
	client *ent.Client
}

// NewUserRepository 创建新的用户Repository
func NewUserRepository(client *ent.Client) UserRepository {
	return &userRepository{client: client}
}

// Create 创建新用户
func (r *userRepository) Create(ctx context.Context, builder func(*ent.UserCreate) *ent.UserCreate) (*ent.User, error) {
	return builder(r.client.User.Create()).Save(ctx)
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	u, err := r.client.User.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return u, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*ent.User, error) {
	u, err := r.client.User.
		Query().
		Where(user.Username(username)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return u, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	u, err := r.client.User.
		Query().
		Where(user.Email(email)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return u, nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, id uuid.UUID, updater func(*ent.UserUpdateOne) *ent.UserUpdateOne) (*ent.User, error) {
	updateOne := r.client.User.UpdateOneID(id)
	updateOne = updater(updateOne)

	u, err := updateOne.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return u, nil
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
