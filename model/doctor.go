package model

import "gorm.io/gorm"

type Doctor struct {
	gorm.Model
	UserId         uint             `json:"user_id" gorm:"not null"`
	Specialization string           `json:"specialization" gorm:"not null"`
	CreatedBy      uint             `json:"created_by" gorm:"not null"`
	UpdatedBy      uint             `json:"updated_by"`
	User           User             `json:"-" gorm:"foreignKey:UserId;references:ID"`
	Bookings       []Booking        `json:"-" gorm:"foreignKey:DoctorId;references:ID"`
	Schedules      []DoctorSchedule `json:"schedules" gorm:"foreignKey:DoctorId;references:ID"`
}

type DoctorSchedule struct {
	gorm.Model
	DoctorId  uint   `json:"doctor_id" gorm:"not null"`
	Day       string `json:"day" gorm:"not null"`
	StartTime string `json:"start_time" gorm:"not null"`
	EndTime   string `json:"end_time" gorm:"not null"`
	CreatedBy uint   `json:"created_by" gorm:"not null"`
	UpdatedBy uint   `json:"updated_by"`
	Doctor    Doctor `json:"-" gorm:"foreignKey:DoctorId;references:ID"`
}

type DoctorResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
}

type DoctorScheduleResponse struct {
	ID        uint   `json:"id"`
	DoctorID  uint   `json:"doctor_id"`
	Day       string `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type DoctorScheduleRequest struct {
	Day       string `json:"day"`
	StartTime string `json:"start_time" `
	EndTime   string `json:"end_time" `
}
