package repository

import (
	"booking-klinik/model"
	"time"

	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(booking *model.Booking) error
	GetAllBookings(limit, offset int) ([]model.Booking, int64, error)
	GetBookingById(id uint) (*model.Booking, error)
	GetBookingsByUserId(userId uint, limit, offset int) ([]model.Booking, error)
	GetBookingsByDoctorId(doctorId uint, limit, offset int) ([]model.Booking, error)
	GetBookingsByDoctorAndDate(doctorId uint, bookingDate time.Time) ([]model.Booking, error)
	UpdateBooking(bookingID uint, booking model.Booking) (*model.Booking, error)
	DeleteBooking(bookingID uint, userID uint) error
}

type BookingRepositoryImpl struct {
	DB                *gorm.DB
	ServiceRepository ServiceRepository
}

func (r *BookingRepositoryImpl) CreateBooking(booking *model.Booking) error {
	tx := r.DB.Begin()
	defer tx.Commit()

	if err := tx.Create(booking).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Preload("User").Preload("Service").First(booking, booking.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *BookingRepositoryImpl) GetAllBookings(limit, offset int) ([]model.Booking, int64, error) {
	var bookings []model.Booking
	var totalRows int64

	if err := r.DB.Model(&model.Booking{}).Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.Limit(limit).Offset(offset).Find(&bookings).Error; err != nil {
		return nil, 0, err
	}
	return bookings, totalRows, nil
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
	defer tx.Commit()
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

	return &existingBooking, nil
}

func (r *BookingRepositoryImpl) DeleteBooking(bookingID uint, userID uint) error {
	tx := r.DB.Begin()
	defer tx.Commit()

	var booking model.Booking
	if err := tx.First(&booking, bookingID).Error; err != nil {
		tx.Rollback()
		return err
	}

	booking.UpdatedBy = userID

	if err := tx.Save(&booking).Error; err != nil {
		tx.Rollback()
		return err
	}
	// Delete the booking
	if err := tx.Delete(&model.Booking{}, bookingID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *BookingRepositoryImpl) GetBookingsByDoctorAndDate(doctorId uint, bookingDate time.Time) ([]model.Booking, error) {
	var bookings []model.Booking
	if err := r.DB.Preload("User").Preload("Service").Where("doctor_id = ? AND booking_date = ? AND status != ?", doctorId, bookingDate, "cancelled").Order("booking_time asc").Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}
