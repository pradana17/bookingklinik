package repository

import (
	"booking-klinik/model"

	"gorm.io/gorm"
)

type DoctorScheduleRepository interface {
	CreateDoctorSchedule(doctorSchedule *model.DoctorSchedule) error
	GetDoctorSchedulesByDoctorId(doctorId uint) ([]model.DoctorSchedule, error)
	GetDoctorSchedulesById(scheduleId uint) (*model.DoctorSchedule, error)
	GetAllDoctorSchedules(limit, offset int) ([]model.DoctorSchedule, error)
	UpdateDoctorSchedule(doctorSchedule *model.DoctorSchedule) error
	DeleteDoctorSchedule(scheduleId uint, userID uint) error
}

type DoctorScheduleRepositoryImpl struct {
	DB *gorm.DB
}

func (r *DoctorScheduleRepositoryImpl) CreateDoctorSchedule(doctorSchedule *model.DoctorSchedule) error {
	tx := r.DB.Begin()
	if err := tx.Create(doctorSchedule).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *DoctorScheduleRepositoryImpl) GetDoctorSchedulesByDoctorId(doctorId uint) ([]model.DoctorSchedule, error) {
	var doctorSchedules []model.DoctorSchedule
	if err := r.DB.Where("doctor_id = ?", doctorId).Find(&doctorSchedules).Error; err != nil {
		return nil, err
	}
	return doctorSchedules, nil
}

func (r *DoctorScheduleRepositoryImpl) GetDoctorSchedulesById(scheduleId uint) (*model.DoctorSchedule, error) {
	var doctorSchedule model.DoctorSchedule
	if err := r.DB.First(&doctorSchedule, scheduleId).Error; err != nil {
		return nil, err
	}
	return &doctorSchedule, nil
}

func (r *DoctorScheduleRepositoryImpl) UpdateDoctorSchedule(doctorSchedule *model.DoctorSchedule) error {
	var existingDoctorSchedule model.DoctorSchedule
	if err := r.DB.First(&existingDoctorSchedule, doctorSchedule.ID).Error; err != nil {
		return err
	}

	existingDoctorSchedule.Date = doctorSchedule.Date
	existingDoctorSchedule.StartTime = doctorSchedule.StartTime
	existingDoctorSchedule.EndTime = doctorSchedule.EndTime
	existingDoctorSchedule.UpdatedBy = doctorSchedule.Doctor.User.ID

	tx := r.DB.Begin()

	if err := tx.Save(&existingDoctorSchedule).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *DoctorScheduleRepositoryImpl) DeleteDoctorSchedule(scheduleId uint, userID uint) error {
	var doctorSchedule model.DoctorSchedule
	tx := r.DB.Begin()

	if err := tx.First(&doctorSchedule, scheduleId).Error; err != nil {
		tx.Rollback()
		return err
	}

	doctorSchedule.UpdatedBy = userID

	if err := tx.Save(&doctorSchedule).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&doctorSchedule).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *DoctorScheduleRepositoryImpl) GetAllDoctorSchedules(limit, offset int) ([]model.DoctorSchedule, error) {
	var doctorSchedules []model.DoctorSchedule
	if err := r.DB.Limit(limit).Offset(offset).Preload("Doctor").Find(&doctorSchedules).Error; err != nil {
		return nil, err
	}
	return doctorSchedules, nil
}
