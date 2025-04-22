package controllers

import (
	"booking-klinik/model"
	"booking-klinik/services"
	"booking-klinik/utils"
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

	createdDoctor, err := dc.DoctorService.CreateDoctor(doctor)
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
	paginator, err := utils.Pagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctors, err := dc.DoctorService.GetAllDoctors(paginator.Limit, paginator.Offset)
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

	c.JSON(http.StatusOK, gin.H{"doctors": doctorResponses})
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

	err = dc.DoctorService.DeleteDoctor(uint(doctorIdUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doctor deleted successfully"})
}
