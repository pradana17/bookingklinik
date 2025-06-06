package model

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	Name            string           `json:"name" gorm:"unique;not null"`
	Description     string           `json:"description" gorm:"type:text"`
	Price           int              `json:"price" gorm:"not null"`
	DurationMinutes int              `json:"duration_minutes" gorm:"not null"`
	IsActive        bool             `json:"is_active" gorm:"not null"`
	CreatedBy       uint             `json:"created_by" gorm:"not null"`
	UpdatedBy       uint             `json:"updated_by"`
	Bookings        []Booking        `json:"-" gorm:"foreignKey:ServiceId;references:ID"`
	Schedules       []DoctorSchedule `json:"doctor_schedule" gorm:"foreignKey:ServiceId;references:ID"`
}

type ServiceResponse struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Price           int    `json:"price"`
	DurationMinutes int    `json:"duration_minutes"`
	IsActive        bool   `json:"is_active"`
}
