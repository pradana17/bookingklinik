package repository

import (
	"booking-klinik/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserById(id uint) (*model.User, error)
	UpdatePassword(userID uint, newPassword string) error
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func (r *UserRepositoryImpl) CreateUser(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserById(id uint) (*model.User, error) {
	var user model.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdatePassword(userID uint, newPassword string) error {
	if err := r.DB.Model(&model.User{}).Where("id = ?", userID).Update("password", newPassword).Error; err != nil {
		return err
	}
	return nil
}
