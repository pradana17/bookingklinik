package model

import (
	"time"

	"gorm.io/gorm"
)

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
	DoctorId  uint      `json:"doctor_id" gorm:"not null"`
	Date      time.Time `json:"date" time_format:"2006-01-02" gorm:"not null;index;uniqueIndex:idx_doctorschedule_date"`
	StartTime time.Time `json:"start_time" time_format:"15:04" gorm:"not null"`
	EndTime   time.Time `json:"end_time" time_format:"15:04" gorm:"not null"`
	CreatedBy uint      `json:"created_by" gorm:"not null"`
	UpdatedBy uint      `json:"updated_by"`
	Doctor    Doctor    `json:"-" gorm:"foreignKey:DoctorId;references:ID"`
}

type DoctorResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
}

type DoctorScheduleResponse struct {
	ID        uint      `json:"id"`
	DoctorID  uint      `json:"doctor_id"`
	Date      time.Time `json:"date" time_format:"2006-01-02"`
	StartTime time.Time `json:"start_time" time_format:"15:04"`
	EndTime   time.Time `json:"end_time" time_format:"15:04"`
}

type DoctorScheduleRequest struct {
	Date      time.Time `json:"date" time_format:"2006-01-02"`
	StartTime time.Time `json:"start_time" time_format:"15:04"`
	EndTime   time.Time `json:"end_time" time_format:"15:04"`
}
