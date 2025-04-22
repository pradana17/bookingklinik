package services

import (
	"booking-klinik/model"
	"booking-klinik/repository"
	"errors"
)

type DoctorScheduleService interface {
	CreateDoctorSchedule(doctorSchedule model.DoctorSchedule) (*model.DoctorSchedule, error)
	GetDoctorSchedulesByDoctorId(doctorId uint) ([]model.DoctorSchedule, error)
	GetAllDoctorSchedules(limit, offset int) ([]model.DoctorSchedule, error)
	GetDoctorScheduleById(scheduleId uint) (*model.DoctorSchedule, error)
	UpdateDoctorSchedule(scheduleID uint, doctorSchedule model.DoctorSchedule) (*model.DoctorSchedule, error)
	DeleteDoctorSchedule(scheduleID uint) error
}

type DoctorScheduleServiceImpl struct {
	DoctorScheduleRepository repository.DoctorScheduleRepository
}

func (ds *DoctorScheduleServiceImpl) CreateDoctorSchedule(doctorSchedule model.DoctorSchedule) (*model.DoctorSchedule, error) {
	if doctorSchedule.DoctorId == 0 {
		return nil, errors.New("doctor id is required")
	}

	existingSchedule, err := ds.DoctorScheduleRepository.GetDoctorSchedulesByDoctorId(doctorSchedule.DoctorId)
	if err != nil {
		return nil, err
	}

	for _, existing := range existingSchedule {
		if existing.Day == doctorSchedule.Day {
			existingStart, existingEnd := existing.StartTime, existing.EndTime
			newStart, newEnd := doctorSchedule.StartTime, doctorSchedule.EndTime

			if (newStart < existingEnd && newEnd > existingStart) || (existingStart < newEnd && existingEnd > newStart) {
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
		if existing.Day == doctorSchedule.Day {
			existingStart, existingEnd := existing.StartTime, existing.EndTime
			newStart, newEnd := doctorSchedule.StartTime, doctorSchedule.EndTime

			if (newStart < existingEnd && newEnd > existingStart) || (existingStart < newEnd && existingEnd > newStart) {
				return nil, errors.New("doctor schedule already exists")
			}

		}
	}

	if err := ds.DoctorScheduleRepository.UpdateDoctorSchedule(&doctorSchedule); err != nil {
		return nil, err
	}

	return &doctorSchedule, nil
}

func (ds *DoctorScheduleServiceImpl) DeleteDoctorSchedule(scheduleID uint) error {
	err := ds.DoctorScheduleRepository.DeleteDoctorSchedule(scheduleID)
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
