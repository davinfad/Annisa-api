package handler

import (
	"annisa-api/helper"
	"annisa-api/models"
	"annisa-api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type inventoryHandler struct {
	inventoryService service.ServiceInventory
}

func NewInventoryHandler(inventoryService service.ServiceInventory) *inventoryHandler {
	return &inventoryHandler{inventoryService}
}

func (h *inventoryHandler) Create(c *gin.Context) {
	var input *models.CreateInventoryDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		c.JSON(http.StatusUnprocessableEntity, helper.APIresponse(http.StatusUnprocessableEntity, errorMessage))
		return
	}

	val, err := h.inventoryService.Create(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}

func (h *inventoryHandler) Update(c *gin.Context) {
	var input *models.CreateInventoryDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		c.JSON(http.StatusUnprocessableEntity, helper.APIresponse(http.StatusUnprocessableEntity, errorMessage))
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	val, err := h.inventoryService.Update(id, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}

func (h *inventoryHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	val, err := h.inventoryService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}

func (h *inventoryHandler) GetByCabang(c *gin.Context) {
	idCabang, _ := strconv.Atoi(c.Param("id_cabang"))

	val, err := h.inventoryService.GetByCabang(idCabang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIresponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}

func (h *inventoryHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	val, err := h.inventoryService.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}

func (h *inventoryHandler) AdjustStok(c *gin.Context) {
	var input *models.AdjustStokDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		c.JSON(http.StatusUnprocessableEntity, helper.APIresponse(http.StatusUnprocessableEntity, errorMessage))
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	val, err := h.inventoryService.AdjustStok(id, input.Delta)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}
