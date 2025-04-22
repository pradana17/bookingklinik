package config

import (
	"booking-klinik/model"
	"log"

	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&model.User{}, &model.Doctor{}, &model.Booking{}, &model.Service{})
	if err != nil {
		panic(err)
	}

	log.Println("Database migrated successfully")
}
