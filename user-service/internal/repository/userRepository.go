package repository

import (
	"context"
	"user-service/internal/repository/model"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	ExistsByEmail(ctx context.Context, email string) bool
	GetByID(ctx context.Context, id string) (*model.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	err := r.db.Where("login = ?", login).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) bool {
	var count int64
	r.db.Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id_user = ?", id).Find(&user).Error
	return &user, err
}
