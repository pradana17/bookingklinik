package services

import (
	"booking-klinik/model"
	"booking-klinik/repository"
	"errors"
	"fmt"
)

type DoctorServices interface {
	CreateDoctor(doctor *model.Doctor) (*model.Doctor, error)
	GetAllDoctors(limit, offset int) ([]model.Doctor, error)
	GetDoctorById(id uint) (*model.Doctor, error)
	UpdateDoctor(doctorID uint, doctor model.Doctor) (*model.Doctor, error)
	DeleteDoctor(doctorID uint, userID uint) error
}

type DoctorServicesImpl struct {
	DoctorRepository repository.DoctorRepository
	UserRepository   repository.UserRepository
}

func (s *DoctorServicesImpl) CreateDoctor(doctor *model.Doctor) (*model.Doctor, error) {
	user, err := s.UserRepository.GetUserById(doctor.UserId)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if doctor.Specialization == "" {
		return nil, errors.New("specialization is required")
	}

	if user.Role != "doctor" {
		return nil, errors.New("user is not a doctor")
	}

	fmt.Println("UserRepository: ", s.UserRepository)
	if err := s.DoctorRepository.CreateDoctor(doctor); err != nil {
		return nil, err
	}

	return doctor, nil
}

func (s *DoctorServicesImpl) GetAllDoctors(limit, offset int) ([]model.Doctor, error) {
	doctors, err := s.DoctorRepository.GetAllDoctors(limit, offset)
	if err != nil {
		return nil, err
	}
	/*
		pagination.TotalRows = totalRows
		pagination.TotalPages = (totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit)*/
	return doctors, nil
}

func (s *DoctorServicesImpl) GetDoctorById(id uint) (*model.Doctor, error) {
	doctor, err := s.DoctorRepository.GetDoctorById(id)
	if err != nil {
		return nil, err
	}
	return doctor, nil
}

func (s *DoctorServicesImpl) UpdateDoctor(doctorID uint, doctor model.Doctor) (*model.Doctor, error) {

	updatedDoctor, err := s.DoctorRepository.UpdateDoctor(doctorID, doctor)
	if err != nil {
		return nil, err
	}

	return updatedDoctor, nil
}

func (s *DoctorServicesImpl) DeleteDoctor(doctorID uint, userID uint) error {

	err := s.DoctorRepository.DeleteDoctor(doctorID, userID)
	if err != nil {
		return err
	}
	return nil
}
