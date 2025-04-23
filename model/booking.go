package model

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserId      uint      `json:"user_id" gorm:"not null"`
	DoctorId    uint      `json:"doctor_id" gorm:"not null"`
	ServiceId   uint      `json:"service_id" gorm:"not null"`
	BookingDate time.Time `json:"booking_date" time_format:"2006-01-02" gorm:"not null"`
	BookingTime time.Time `json:"booking_time" time_format:"15:04" gorm:"not null"`
	Status      string    `json:"status" gorm:"not null;default:pending"`
	Notes       string    `json:"notes" gorm:"type:text"`
	CreatedBy   uint      `json:"created_by" gorm:"not null"`
	UpdatedBy   uint      `json:"updated_by"`
	User        User      `json:"-" gorm:"foreignKey:UserId;references:ID"`
	Doctor      Doctor    `json:"-" gorm:"foreignKey:DoctorId;references:ID"`
	Service     Service   `json:"-" gorm:"foreignKey:ServiceId;references:ID"`
}

type BookingRequest struct {
	DoctorId    uint      `json:"doctor_id"`
	ServiceId   uint      `json:"service_id"`
	BookingDate time.Time `json:"booking_date" time_format:"2006-01-02"`
	BookingTime time.Time `json:"booking_time" time_format:"15:04"`
	Notes       string    `json:"notes"`
}

type BookingResponse struct {
	ID          uint      `json:"id"`
	DoctorName  string    `json:"doctor_name"`
	ServiceName string    `json:"service_name"`
	BookingDate time.Time `json:"booking_date" time_format:"2006-01-02"`
	BookingTime time.Time `json:"booking_time" time_format:"15:04"`
	Status      string    `json:"status"`
	Notes       string    `json:"notes"`
}

type UpdateRequest struct {
	BookingDate time.Time `json:"booking_date" time_format:"2006-01-02"`
	BookingTime time.Time `json:"booking_time" time_format:"15:04"`
	Status      string    `json:"status"`
	Notes       string    `json:"notes"`
}
