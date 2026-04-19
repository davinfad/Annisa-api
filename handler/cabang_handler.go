package handler

import (
	"annisa-api/helper"
	"annisa-api/models"
	"annisa-api/service"
	"net/http"
	"strconv"

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
		c.JSON(http.StatusUnprocessableEntity, helper.APIresponse(http.StatusUnprocessableEntity, gin.H{"errors": helper.FormatValidationError(err)}))
		return
	}

	val, err := h.cabangService.Create(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, val))
}

func (h *cabangHandler) GetAll(c *gin.Context) {
	cabangs, err := h.cabangService.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, cabangs))
}

func (h *cabangHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, "invalid id"))
		return
	}

	cabang, err := h.cabangService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.APIresponse(http.StatusNotFound, err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, cabang))
}

func (h *cabangHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, "invalid id"))
		return
	}

	var input models.CabangWithUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, helper.APIresponse(http.StatusUnprocessableEntity, gin.H{"errors": helper.FormatValidationError(err)}))
		return
	}

	cabangDTO := &models.CabangDTO{
		NamaCabang: input.NamaCabang,
		KodeCabang: input.KodeCabang,
		JamBuka:    input.JamBuka,
		JamTutup:   input.JamTutup,
	}

	var userDTO *models.UpdateUserDTO
	if input.Username != "" {
		userDTO = &models.UpdateUserDTO{
			Username:   input.Username,
			Password:   input.Password,
			AccessCode: input.AccessCode,
		}
	}

	cabang, err := h.cabangService.Update(id, cabangDTO, userDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, cabang))
}

func (h *cabangHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, "invalid id"))
		return
	}

	if err := h.cabangService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, "cabang berhasil dihapus"))
}
