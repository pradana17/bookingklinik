package repository

import (
	"booking-klinik/model"

	"gorm.io/gorm"
)

type ServiceRepository interface {
	CreateService(service *model.Service) error
	GetAllServices(limit, offset int) ([]model.Service, error)
	GetServiceById(id uint) (*model.Service, error)
	UpdateService(serviceID uint, service model.Service) (*model.Service, error)
	DeleteService(serviceID uint, deletedBy uint) error
}

type ServiceRepositoryImpl struct {
	DB *gorm.DB
}

func (r *ServiceRepositoryImpl) CreateService(service *model.Service) error {
	tx := r.DB.Begin()
	if err := tx.Create(service).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *ServiceRepositoryImpl) GetAllServices(limit, offset int) ([]model.Service, error) {
	var services []model.Service
	if err := r.DB.Limit(limit).Offset(offset).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (r *ServiceRepositoryImpl) GetServiceById(id uint) (*model.Service, error) {
	var service model.Service
	if err := r.DB.First(&service, id).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *ServiceRepositoryImpl) UpdateService(serviceID uint, service model.Service) (*model.Service, error) {
	tx := r.DB.Begin()

	var existingService model.Service
	if err := r.DB.First(&existingService, serviceID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	existingService.Name = service.Name
	existingService.Description = service.Description
	existingService.Price = service.Price
	existingService.DurationMinutes = service.DurationMinutes
	existingService.IsActive = service.IsActive
	existingService.UpdatedBy = service.UpdatedBy

	if err := r.DB.Save(&existingService).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &existingService, nil
}

func (r *ServiceRepositoryImpl) DeleteService(serviceID uint, deletedBy uint) error {
	tx := r.DB.Begin()

	var service model.Service
	if err := r.DB.First(&service, serviceID).Error; err != nil {
		tx.Rollback()
		return err
	}

	service.UpdatedBy = deletedBy

	if err := r.DB.Save(&service).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&model.Service{}, serviceID).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
