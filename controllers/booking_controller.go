package controllers

import (
	"booking-klinik/model"
	"booking-klinik/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	BookingService services.BookingService
	DoctorService  services.DoctorServices
	UserService    services.UserService
	ServiceService services.ServiceService
	DoctorSchedule services.DoctorScheduleService
}

func (bc *BookingController) CreateBooking(c *gin.Context) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	time.Local = loc
	var bookingRequest model.BookingRequest
	if err := c.ShouldBindJSON(&bookingRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookingDate, err := time.ParseInLocation("2006-01-02", bookingRequest.BookingDate, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	bookingTime, err := time.ParseInLocation("2006-01-02 15:04", bookingRequest.BookingDate+" "+bookingRequest.BookingTime, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
		return
	}
	userID := c.MustGet("userID").(uint)

	newBooking := model.Booking{
		UserId:      userID,
		DoctorId:    bookingRequest.DoctorId,
		ServiceId:   bookingRequest.ServiceId,
		BookingDate: bookingDate,
		BookingTime: bookingTime,
		Status:      "pending",
		Notes:       bookingRequest.Notes,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	doctor, err := bc.DoctorService.GetDoctorById(bookingRequest.DoctorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := bc.UserService.GetUserById(doctor.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createdBooking, err := bc.BookingService.CreateBooking(newBooking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookingResponse := model.BookingResponse{
		ID:          createdBooking.ID,
		PatientName: createdBooking.User.Name,
		DoctorName:  user.Name,
		ServiceName: createdBooking.Service.Name,
		BookingDate: createdBooking.BookingDate,
		BookingTime: createdBooking.BookingTime,
		Status:      createdBooking.Status,
		Notes:       createdBooking.Notes,
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking created successfully", "booking": bookingResponse})
}

func (bc *BookingController) GetAllBookings(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	userRole := c.MustGet("role").(string)
	userID := c.MustGet("userID").(uint)

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	bookings, pagination, err := bc.BookingService.GetAllBookings(limit, offset, userRole, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var bookingResponses []model.BookingResponse
	for _, booking := range bookings {

		doctor, err := bc.DoctorService.GetDoctorById(booking.DoctorId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctor details"})
			return
		}

		// Ambil nama user (dokter) berdasarkan UserId dokter
		doc, err := bc.UserService.GetUserById(doctor.UserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctor user details"})
			return
		}
		bookingResponses = append(bookingResponses, model.BookingResponse{
			ID:          booking.ID,
			PatientName: booking.User.Name,
			DoctorName:  doc.Name,
			ServiceName: booking.Service.Name,
			BookingDate: booking.BookingDate,
			BookingTime: booking.BookingTime,
			Status:      booking.Status,
			Notes:       booking.Notes,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":         bookingResponses,
		"total_rows":   pagination.TotalRows,
		"total_pages":  pagination.TotalPages,
		"current_page": pagination.Page,
		"limit":        pagination.Limit})
}

func (bc *BookingController) GetBookingsById(c *gin.Context) {
	bookingId := c.Param("id")
	userID := c.MustGet("userID").(uint)
	userRole := c.MustGet("role").(string)
	bookingIdUint, err := strconv.ParseUint(bookingId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	booking, err := bc.BookingService.GetBookingById(uint(bookingIdUint), userID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	doctor, err := bc.DoctorService.GetDoctorById(booking.DoctorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := bc.UserService.GetUserById(doctor.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookingResponse := model.BookingResponse{
		ID:          booking.ID,
		PatientName: booking.User.Name,
		DoctorName:  user.Name,
		ServiceName: booking.Service.Name,
		BookingDate: booking.BookingDate,
		BookingTime: booking.BookingTime,
		Status:      booking.Status,
		Notes:       booking.Notes,
	}

	c.JSON(http.StatusOK, gin.H{"booking": bookingResponse})
}

func (bc *BookingController) GetBookingsByUserId(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")
	userID := c.Param("user_id")

	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	bookings, pagination, err := bc.BookingService.GetBookingsByUserId(uint(userIDUint), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var bookingResponses []model.BookingResponse
	for _, booking := range bookings {
		doctor, err := bc.DoctorService.GetDoctorById(booking.DoctorId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user, err := bc.UserService.GetUserById(doctor.UserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		bookingResponses = append(bookingResponses, model.BookingResponse{
			ID:          booking.ID,
			PatientName: booking.User.Name,
			DoctorName:  user.Name,
			ServiceName: booking.Service.Name,
			BookingDate: booking.BookingDate,
			BookingTime: booking.BookingTime,
			Status:      booking.Status,
			Notes:       booking.Notes,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": bookingResponses, "total_rows": pagination.TotalRows, "total_pages": pagination.TotalPages, "current_page": pagination.Page, "limit": pagination.Limit})
}

func (bc *BookingController) GetBookingsByDoctorId(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")
	doctorId := c.Param("doctor_id")
	doctorIdUint, err := strconv.ParseUint(doctorId, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	bookings, pagination, err := bc.BookingService.GetBookingsByDoctorId(uint(doctorIdUint), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var bookingResponses []model.BookingResponse
	for _, booking := range bookings {
		doctor, err := bc.DoctorService.GetDoctorById(booking.DoctorId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user, err := bc.UserService.GetUserById(doctor.UserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		bookingResponses = append(bookingResponses, model.BookingResponse{
			ID:          booking.ID,
			PatientName: booking.User.Name,
			DoctorName:  user.Name,
			ServiceName: booking.Service.Name,
			BookingDate: booking.BookingDate,
			BookingTime: booking.BookingTime,
			Status:      booking.Status,
			Notes:       booking.Notes,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": bookingResponses, "total_rows": pagination.TotalRows, "total_pages": pagination.TotalPages, "current_page": pagination.Page, "limit": pagination.Limit})
}

func (bc *BookingController) UpdateBooking(c *gin.Context) {
	bookingId := c.Param("id")
	bookingIdUint, err := strconv.ParseUint(bookingId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var updateRequest model.UpdateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRole := c.MustGet("role").(string)

	updatedBooking, err := bc.BookingService.UpdateBooking(uint(bookingIdUint), model.Booking{
		Notes:       updateRequest.Notes,
		Status:      updateRequest.Status,
		BookingDate: updateRequest.BookingDate,
		BookingTime: updateRequest.BookingTime,
	}, userRole)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	doctor, err := bc.DoctorService.GetDoctorById(updatedBooking.DoctorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := bc.UserService.GetUserById(doctor.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookingResponse := model.BookingResponse{
		ID:          updatedBooking.ID,
		DoctorName:  user.Name,
		BookingDate: updatedBooking.BookingDate,
		BookingTime: updatedBooking.BookingTime,
		Status:      updatedBooking.Status,
		Notes:       updatedBooking.Notes,
	}

	c.JSON(http.StatusOK, gin.H{"booking": bookingResponse})
}

func (bc *BookingController) DeleteBooking(c *gin.Context) {
	bookingId := c.Param("id")
	bookingIdUint, err := strconv.ParseUint(bookingId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	userRole := c.MustGet("role").(string)
	userID := c.MustGet("userID").(uint)

	err = bc.BookingService.DeleteBooking(uint(bookingIdUint), userRole, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
}
