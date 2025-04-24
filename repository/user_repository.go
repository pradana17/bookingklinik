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
	tx := r.DB.Begin()
	defer tx.Commit()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
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

	tx := r.DB.Begin()
	defer tx.Commit()
	var user model.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return err
	}
	user.Password = newPassword
	user.UpdatedBy = userID
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
