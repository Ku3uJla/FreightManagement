package repository

import (
	"auth-service/internal/repository/model"
	"context"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Create(ctx context.Context, user *model.Auth) error {
	return r.db.WithContext(ctx).Model(&model.Auth{}).Create(user).Error
}

func (r *AuthRepository) GetByLogin(ctx context.Context, login string) (*model.Auth, error) {
	var user model.Auth
	err := r.db.WithContext(ctx).Model(&model.Auth{}).Where("login = ?", login).First(&user).Error
	return &user, err
}

func (r *AuthRepository) GetByEmail(ctx context.Context, email string) (*model.Auth, error) {
	var user model.Auth
	err := r.db.WithContext(ctx).Model(&model.Auth{}).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *AuthRepository) ExistsByEmail(ctx context.Context, email string) bool {
	var count int64
	r.db.WithContext(ctx).Model(&model.Auth{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (r *AuthRepository) GetByID(ctx context.Context, id string) (*model.Auth, error) {
	var user model.Auth
	err := r.db.WithContext(ctx).Model(&model.Auth{}).Where("id = ?", id).Find(&user).Error
	return &user, err
}
