package services

import (
	"booking-klinik/model"
	"booking-klinik/repository"
	"errors"
	"time"
)

type DoctorScheduleService interface {
	CreateDoctorSchedule(doctorSchedule model.DoctorSchedule) (*model.DoctorSchedule, error)
	GetDoctorSchedulesByDoctorId(doctorId uint) ([]model.DoctorSchedule, error)
	GetAllDoctorSchedules(limit, offset int) ([]model.DoctorSchedule, error)
	GetDoctorScheduleById(scheduleId uint) (*model.DoctorSchedule, error)
	UpdateDoctorSchedule(scheduleID uint, doctorSchedule model.DoctorSchedule) (*model.DoctorSchedule, error)
	DeleteDoctorSchedule(scheduleID uint, userID uint) error
}

type DoctorScheduleServiceImpl struct {
	DoctorScheduleRepository repository.DoctorScheduleRepository
	DoctorRepository         repository.DoctorRepository
	ServiceRepository        repository.ServiceRepository
}

func (ds *DoctorScheduleServiceImpl) CreateDoctorSchedule(doctorSchedule model.DoctorSchedule) (*model.DoctorSchedule, error) {
	if doctorSchedule.DoctorId == 0 || doctorSchedule.ServiceId == 0 || doctorSchedule.Date.IsZero() || doctorSchedule.StartTime.IsZero() || doctorSchedule.EndTime.IsZero() {
		return nil, errors.New("invalid doctor schedule data")
	}

	if doctorSchedule.Date.Before(time.Now()) {
		return nil, errors.New("date must be in the future")
	}

	if _, err := ds.DoctorRepository.GetDoctorById(doctorSchedule.DoctorId); err != nil {
		return nil, errors.New("doctor not found")
	}

	if _, err := ds.ServiceRepository.GetServiceById(doctorSchedule.ServiceId); err != nil {
		return nil, errors.New("service not found")
	}

	existingSchedule, err := ds.DoctorScheduleRepository.GetDoctorSchedulesByDoctorId(doctorSchedule.DoctorId)
	if err != nil {
		return nil, errors.New("error getting doctor schedules")
	}

	for _, existing := range existingSchedule {
		if existing.Date.Equal(doctorSchedule.Date) && existing.ServiceId == doctorSchedule.ServiceId {
			existingStart, existingEnd := existing.StartTime, existing.EndTime
			newStart, newEnd := doctorSchedule.StartTime, doctorSchedule.EndTime

			if (newStart.Before(existingEnd) && newEnd.After(existingStart)) || (existingStart.Before(newEnd) && existingEnd.After(newStart)) || (newStart.Equal(existingStart) && newEnd.Equal(existingEnd)) || (existingStart.Equal(newStart) || existingEnd.Equal(newEnd)) {
				return nil, errors.New("doctor schedule already exists")
			}
		}
	}

	if err := ds.DoctorScheduleRepository.CreateDoctorSchedule(&doctorSchedule); err != nil {
		return nil, err
	}

	return &doctorSchedule, nil
}

func (ds *DoctorScheduleServiceImpl) GetDoctorSchedulesByDoctorId(doctorId uint) ([]model.DoctorSchedule, error) {
	schedules, err := ds.DoctorScheduleRepository.GetDoctorSchedulesByDoctorId(doctorId)
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (ds *DoctorScheduleServiceImpl) UpdateDoctorSchedule(scheduleID uint, doctorSchedule model.DoctorSchedule) (*model.DoctorSchedule, error) {
	existingSchedule, err := ds.DoctorScheduleRepository.GetDoctorSchedulesByDoctorId(doctorSchedule.DoctorId)
	if err != nil {
		return nil, err
	}

	for _, existing := range existingSchedule {
		if existing.Date == doctorSchedule.Date {
			existingStart, existingEnd := existing.StartTime, existing.EndTime
			newStart, newEnd := doctorSchedule.StartTime, doctorSchedule.EndTime

			if (newStart.Before(existingEnd) && newEnd.After(existingStart)) || (existingStart.Before(newEnd) && existingEnd.After(newStart)) {
				return nil, errors.New("doctor schedule already exists")
			}

		}
	}

	if err := ds.DoctorScheduleRepository.UpdateDoctorSchedule(&doctorSchedule); err != nil {
		return nil, err
	}

	return &doctorSchedule, nil
}

func (ds *DoctorScheduleServiceImpl) DeleteDoctorSchedule(scheduleID uint, userID uint) error {
	err := ds.DoctorScheduleRepository.DeleteDoctorSchedule(scheduleID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DoctorScheduleServiceImpl) GetAllDoctorSchedules(limit, offset int) ([]model.DoctorSchedule, error) {
	schedules, err := ds.DoctorScheduleRepository.GetAllDoctorSchedules(limit, offset)
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (ds *DoctorScheduleServiceImpl) GetDoctorScheduleById(scheduleId uint) (*model.DoctorSchedule, error) {
	schedule, err := ds.DoctorScheduleRepository.GetDoctorSchedulesById(scheduleId)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}
