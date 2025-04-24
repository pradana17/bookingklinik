package controllers

import (
	"booking-klinik/model"
	"booking-klinik/services"
	"booking-klinik/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	BookingService services.BookingService
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

	createdBooking, err := bc.BookingService.CreateBooking(newBooking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookingResponse := model.BookingResponse{
		ID:          createdBooking.ID,
		DoctorName:  createdBooking.User.Name,
		ServiceName: createdBooking.Service.Name,
		BookingDate: createdBooking.BookingDate,
		BookingTime: createdBooking.BookingTime,
		Status:      createdBooking.Status,
		Notes:       createdBooking.Notes,
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking created successfully", "booking": bookingResponse})
}

func (bc *BookingController) GetAllBookings(c *gin.Context) {
	paginator, err := utils.Pagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookings, pagination, err := bc.BookingService.GetAllBookings(paginator.Limit, paginator.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var bookingResponses []model.BookingResponse
	for _, booking := range bookings {
		bookingResponses = append(bookingResponses, model.BookingResponse{
			ID:          booking.ID,
			DoctorName:  booking.User.Name,
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
	bookingIdUint, err := strconv.ParseUint(bookingId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	booking, err := bc.BookingService.GetBookingById(uint(bookingIdUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookingResponse := model.BookingResponse{
		ID:          booking.ID,
		DoctorName:  booking.User.Name,
		ServiceName: booking.Service.Name,
		BookingDate: booking.BookingDate,
		BookingTime: booking.BookingTime,
		Status:      booking.Status,
		Notes:       booking.Notes,
	}

	c.JSON(http.StatusOK, gin.H{"booking": bookingResponse})
}

func (bc *BookingController) GetBookingsByUserId(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	paginator, err := utils.Pagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookings, err := bc.BookingService.GetBookingsByUserId(userID, paginator.Limit, paginator.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var bookingResponses []model.BookingResponse
	for _, booking := range bookings {
		bookingResponses = append(bookingResponses, model.BookingResponse{
			ID:          booking.ID,
			DoctorName:  booking.User.Name,
			ServiceName: booking.Service.Name,
			BookingDate: booking.BookingDate,
			BookingTime: booking.BookingTime,
			Status:      booking.Status,
			Notes:       booking.Notes,
		})
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookingResponses})
}

func (bc *BookingController) GetBookingsByDoctorId(c *gin.Context) {
	doctorId := c.Param("id")
	doctorIdUint, err := strconv.ParseUint(doctorId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	paginator, err := utils.Pagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookings, err := bc.BookingService.GetBookingsByDoctorId(uint(doctorIdUint), paginator.Limit, paginator.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var bookingResponses []model.BookingResponse
	for _, booking := range bookings {
		bookingResponses = append(bookingResponses, model.BookingResponse{
			ID:          booking.ID,
			DoctorName:  booking.User.Name,
			ServiceName: booking.Service.Name,
			BookingDate: booking.BookingDate,
			BookingTime: booking.BookingTime,
			Status:      booking.Status,
			Notes:       booking.Notes,
		})
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookingResponses})
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

	userRole := c.MustGet("userRole").(string)

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

	bookingResponse := model.BookingResponse{
		ID:          updatedBooking.ID,
		DoctorName:  updatedBooking.User.Name,
		ServiceName: updatedBooking.Service.Name,
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

	userRole := c.MustGet("userRole").(string)
	userID := c.MustGet("userID").(uint)

	err = bc.BookingService.DeleteBooking(uint(bookingIdUint), userRole, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
}
