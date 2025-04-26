package model

import (
	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	UserId         uint             `json:"user_id" gorm:"not null;uniqueIndex:idx_user_id"`
	Specialization string           `json:"specialization" gorm:"not null"`
	CreatedBy      uint             `json:"created_by" gorm:"not null"`
	UpdatedBy      uint             `json:"updated_by"`
	User           User             `json:"-" gorm:"foreignKey:UserId;references:ID"`
	Bookings       []Booking        `json:"-" gorm:"foreignKey:DoctorId;references:ID"`
	Schedules      []DoctorSchedule `json:"schedules" gorm:"foreignKey:DoctorId;references:ID"`
}

type DoctorResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
}

type DoctorRequest struct {
	ID uint `json:"id"`
}

type GetDoctorNameRequest struct {
	UserID uint `json:"user_id"`
}
