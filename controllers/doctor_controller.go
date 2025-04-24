package controllers

import (
	"booking-klinik/model"
	"booking-klinik/services"
	"booking-klinik/utils"
	"fmt"
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
	fmt.Println("test")
	paginator, err := utils.Pagination(c)
	fmt.Println("test2")
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(paginator.Limit)
	fmt.Println(paginator.Offset)
	fmt.Println(paginator.Page)
	fmt.Println("test2")
	doctors, err := dc.DoctorService.GetAllDoctors(paginator.Limit, paginator.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	/*
		paginator.TotalRows = totalRows
		paginator.TotalPages = (totalRows + int64(paginator.Limit) - 1) / int64(paginator.Limit)*/
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
		"total_rows":   paginator.TotalRows,
		"total_pages":  paginator.TotalPages,
		"current_page": paginator.Page,
		"limit":        paginator.Limit,
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
