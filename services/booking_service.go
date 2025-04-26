package services

import (
	"booking-klinik/model"
	"booking-klinik/repository"
	"booking-klinik/utils"
	"errors"
	"fmt"
	"time"
)

type BookingService interface {
	CreateBooking(booking model.Booking) (*model.Booking, error)
	GetAllBookings(limit, offset int, userRole string, userId uint) ([]model.Booking, *utils.Paginator, error)
	GetBookingById(id uint, userID uint, userRole string) (*model.Booking, error)
	GetBookingsByUserId(userId uint, limit, offset int) ([]model.Booking, *utils.Paginator, error)
	GetBookingsByDoctorId(doctorId uint, limit, offset int) ([]model.Booking, *utils.Paginator, error)
	GetDoctorName(doctorId uint) (string, error)
	UpdateBooking(bookingID uint, booking model.Booking, userRole string) (*model.Booking, error)
	DeleteBooking(bookingID uint, userRole string, userID uint) error
}

type BookingServicesImpl struct {
	BookingRepository        repository.BookingRepository
	DoctorRepository         repository.DoctorRepository
	ServiceRepository        repository.ServiceRepository
	DoctorScheduleRepository repository.DoctorScheduleRepository
	UserRepository           repository.UserRepository
}

func (s *BookingServicesImpl) CreateBooking(booking model.Booking) (*model.Booking, error) {
	// Validate the booking
	doctor, err := s.DoctorRepository.GetDoctorById(booking.DoctorId)
	if err != nil {
		return nil, errors.New("doctor not found")
	}
	service, err := s.ServiceRepository.GetServiceById(booking.ServiceId)
	if err != nil || !service.IsActive {
		return nil, errors.New("service is inactive or not found")
	}

	schedules, err := s.DoctorScheduleRepository.GetDoctorSchedulesByDoctorId(doctor.ID)
	if err != nil {
		return nil, errors.New("doctor schedule not found")
	}

	if len(schedules) == 0 {
		return nil, errors.New("doctor schedule not found")
	}

	available := false

	for _, schedule := range schedules {
		fmt.Println("Schedule date:", schedule.Date.Format("2006-01-02"))
		fmt.Println("Booking date :", booking.BookingDate.Format("2006-01-02"))
		fmt.Println("Schedule start:", schedule.StartTime)
		fmt.Println("Booking time  :", booking.BookingTime)
		if schedule.Date.Format("2006-01-02") == booking.BookingDate.Format("2006-01-02") {
			if (booking.BookingTime.After(schedule.StartTime) || booking.BookingTime.Equal(schedule.StartTime)) && (booking.BookingTime.Before(schedule.EndTime) || booking.BookingTime.Equal(schedule.EndTime)) {
				available = true
				break
			}
		}
	}

	if !available {
		return nil, errors.New("doctor is not available at this date and time")
	}

	conflict, nextAvailableTime, err := s.CheckBookingConflict(booking.DoctorId, booking.BookingDate, booking.BookingTime, service.DurationMinutes)
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

func (s *BookingServicesImpl) GetAllBookings(limit, offset int, userRole string, userId uint) ([]model.Booking, *utils.Paginator, error) {
	var bookings []model.Booking
	var totalRows int64
	var err error

	switch userRole {
	case "patient":
		bookings, totalRows, err = s.BookingRepository.GetBookingsByUserId(userId, limit, offset)
	case "doctor":
		doctorID, err := s.DoctorRepository.GetDoctorIDbyUserID(userId)
		if err != nil {
			return nil, nil, err
		}
		bookings, totalRows, err = s.BookingRepository.GetBookingsByDoctorId(doctorID, limit, offset)
		if err != nil {
			return nil, nil, err
		}
	case "admin":
		bookings, totalRows, err = s.BookingRepository.GetAllBookings(limit, offset)
	default:
		return nil, nil, errors.New("invalid user role")
	}

	if err != nil {
		return nil, nil, err
	}

	pagination := &utils.Paginator{Limit: limit, Offset: offset, Page: (offset / limit) + 1, TotalRows: totalRows}

	pagination.TotalPages = (totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit)
	return bookings, pagination, nil
}

func (s *BookingServicesImpl) GetBookingById(id uint, userID uint, userRole string) (*model.Booking, error) {
	if userRole == "patient" {
		booking, err := s.BookingRepository.GetBookingById(id)
		if err != nil {
			return nil, err
		}
		if booking.UserId != userID {
			return nil, errors.New("you can only access your own bookings")
		}
		return booking, nil
	} else if userRole == "doctor" {
		doctorID, err := s.DoctorRepository.GetDoctorIDbyUserID(userID)
		if err != nil {
			return nil, err
		}
		booking, err := s.BookingRepository.GetBookingById(id)
		if err != nil {
			return nil, err
		}
		if booking.DoctorId != doctorID {
			return nil, errors.New("you can only access your patients bookings")
		}
		return booking, nil
	} else if userRole == "admin" {
		return s.BookingRepository.GetBookingById(id)
	}
	return nil, errors.New("invalid user role")
}

func (s *BookingServicesImpl) GetBookingsByUserId(userId uint, limit, offset int) ([]model.Booking, *utils.Paginator, error) {
	bookings, totalRows, err := s.BookingRepository.GetBookingsByUserId(userId, limit, offset)
	if err != nil {
		return nil, nil, err
	}
	pagination := &utils.Paginator{Limit: limit, Offset: offset, Page: (offset / limit) + 1, TotalRows: totalRows}

	pagination.TotalPages = (totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit)
	return bookings, pagination, nil
}

func (s *BookingServicesImpl) GetBookingsByDoctorId(doctorId uint, limit, offset int) ([]model.Booking, *utils.Paginator, error) {
	booking, totalRows, err := s.BookingRepository.GetBookingsByDoctorId(doctorId, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	pagination := &utils.Paginator{Limit: limit, Offset: offset, Page: (offset / limit) + 1, TotalRows: totalRows}

	pagination.TotalPages = (totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit)

	return booking, pagination, nil
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

	err = s.BookingRepository.DeleteBooking(bookingID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *BookingServicesImpl) CheckBookingConflict(doctorId uint, bookingDate time.Time, bookingTime time.Time, durationMinutes int) (bool, time.Time, error) {
	bookings, err := s.BookingRepository.GetBookingsByDoctorAndDate(doctorId, bookingDate)
	if err != nil {
		return false, time.Time{}, err
	}

	newEnd := bookingTime.Add(time.Duration(durationMinutes) * time.Minute)
	for _, booking := range bookings {
		service, err := s.ServiceRepository.GetServiceById(booking.ServiceId)
		if err != nil {
			return false, time.Time{}, err
		}

		existingEnd := booking.BookingTime.Add(time.Duration(service.DurationMinutes) * time.Minute)

		if bookingTime.Before(existingEnd) && newEnd.After(booking.BookingTime) {
			return true, existingEnd, nil
		}
	}

	return false, time.Time{}, nil
}

func (s *BookingServicesImpl) GetDoctorName(doctorId uint) (string, error) {
	doctor, err := s.DoctorRepository.GetDoctorById(doctorId)
	if err != nil {
		return "", err
	}

	user, err := s.UserRepository.GetUserById(doctor.UserId)
	if err != nil {
		return "", err
	}
	return user.Name, nil
}
