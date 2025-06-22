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

	if err := c.ShouldBindJSON(&input); err != nil {
		response := helper.APIresponse(http.StatusUnprocessableEntity, gin.H{"errors": helper.FormatValidationError(err)})
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	val, err := h.cabangService.Create(input)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}
