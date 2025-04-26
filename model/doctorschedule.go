package model

import (
	"time"

	"gorm.io/gorm"
)

type DoctorSchedule struct {
	gorm.Model
	DoctorId  uint      `json:"doctor_id" gorm:"not null"`
	ServiceId uint      `json:"service_id" gorm:"not null"`
	Date      time.Time `json:"date" time_format:"YYYY-MM-DD" gorm:"not null"`
	StartTime time.Time `json:"start_time" time_format:"15:04" gorm:"not null"`
	EndTime   time.Time `json:"end_time" time_format:"15:04" gorm:"not null"`
	CreatedBy uint      `json:"created_by" gorm:"not null"`
	UpdatedBy uint      `json:"updated_by"`
	Doctor    Doctor    `json:"doctor" gorm:"foreignKey:DoctorId;references:ID"`
	Service   Service   `json:"service" gorm:"foreignKey:ServiceId;references:ID"`
}

type DoctorScheduleResponse struct {
	ID        uint      `json:"id"`
	DoctorID  uint      `json:"doctor_id"`
	Date      time.Time `json:"date" time_format:"YYYY-MM-DD"`
	StartTime time.Time `json:"start_time" time_format:"15:04"`
	EndTime   time.Time `json:"end_time" time_format:"15:04"`
}

type DoctorScheduleRequest struct {
	DoctorID  uint   `json:"doctor_id"`
	ServiceID uint   `json:"service_id"`
	Date      string `json:"date" time_format:"2006-01-02"`
	StartTime string `json:"start_time" time_format:"15:04"`
	EndTime   string `json:"end_time" time_format:"15:04"`
}
