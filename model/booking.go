package model

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	UserId      uint    `json:"user_id" gorm:"not null"`
	DoctorId    uint    `json:"doctor_id" gorm:"not null"`
	ServiceId   uint    `json:"service_id" gorm:"not null"`
	BookingDate string  `json:"booking_date" gorm:"not null"`
	BookingTime string  `json:"booking_time" gorm:"not null"`
	Status      string  `json:"status" gorm:"not null;default:pending"`
	Notes       string  `json:"notes" gorm:"type:text"`
	CreatedBy   uint    `json:"created_by" gorm:"not null"`
	UpdatedBy   uint    `json:"updated_by"`
	User        User    `json:"-" gorm:"foreignKey:UserId;references:ID"`
	Doctor      Doctor  `json:"-" gorm:"foreignKey:DoctorId;references:ID"`
	Service     Service `json:"-" gorm:"foreignKey:ServiceId;references:ID"`
}

type BookingRequest struct {
	DoctorId    uint   `json:"doctor_id"`
	ServiceId   uint   `json:"service_id"`
	BookingDate string `json:"booking_date"`
	BookingTime string `json:"booking_time"`
	Notes       string `json:"notes"`
}

type BookingResponse struct {
	ID          uint   `json:"id"`
	DoctorName  string `json:"doctor_name"`
	ServiceName string `json:"service_name"`
	BookingDate string `json:"booking_date"`
	BookingTime string `json:"booking_time"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
}

type UpdateRequest struct {
	BookingTime string `json:"booking_time"`
	BookingDate string `json:"booking_date"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
}
