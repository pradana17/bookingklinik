package controllers

import (
	"booking-klinik/model"
	"booking-klinik/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DoctorController struct {
	DoctorService services.DoctorServices
}

func (dc *DoctorController) CreateDoctor(c *gin.Context) {
	var doctor model.Doctor

	if err := c.ShouldBindJSON(&doctor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}
	doctor.CreatedBy = userID.(uint)

	createdDoctor, err := dc.DoctorService.CreateDoctor(&doctor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	doctorResponse := model.DoctorResponse{
		ID:             createdDoctor.ID,
		Name:           createdDoctor.User.Name,
		Specialization: createdDoctor.Specialization,
	}
	c.JSON(http.StatusOK, gin.H{"message": "Doctor created successfully", "doctor": doctorResponse})
}

func (dc *DoctorController) GetAllDoctors(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	doctors, pagination, err := dc.DoctorService.GetAllDoctors(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var doctorResponses []model.DoctorResponse
	for _, doctor := range doctors {
		doctorResponses = append(doctorResponses, model.DoctorResponse{
			ID:             doctor.ID,
			Name:           doctor.User.Name,
			Specialization: doctor.Specialization,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":         doctorResponses,
		"total_rows":   pagination.TotalRows,
		"total_pages":  pagination.TotalPages,
		"current_page": pagination.Page,
		"limit":        pagination.Limit,
	})
}

func (dc *DoctorController) GetDoctorById(c *gin.Context) {
	doctorId := c.Param("id")
	doctorIdUint, err := strconv.ParseUint(doctorId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	doctor, err := dc.DoctorService.GetDoctorById(uint(doctorIdUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	doctorResponse := model.DoctorResponse{
		ID:             doctor.ID,
		Name:           doctor.User.Name,
		Specialization: doctor.Specialization,
	}

	c.JSON(http.StatusOK, gin.H{"doctor": doctorResponse})
}

func (dc *DoctorController) UpdateDoctor(c *gin.Context) {
	doctorId := c.Param("id")

	doctorIdUint, err := strconv.ParseUint(doctorId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	var doctor model.Doctor
	if err := c.ShouldBindJSON(&doctor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)
	doctor.UpdatedBy = userID

	updatedDoctor, err := dc.DoctorService.UpdateDoctor(uint(doctorIdUint), doctor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	doctorResponse := model.DoctorResponse{
		ID:             updatedDoctor.ID,
		Name:           updatedDoctor.User.Name,
		Specialization: updatedDoctor.Specialization,
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doctor updated successfully", "doctor": doctorResponse})
}

func (dc *DoctorController) DeleteDoctor(c *gin.Context) {
	doctorId := c.Param("id")
	doctorIdUint, err := strconv.ParseUint(doctorId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	userID := c.MustGet("userID").(uint)
	err = dc.DoctorService.DeleteDoctor(uint(doctorIdUint), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doctor deleted successfully"})
}
