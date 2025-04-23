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

	//User Routes
	userRepository := &repository.UserRepositoryImpl{DB: db}
	userService := &services.UserServicesImpl{UserRepository: userRepository}
	userController := &controllers.UserController{UserService: userService}
	r.POST("/register", userController.RegisterUser)
	r.POST("/login", userController.LoginUser)

	userGroup := r.Group("/user")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.PUT("/password", userController.UpdatePassword)
	}

	//Booking Routes
	bookingRepository := &repository.BookingRepositoryImpl{DB: db}
	bookingService := &services.BookingServicesImpl{BookingRepository: bookingRepository}
	bookingController := &controllers.BookingController{BookingService: bookingService}

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
	doctorRepository := &repository.DoctorRepositoryImpl{DB: db}
	doctorService := &services.DoctorServicesImpl{DoctorRepository: doctorRepository}
	doctorController := &controllers.DoctorController{DoctorService: doctorService}
	doctorGroup := r.Group("/doctor")
	doctorGroup.Use(middleware.RoleCheckMiddleware("admin", "doctor"), middleware.AuthMiddleware())
	{
		doctorGroup.POST("/", doctorController.CreateDoctor)
		doctorGroup.GET("/", doctorController.GetAllDoctors)
		doctorGroup.GET("/:id", doctorController.GetDoctorById)
		doctorGroup.PUT("/:id", doctorController.UpdateDoctor)
		doctorGroup.DELETE("/:id", doctorController.DeleteDoctor)
	}

	//Doctor Schedule Routes
	doctorScheduleRepository := &repository.DoctorScheduleRepositoryImpl{DB: db}
	doctorScheduleService := &services.DoctorScheduleServiceImpl{DoctorScheduleRepository: doctorScheduleRepository}
	doctorScheduleController := &controllers.DoctorScheduleController{DoctorScheduleService: doctorScheduleService}
	doctorScheduleGroup := r.Group("/doctorschedule")
	doctorScheduleGroup.Use(middleware.RoleCheckMiddleware("admin", "doctor"), middleware.AuthMiddleware())
	{
		doctorScheduleGroup.POST("/", doctorScheduleController.CreateDoctorSchedule)
		doctorScheduleGroup.GET("/", doctorScheduleController.GetAllDoctorSchedules)
		doctorScheduleGroup.GET("/:id", doctorScheduleController.GetDoctorScheduleById)
		doctorScheduleGroup.PUT("/:id", doctorScheduleController.UpdateDoctorSchedule)
		doctorScheduleGroup.DELETE("/:id", doctorScheduleController.DeleteDoctorSchedule)
	}

	//Service Routes
	serviceRepository := &repository.ServiceRepositoryImpl{DB: db}
	serviceService := &services.ServiceServiceImpl{ServiceRepository: serviceRepository}
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
