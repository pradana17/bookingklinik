package repository

import (
	"booking-klinik/model"

	"gorm.io/gorm"
)

type DoctorRepository interface {
	CreateDoctor(doctor *model.Doctor) error
	GetAllDoctors(limit, offset int) ([]model.Doctor, error)
	GetDoctorById(id uint) (*model.Doctor, error)
	UpdateDoctor(doctorID uint, doctor model.Doctor) (*model.Doctor, error)
	DeleteDoctor(doctorID uint) error
}

type DoctorRepositoryImpl struct {
	DB *gorm.DB
}

func (r *DoctorRepositoryImpl) CreateDoctor(doctor *model.Doctor) error {
	tx := r.DB.Begin()
	if err := tx.Create(doctor).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// GetAllDoctors gets all doctors with the given limit and offset.
//
// This function will return a list of doctors with the given limit and offset.
// If the query is not successful, it will return an error.
func (r *DoctorRepositoryImpl) GetAllDoctors(limit, offset int) ([]model.Doctor, error) {
	var doctors []model.Doctor
	if err := r.DB.Limit(limit).Offset(offset).Preload("User").Find(&doctors).Error; err != nil {
		return nil, err
	}
	return doctors, nil
}

// GetDoctorById gets a doctor by given ID.
//
// This function will return a doctor that matches the given ID.
// If the doctor is not found, it will return an error.
func (r *DoctorRepositoryImpl) GetDoctorById(id uint) (*model.Doctor, error) {
	var doctor model.Doctor
	if err := r.DB.Preload("User").First(&doctor, id).Error; err != nil {
		return nil, err
	}
	return &doctor, nil
}

// UpdateDoctor updates the specialization of a doctor with the given doctorID.
//
// This function retrieves the existing doctor from the database by doctorID,
// updates the specialization, and saves the changes. If any errors occur during
// the process, it rolls back the transaction and returns an error. Upon
// successful update, it commits the transaction and returns the updated doctor.
//
// Parameters:
//  - doctorID: the ID of the doctor to be updated.
//  - doctor: a model.Doctor object containing the new specialization details.
//
// Returns:
//  - A pointer to the updated model.Doctor object.
//  - An error if the update process fails.

func (r *DoctorRepositoryImpl) UpdateDoctor(doctorID uint, doctor model.Doctor) (*model.Doctor, error) {
	tx := r.DB.Begin()

	var existingDoctor model.Doctor
	if err := r.DB.First(&existingDoctor, doctorID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	existingDoctor.Specialization = doctor.Specialization

	if err := r.DB.Save(&existingDoctor).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &existingDoctor, nil
}

func (r *DoctorRepositoryImpl) DeleteDoctor(doctorID uint) error {
	tx := r.DB.Begin()

	if err := tx.Delete(&model.Doctor{}, doctorID).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
