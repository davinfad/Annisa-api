package handler

import (
	"annisa-api/helper"
	"annisa-api/models"
	"annisa-api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type layananHandler struct {
	layananService service.ServiceLayanan
}

func NewLayananHandler(layananService service.ServiceLayanan) *layananHandler {
	return &layananHandler{layananService}
}

func (h *layananHandler) Create(c *gin.Context) {
	var input *models.CreateLayananDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIresponse(http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	val, err := h.layananService.Create(input)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)

}

func (h *layananHandler) Update(c *gin.Context) {
	var input *models.CreateLayananDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIresponse(http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	param := c.Param("id")
	params, _ := strconv.Atoi(param)

	val, err := h.layananService.Update(params, input)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)

}

func (h *layananHandler) GetByID(c *gin.Context) {
	param := c.Param("id")
	params, _ := strconv.Atoi(param)

	val, err := h.layananService.GetByID(params)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}

func (h *layananHandler) Delete(c *gin.Context) {
	param := c.Param("id")
	params, _ := strconv.Atoi(param)

	val, err := h.layananService.Delete(params)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}

func (h *layananHandler) GetAll(c *gin.Context) {
	val, err := h.layananService.GetAll()
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}
