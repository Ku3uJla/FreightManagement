package repository

import (
	"context"
	"user-service/internal/repository/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	UpdateRole(ctx context.Context, id, status int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Find(&user).Error
	return &user, err
}

func (r *userRepository) UpdateRole(ctx context.Context, id, status int) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}
