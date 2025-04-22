package controllers

import (
	"booking-klinik/model"
	"booking-klinik/services"
	"booking-klinik/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DoctorScheduleController struct {
	DoctorScheduleService services.DoctorScheduleService
}

func (dsc *DoctorScheduleController) CreateDoctorSchedule(c *gin.Context) {
	var doctorSchedule model.DoctorSchedule

	if err := c.ShouldBindJSON(&doctorSchedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdDoctorSchedule, err := dsc.DoctorScheduleService.CreateDoctorSchedule(doctorSchedule)
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
			Day:       doctorSchedule.Day,
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

	doctorSchedule, err := dsc.DoctorScheduleService.UpdateDoctorSchedule(uint(doctorScheduleIdUint), model.DoctorSchedule{
		Day:       updateRequest.Day,
		StartTime: updateRequest.StartTime,
		EndTime:   updateRequest.EndTime,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doctor schedule updated successfully", "doctorSchedule": doctorSchedule})
}

func (dsc *DoctorScheduleController) DeleteDoctorSchedule(c *gin.Context) {
	doctorScheduleId := c.Param("id")
	doctorScheduleIdUint, err := strconv.ParseUint(doctorScheduleId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor schedule ID"})
		return
	}

	err = dsc.DoctorScheduleService.DeleteDoctorSchedule(uint(doctorScheduleIdUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
