package repository

import (
	"booking-klinik/model"

	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(booking *model.Booking) error
	GetAllBookings(limit, offset int) ([]model.Booking, error)
	GetBookingById(id uint) (*model.Booking, error)
	GetBookingsByUserId(userId uint, limit, offset int) ([]model.Booking, error)
	GetBookingsByDoctorId(doctorId uint, limit, offset int) ([]model.Booking, error)
	UpdateBooking(bookingID uint, booking model.Booking) (*model.Booking, error)
	DeleteBooking(bookingID uint) error
}

type BookingRepositoryImpl struct {
	DB *gorm.DB
}

func (r *BookingRepositoryImpl) CreateBooking(booking *model.Booking) error {
	tx := r.DB.Begin()
	if err := tx.Create(booking).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *BookingRepositoryImpl) GetAllBookings(limit, offset int) ([]model.Booking, error) {
	var bookings []model.Booking
	if err := r.DB.Limit(limit).Offset(offset).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepositoryImpl) GetBookingById(id uint) (*model.Booking, error) {
	var booking model.Booking
	if err := r.DB.First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepositoryImpl) GetBookingsByUserId(userId uint, limit, offset int) ([]model.Booking, error) {
	var bookings []model.Booking
	if err := r.DB.Preload("Doctor").Preload("Service").Where("user_id = ?", userId).Limit(limit).Offset(offset).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepositoryImpl) GetBookingsByDoctorId(doctorId uint, limit, offset int) ([]model.Booking, error) {
	var bookings []model.Booking
	if err := r.DB.Preload("User").Preload("Service").Where("doctor_id = ?", doctorId).Limit(limit).Offset(offset).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepositoryImpl) UpdateBooking(bookingID uint, booking model.Booking) (*model.Booking, error) {
	tx := r.DB.Begin()

	var existingBooking model.Booking
	if err := r.DB.First(&existingBooking, bookingID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	existingBooking.Status = booking.Status
	existingBooking.Notes = booking.Notes

	if err := r.DB.Save(&existingBooking).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &existingBooking, nil
}

func (r *BookingRepositoryImpl) DeleteBooking(bookingID uint) error {
	tx := r.DB.Begin()

	// Delete the booking
	if err := tx.Delete(&model.Booking{}, bookingID).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
