package handler

import (
	"annisa-api/helper"
	"annisa-api/models"
	"annisa-api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type karyawanHandler struct {
	karyawanService service.ServiceKaryawan
}

func NewKaryawanHandler(karyawanService service.ServiceKaryawan) *karyawanHandler {
	return &karyawanHandler{karyawanService}
}

func (h *karyawanHandler) Create(c *gin.Context) {
	var input *models.CreateKaryawanDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIresponse(http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	val, err := h.karyawanService.Create(input)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}

func (h *karyawanHandler) Update(c *gin.Context) {

	param := c.Param("id")
	params, _ := strconv.Atoi(param)

	var input *models.CreateKaryawanDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIresponse(http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	val, err := h.karyawanService.Update(params, input)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}

func (h *karyawanHandler) GetByID(c *gin.Context) {
	param := c.Param("id")
	params, _ := strconv.Atoi(param)

	val, err := h.karyawanService.GetByID(params)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}

func (h *karyawanHandler) GetByIDCabang(c *gin.Context) {
	param := c.Param("id_cabang")
	params, _ := strconv.Atoi(param)

	val, err := h.karyawanService.GetByIDCabang(params)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}

func (h *karyawanHandler) Delete(c *gin.Context) {
	param := c.Param("id")
	params, _ := strconv.Atoi(param)

	val, err := h.karyawanService.Delete(params)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}
