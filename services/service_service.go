package services

import (
	"booking-klinik/model"
	"booking-klinik/repository"
	"errors"
)

type ServiceService interface {
	CreateService(service model.Service) (*model.Service, error)
	GetAllServices(limit, offset int) ([]model.Service, error)
	GetServiceById(id uint) (*model.Service, error)
	UpdateService(serviceID uint, service model.Service) (*model.Service, error)
	DeleteService(serviceID uint, deletedBy uint) error
}

type ServiceServiceImpl struct {
	ServiceRepository repository.ServiceRepository
}

func (s *ServiceServiceImpl) CreateService(service model.Service) (*model.Service, error) {

	if service.Name == "" {
		return nil, errors.New("service name is required")
	}

	if service.Price <= 0 || service.DurationMinutes <= 0 {
		return nil, errors.New("service price and duration minutes must be greater than 0")
	}

	if err := s.ServiceRepository.CreateService(&service); err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *ServiceServiceImpl) GetAllServices(limit, offset int) ([]model.Service, error) {
	if services, err := s.ServiceRepository.GetAllServices(limit, offset); err != nil {
		return nil, err
	} else {
		return services, nil
	}
}

func (s *ServiceServiceImpl) GetServiceById(id uint) (*model.Service, error) {
	if service, err := s.ServiceRepository.GetServiceById(id); err != nil {
		return nil, err
	} else {
		return service, nil
	}
}

func (s *ServiceServiceImpl) UpdateService(serviceID uint, service model.Service) (*model.Service, error) {
	if service.Name == "" {
		return nil, errors.New("service name is required")
	}

	if service.Price <= 0 || service.DurationMinutes <= 0 {
		return nil, errors.New("service price and duration minutes must be greater than 0")
	}

	if updatedService, err := s.ServiceRepository.UpdateService(serviceID, service); err != nil {
		return nil, err
	} else {
		return updatedService, nil
	}
}

func (s *ServiceServiceImpl) DeleteService(serviceID uint, deletedBy uint) error {
	if err := s.ServiceRepository.DeleteService(serviceID, deletedBy); err != nil {
		return err
	}
	return nil
}
