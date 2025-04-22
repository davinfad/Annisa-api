package handler

import (
	"annisa-api/helper"
	"annisa-api/models"
	"annisa-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type cabangHandler struct {
	cabangService service.ServiceCabang
}

func NewCabangHandler(cabangService service.ServiceCabang) *cabangHandler {
	return &cabangHandler{cabangService}
}

func (h *cabangHandler) Create(c *gin.Context) {
	var input *models.CabangDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIresponse(http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	val, err := h.cabangService.Create(input)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}
