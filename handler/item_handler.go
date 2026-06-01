package handler

import (
	"annisa-api/helper"
	"annisa-api/models"
	"annisa-api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type itemHandler struct {
	itemService service.ServiceItem
}

func NewItemHandler(itemService service.ServiceItem) *itemHandler {
	return &itemHandler{itemService}
}

func (h *itemHandler) Create(c *gin.Context) {
	var input *models.CreateItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		c.JSON(http.StatusUnprocessableEntity, helper.APIresponse(http.StatusUnprocessableEntity, errorMessage))
		return
	}

	val, err := h.itemService.Create(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}

func (h *itemHandler) Update(c *gin.Context) {
	var input *models.CreateItemDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		c.JSON(http.StatusUnprocessableEntity, helper.APIresponse(http.StatusUnprocessableEntity, errorMessage))
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	val, err := h.itemService.Update(id, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}

func (h *itemHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	val, err := h.itemService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}

func (h *itemHandler) GetAll(c *gin.Context) {
	val, err := h.itemService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIresponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}

func (h *itemHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	val, err := h.itemService.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}
