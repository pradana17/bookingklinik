package main

import (
	"booking-klinik/config"
	"booking-klinik/routes"

	"github.com/joho/godotenv"
)

func main() {

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
