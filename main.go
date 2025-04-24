package main

import (
	"booking-klinik/config"
	"booking-klinik/routes"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	loc, _ := time.LoadLocation("Asia/Jakarta")
	time.Local = loc
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	//Connect DB
	db := config.ConnectDB()
	//Migrate DB
	config.MigrateDB(db)

	//Setup Router
	r := routes.SetupRouter(db)
	r.Run(":8080")
}
