package services

import (
	"booking-klinik/model"
	"booking-klinik/repository"
	"errors"
	"fmt"
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
	BookingRepository        repository.BookingRepository
	DoctorRepository         repository.DoctorRepository
	ServiceRepository        repository.ServiceRepository
	DoctorScheduleRepository repository.DoctorScheduleRepository
}

func (s *BookingServicesImpl) CreateBooking(booking model.Booking) (*model.Booking, error) {
	// Validate the booking
	_, err := s.DoctorRepository.GetDoctorById(booking.DoctorId)
	if err != nil {
		return nil, errors.New("doctor not found")
	}
	service, err := s.ServiceRepository.GetServiceById(booking.ServiceId)
	if err != nil || !service.IsActive {
		return nil, errors.New("service is inactive or not found")
	}

	schedules, err := s.DoctorScheduleRepository.GetDoctorSchedulesByDoctorId(booking.DoctorId)
	if err != nil {
		return nil, errors.New("doctor schedule not found")
	}

	for _, schedule := range schedules {
		if schedule.Date == booking.BookingDate {
			return nil, errors.New("doctor is not available at this date")
		}
		if booking.BookingTime.Before(schedule.StartTime) && booking.BookingTime.After(schedule.EndTime) {
			return nil, errors.New("booking time is outside doctor's available hours")
		}
	}

	conflict, nextAvailableTime, err := s.BookingRepository.CheckBookingConflict(booking.DoctorId, booking.BookingDate, booking.BookingTime, service.DurationMinutes)
	if err != nil {
		return nil, errors.New("error checking booking conflict")
	}

	if conflict {
		return nil, fmt.Errorf("doctor is already booked at this time. Next available slot starts from %s", nextAvailableTime)
	}

	// Create the booking
	if err := s.BookingRepository.CreateBooking(&booking); err != nil {
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
		if !booking.BookingTime.IsZero() {
			existingBooking.BookingTime = booking.BookingTime
		}
		if !booking.BookingDate.IsZero() {
			existingBooking.BookingDate = booking.BookingDate
		}
	} else if userRole == "doctor" || userRole == "admin" {
		existingBooking.Status = booking.Status
		existingBooking.Notes = booking.Notes
	}

	existingBooking.UpdatedBy = booking.UserId
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
