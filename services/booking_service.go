package services

import (
	"booking-klinik/model"
	"booking-klinik/repository"
	"errors"
)

type BookingService interface {
	CreateBooking(booking model.Booking) (*model.Booking, error)
	GetAllBookings(limit, offset int) ([]model.Booking, error)
	GetBookingById(id uint) (*model.Booking, error)
	GetBookingsByUserId(userId uint, limit, offset int) ([]model.Booking, error)
	GetBookingsByDoctorId(doctorId uint, limit, offset int) ([]model.Booking, error)
	UpdateBooking(bookingID uint, booking model.Booking, userRole string) (*model.Booking, error)
	DeleteBooking(bookingID uint, userRole string, userID uint) error
}

type BookingServicesImpl struct {
	BookingRepository repository.BookingRepository
}

func (s *BookingServicesImpl) CreateBooking(booking model.Booking) (*model.Booking, error) {
	err := s.BookingRepository.CreateBooking(&booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *BookingServicesImpl) GetAllBookings(limit, offset int) ([]model.Booking, error) {
	return s.BookingRepository.GetAllBookings(limit, offset)
}

func (s *BookingServicesImpl) GetBookingById(id uint) (*model.Booking, error) {
	return s.BookingRepository.GetBookingById(id)
}

func (s *BookingServicesImpl) GetBookingsByUserId(userId uint, limit, offset int) ([]model.Booking, error) {
	return s.BookingRepository.GetBookingsByUserId(userId, limit, offset)
}

func (s *BookingServicesImpl) GetBookingsByDoctorId(doctorId uint, limit, offset int) ([]model.Booking, error) {
	return s.BookingRepository.GetBookingsByDoctorId(doctorId, limit, offset)
}

func (s *BookingServicesImpl) UpdateBooking(bookingID uint, booking model.Booking, userRole string) (*model.Booking, error) {
	existingBooking, err := s.BookingRepository.GetBookingById(bookingID)
	if err != nil {
		return nil, err
	}

	if userRole == "patient" && existingBooking.UserId != booking.UserId {
		return nil, errors.New("you can only update your own bookings")
	}

	if existingBooking.Status == "confirmed" && userRole == "patient" {
		return nil, errors.New("booking already confirmed")
	}

	if userRole == "patient" {
		existingBooking.Notes = booking.Notes
		if booking.BookingTime != "" {
			existingBooking.BookingTime = booking.BookingTime
		}
		if booking.BookingDate != "" {
			existingBooking.BookingDate = booking.BookingDate
		}
	} else if userRole == "doctor" || userRole == "admin" {
		existingBooking.Status = booking.Status
		existingBooking.Notes = booking.Notes
	}

	existingBooking, err = s.BookingRepository.UpdateBooking(bookingID, *existingBooking)
	if err != nil {
		return nil, err
	}

	return existingBooking, nil
}

func (s *BookingServicesImpl) DeleteBooking(bookingID uint, userRole string, userID uint) error {
	booking, err := s.BookingRepository.GetBookingById(bookingID)
	if err != nil {
		return err
	}
	if booking.Status == "confirmed" {
		if userRole == "patient" && booking.UserId != userID {
			return errors.New("you can only delete your own bookings")
		}
		return errors.New("booking already confirmed")
	}

	if userRole == "doctor" {
		return errors.New("doctors cannot delete bookings")
	}

	err = s.BookingRepository.DeleteBooking(bookingID)
	if err != nil {
		return err
	}
	return nil
}
