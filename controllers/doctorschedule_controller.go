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

type DoctorScheduleController struct {
	DoctorScheduleService services.DoctorScheduleService
}

func (dsc *DoctorScheduleController) CreateDoctorSchedule(c *gin.Context) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	time.Local = loc
	var doctorScheduleRequest model.DoctorScheduleRequest

	if err := c.ShouldBindJSON(&doctorScheduleRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.ParseInLocation("2006-01-02", doctorScheduleRequest.Date, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
	startTime, err := time.ParseInLocation("2006-01-02 15:04", doctorScheduleRequest.Date+" "+doctorScheduleRequest.StartTime, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
		return
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04", doctorScheduleRequest.Date+" "+doctorScheduleRequest.EndTime, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format"})
		return
	}

	userID := c.MustGet("userID").(uint)

	schedule := model.DoctorSchedule{
		DoctorId:  doctorScheduleRequest.DoctorID,
		ServiceId: doctorScheduleRequest.ServiceID,
		Date:      date,
		StartTime: startTime,
		EndTime:   endTime,
		CreatedBy: userID,
	}
	createdDoctorSchedule, err := dsc.DoctorScheduleService.CreateDoctorSchedule(schedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doctor schedule created successfully", "doctorSchedule": createdDoctorSchedule})
}

func (dsc *DoctorScheduleController) GetAllDoctorSchedules(c *gin.Context) {
	paginator, err := utils.Pagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctorSchedules, err := dsc.DoctorScheduleService.GetAllDoctorSchedules(paginator.Limit, paginator.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var doctorScheduleResponses []model.DoctorScheduleResponse
	for _, doctorSchedule := range doctorSchedules {
		doctorScheduleResponses = append(doctorScheduleResponses, model.DoctorScheduleResponse{
			ID:        doctorSchedule.ID,
			DoctorID:  doctorSchedule.Doctor.ID,
			Date:      doctorSchedule.Date,
			StartTime: doctorSchedule.StartTime,
			EndTime:   doctorSchedule.EndTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"doctorSchedules": doctorScheduleResponses})
}

func (dsc *DoctorScheduleController) GetDoctorScheduleById(c *gin.Context) {
	doctorScheduleId := c.Param("id")
	doctorScheduleIdUint, err := strconv.ParseUint(doctorScheduleId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor schedule ID"})
		return
	}

	doctorSchedule, err := dsc.DoctorScheduleService.GetDoctorScheduleById(uint(doctorScheduleIdUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"doctorSchedule": doctorSchedule})
}

func (dsc *DoctorScheduleController) UpdateDoctorSchedule(c *gin.Context) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	time.Local = loc
	doctorScheduleId := c.Param("id")
	doctorScheduleIdUint, err := strconv.ParseUint(doctorScheduleId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor schedule ID"})
		return
	}

	var updateRequest model.DoctorScheduleRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.ParseInLocation("2006-01-02", updateRequest.Date, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
	startTime, err := time.ParseInLocation("2006-01-02 15:04", updateRequest.Date+" "+updateRequest.StartTime, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
		return
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04", updateRequest.Date+" "+updateRequest.EndTime, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format"})
		return
	}

	doctorSchedule, err := dsc.DoctorScheduleService.UpdateDoctorSchedule(uint(doctorScheduleIdUint), model.DoctorSchedule{
		Date:      date,
		StartTime: startTime,
		EndTime:   endTime,
		UpdatedBy: c.MustGet("userID").(uint),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doctor schedule updated successfully", "doctorSchedule": doctorSchedule})
}

func (dsc *DoctorScheduleController) DeleteDoctorSchedule(c *gin.Context) {
	doctorScheduleId := c.Param("id")
	userID := c.MustGet("userID").(uint)
	doctorScheduleIdUint, err := strconv.ParseUint(doctorScheduleId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor schedule ID"})
		return
	}

	err = dsc.DoctorScheduleService.DeleteDoctorSchedule(uint(doctorScheduleIdUint), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doctor schedule deleted successfully"})
}
