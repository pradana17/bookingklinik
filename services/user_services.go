package services

import (
	"booking-klinik/model"
	"booking-klinik/repository"
	"booking-klinik/utils"
	"errors"
)

type UserService interface {
	RegisterUser(user *model.User) (*model.User, error)
	LoginUser(email, password string) (string, error)
	UpdatePassword(userID uint, OldPassword, newPassword string) error
	GetUserById(id uint) (*model.User, error)
}

type UserServicesImpl struct {
	UserRepository repository.UserRepository
}

func (s *UserServicesImpl) RegisterUser(user *model.User) (*model.User, error) {
	if user.Role == "" {
		user.Role = "patient"
	}
	existingUser, _ := s.UserRepository.GetUserByEmail(user.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	if user.Password == "" {
		return nil, errors.New("password is required")
	}

	if len(user.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword
	if err := s.UserRepository.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServicesImpl) LoginUser(email, password string) (string, error) {
	user, err := s.UserRepository.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(*user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserServicesImpl) UpdatePassword(userID uint, OldPassword, newPassword string) error {
	user, err := s.UserRepository.GetUserById(userID)
	if err != nil {
		return err
	}

	if !utils.CheckPassword(OldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	if newPassword == OldPassword {
		return errors.New("new password cannot be the same as the old password")
	}

	if len(newPassword) < 6 {
		return errors.New("new password must be at least 6 characters long")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.UserRepository.UpdatePassword(userID, hashedPassword); err != nil {
		return err
	}

	return nil
}

func (s *UserServicesImpl) GetUserById(id uint) (*model.User, error) {
	user, err := s.UserRepository.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
