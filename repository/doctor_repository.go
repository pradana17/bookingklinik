package repository

import (
	"booking-klinik/model"

	"gorm.io/gorm"
)

type DoctorRepository interface {
	CreateDoctor(doctor *model.Doctor) error
	GetAllDoctors(limit, offset int) ([]model.Doctor, int64, error)
	GetDoctorById(id uint) (*model.Doctor, error)
	GetDoctorIDbyUserID(userID uint) (uint, error)
	UpdateDoctor(doctorID uint, doctor model.Doctor) (*model.Doctor, error)
	DeleteDoctor(doctorID uint, userID uint) error
}

type DoctorRepositoryImpl struct {
	DB *gorm.DB
}

func (r *DoctorRepositoryImpl) CreateDoctor(doctor *model.Doctor) error {
	tx := r.DB.Begin()
	defer tx.Commit()
	if err := tx.Create(doctor).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// GetAllDoctors gets all doctors with the given limit and offset.
//
// This function will return a list of doctors with the given limit and offset.
// If the query is not successful, it will return an error.
func (r *DoctorRepositoryImpl) GetAllDoctors(limit, offset int) ([]model.Doctor, int64, error) {
	var totalRows int64
	var doctors []model.Doctor
	if err := r.DB.Model(&model.Doctor{}).Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.Limit(limit).Offset(offset).Preload("User").Find(&doctors).Error; err != nil {
		return nil, 0, err
	}
	return doctors, totalRows, nil
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
	defer tx.Commit()

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
	return &existingDoctor, nil
}

func (r *DoctorRepositoryImpl) DeleteDoctor(doctorID uint, userID uint) error {
	tx := r.DB.Begin()
	defer tx.Commit()

	var doctor model.Doctor
	if err := tx.First(&doctor, doctorID).Error; err != nil {
		tx.Rollback()
		return err
	}

	doctor.UpdatedBy = userID

	if err := tx.Save(&doctor).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&model.Doctor{}, doctorID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *DoctorRepositoryImpl) GetDoctorIDbyUserID(userID uint) (uint, error) {
	var doctor model.Doctor
	if err := r.DB.Preload("User").First(&doctor, "user_id = ?", userID).Error; err != nil {
		return 0, err
	}
	return doctor.ID, nil
}
