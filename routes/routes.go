package routes

import (
	"booking-klinik/controllers"
	"booking-klinik/middleware"
	"booking-klinik/repository"
	"booking-klinik/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	userRepository := &repository.UserRepositoryImpl{DB: db}
	serviceRepository := &repository.ServiceRepositoryImpl{DB: db}
	doctorScheduleRepository := &repository.DoctorScheduleRepositoryImpl{DB: db}
	bookingRepository := &repository.BookingRepositoryImpl{DB: db}
	doctorRepository := &repository.DoctorRepositoryImpl{DB: db}

	userService := &services.UserServicesImpl{UserRepository: userRepository}
	doctorService := &services.DoctorServicesImpl{
		DoctorRepository: doctorRepository,
		UserRepository:   userRepository,
	}
	bookingService := &services.BookingServicesImpl{
		BookingRepository:        bookingRepository,
		DoctorRepository:         doctorRepository,
		ServiceRepository:        serviceRepository,
		DoctorScheduleRepository: doctorScheduleRepository,
		UserRepository:           userRepository}
	doctorScheduleService := &services.DoctorScheduleServiceImpl{DoctorScheduleRepository: doctorScheduleRepository, DoctorRepository: doctorRepository, ServiceRepository: serviceRepository}
	serviceService := &services.ServiceServiceImpl{ServiceRepository: serviceRepository}

	//User Routes
	userController := &controllers.UserController{UserService: userService}
	r.POST("/register", userController.RegisterUser)
	r.POST("/login", userController.LoginUser)

	userGroup := r.Group("/user")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.PUT("/password", userController.UpdatePassword)
	}

	//Booking Routes
	bookingController := &controllers.BookingController{BookingService: bookingService, DoctorService: doctorService, UserService: userService}
	bookingGroup := r.Group("/booking")
	bookingGroup.Use(middleware.AuthMiddleware())
	{
		bookingGroup.POST("/", bookingController.CreateBooking)
		bookingGroup.GET("/", bookingController.GetAllBookings)
		bookingGroup.GET("/:id", bookingController.GetBookingsById)
		bookingGroup.GET("/user/:user_id", bookingController.GetBookingsByUserId)
		bookingGroup.GET("/doctor/:doctor_id", bookingController.GetBookingsByDoctorId)
		bookingGroup.PUT("/:id", bookingController.UpdateBooking)
		bookingGroup.DELETE("/:id", bookingController.DeleteBooking)
	}

	//Doctor Routes

	doctorController := &controllers.DoctorController{DoctorService: doctorService}
	doctorGroup := r.Group("/doctor")
	doctorGroup.Use(middleware.AuthMiddleware(), middleware.RoleCheckMiddleware("admin", "doctor"))
	{
		doctorGroup.POST("/", doctorController.CreateDoctor)
		doctorGroup.GET("/", doctorController.GetAllDoctors)
		doctorGroup.GET("/:id", doctorController.GetDoctorById)
		doctorGroup.PUT("/:id", doctorController.UpdateDoctor)
		doctorGroup.DELETE("/:id", doctorController.DeleteDoctor)
	}

	//Doctor Schedule Routes

	doctorScheduleController := &controllers.DoctorScheduleController{DoctorScheduleService: doctorScheduleService}
	doctorScheduleGroup := r.Group("/doctorschedule")
	doctorScheduleGroup.Use(middleware.AuthMiddleware(), middleware.RoleCheckMiddleware("admin", "doctor"))
	{
		doctorScheduleGroup.POST("/", doctorScheduleController.CreateDoctorSchedule)
		doctorScheduleGroup.GET("/", doctorScheduleController.GetAllDoctorSchedules)
		doctorScheduleGroup.GET("/:id", doctorScheduleController.GetDoctorScheduleById)
		doctorScheduleGroup.PUT("/:id", doctorScheduleController.UpdateDoctorSchedule)
		doctorScheduleGroup.DELETE("/:id", doctorScheduleController.DeleteDoctorSchedule)
	}

	//Service Routes

	serviceController := &controllers.ServiceController{ServiceService: serviceService}

	serviceGroup := r.Group("/service")
	serviceGroup.Use(middleware.AuthMiddleware())
	{
		serviceGroup.GET("/", serviceController.GetAllServices)
		serviceGroup.GET("/:id", serviceController.GetServiceById)
		serviceGroup.Use(middleware.RoleCheckMiddleware("admin"))
		{
			serviceGroup.POST("/", serviceController.CreateService)
			serviceGroup.PUT("/:id", serviceController.UpdateService)
			serviceGroup.DELETE("/:id", serviceController.DeleteService)
		}
	}

	return r
}
