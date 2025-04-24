package controllers

import (
	"booking-klinik/model"
	"booking-klinik/services"
	"booking-klinik/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceController struct {
	ServiceService services.ServiceService
}

func (sc *ServiceController) CreateService(c *gin.Context) {
	var service model.Service

	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Default isActive = true

	userID := c.MustGet("userID").(uint)
	service.CreatedBy = userID
	service.IsActive = true
	createdService, err := sc.ServiceService.CreateService(service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service created successfully", "service": createdService})
}

func (sc *ServiceController) GetAllServices(c *gin.Context) {
	paginator, err := utils.Pagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	services, err := sc.ServiceService.GetAllServices(paginator.Limit, paginator.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var serviceResponses []model.ServiceResponse
	for _, service := range services {
		serviceResponses = append(serviceResponses, model.ServiceResponse{
			Name:            service.Name,
			Description:     service.Description,
			Price:           service.Price,
			DurationMinutes: service.DurationMinutes,
			IsActive:        service.IsActive,
		})
	}

	c.JSON(http.StatusOK, gin.H{"services": serviceResponses})
}

func (sc *ServiceController) GetServiceById(c *gin.Context) {
	serviceId := c.Param("id")
	serviceIdUint, err := strconv.ParseUint(serviceId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	service, err := sc.ServiceService.GetServiceById(uint(serviceIdUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var serviceResponse model.ServiceResponse
	serviceResponse.Name = service.Name
	serviceResponse.Description = service.Description
	serviceResponse.Price = service.Price
	serviceResponse.DurationMinutes = service.DurationMinutes
	serviceResponse.IsActive = service.IsActive

	c.JSON(http.StatusOK, gin.H{"service": serviceResponse})
}

func (sc *ServiceController) UpdateService(c *gin.Context) {
	serviceId := c.Param("id")

	serviceIdUint, err := strconv.ParseUint(serviceId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	var service model.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)

	service.UpdatedBy = userID

	updatedService, err := sc.ServiceService.UpdateService(uint(serviceIdUint), service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service updated successfully", "service": updatedService})
}

func (sc *ServiceController) DeleteService(c *gin.Context) {
	serviceId := c.Param("id")
	serviceIdUint, err := strconv.ParseUint(serviceId, 10, 32)
	userID := c.MustGet("userID").(uint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	err = sc.ServiceService.DeleteService(uint(serviceIdUint), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}
