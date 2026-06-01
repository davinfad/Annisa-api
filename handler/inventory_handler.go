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

// GET /inventory/cabang/:id_cabang -> full catalog with this branch's stock.
func (h *inventoryHandler) GetByCabang(c *gin.Context) {
	idCabang, _ := strconv.Atoi(c.Param("id_cabang"))

	val, err := h.inventoryService.GetByCabang(idCabang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIresponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}

// PATCH /inventory/cabang/:id_cabang/item/:id_item/stok -> relative change.
func (h *inventoryHandler) AdjustStok(c *gin.Context) {
	var input *models.AdjustStokDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		c.JSON(http.StatusUnprocessableEntity, helper.APIresponse(http.StatusUnprocessableEntity, errorMessage))
		return
	}

	idCabang, _ := strconv.Atoi(c.Param("id_cabang"))
	idItem, _ := strconv.Atoi(c.Param("id_item"))

	val, err := h.inventoryService.AdjustStok(idCabang, idItem, input.Delta)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}

// PUT /inventory/cabang/:id_cabang/item/:id_item -> absolute set.
func (h *inventoryHandler) SetStok(c *gin.Context) {
	var input *models.SetStokDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		c.JSON(http.StatusUnprocessableEntity, helper.APIresponse(http.StatusUnprocessableEntity, errorMessage))
		return
	}

	idCabang, _ := strconv.Atoi(c.Param("id_cabang"))
	idItem, _ := strconv.Atoi(c.Param("id_item"))

	val, err := h.inventoryService.SetStok(idCabang, idItem, *input.Stok)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}
