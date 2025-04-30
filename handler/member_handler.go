package handler

import (
	"annisa-api/helper"
	"annisa-api/models"
	"annisa-api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type memberHandler struct {
	memberService service.ServiceMember
}

func NewMemberHandler(memberService service.ServiceMember) *memberHandler {
	return &memberHandler{memberService}
}

func (h *memberHandler) Create(c *gin.Context) {
	var input *models.CreateMemberDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIresponse(http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	val, err := h.memberService.Create(input)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}

func (h *memberHandler) Update(c *gin.Context) {

	param := c.Param("id")
	params, _ := strconv.Atoi(param)

	var input *models.CreateMemberDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIresponse(http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	val, err := h.memberService.Update(params, input)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}

func (h *memberHandler) GetByID(c *gin.Context) {
	param := c.Param("id")
	params, _ := strconv.Atoi(param)

	val, err := h.memberService.GetByID(params)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}

func (h *memberHandler) GetAll(c *gin.Context) {
	val, err := h.memberService.GetAll()
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, val)
}

func (h *memberHandler) GetMemberByIDCabang(c *gin.Context) {
	param := c.Param("id_cabang")
	params, _ := strconv.Atoi(param)

	val, err := h.memberService.GetMemberByCabangID(params)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}

func (h *memberHandler) Delete(c *gin.Context) {
	param := c.Param("id")
	params, _ := strconv.Atoi(param)

	val, err := h.memberService.Delete(params)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, val)
	c.JSON(http.StatusOK, response)
}
