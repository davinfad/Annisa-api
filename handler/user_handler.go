package handler

import (
	"annisa-api/auth"
	"annisa-api/helper"
	"annisa-api/models"
	"annisa-api/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService   service.ServiceUser
	cabangService service.ServiceCabang
	authService   auth.UserAuthService
}

func NewUserHandler(userService service.ServiceUser, cabangService service.ServiceCabang, authService auth.UserAuthService) *userHandler {
	return &userHandler{userService, cabangService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input models.UserRegisterDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIresponse(http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsUsernameAvailability(input.Username)
	if err != nil {
		response := helper.APIresponse(http.StatusInternalServerError, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if !isEmailAvailable {
		response := helper.APIresponse(http.StatusConflict, err.Error())
		c.JSON(http.StatusConflict, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIresponse(http.StatusUnprocessableEntity, err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, newUser)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input models.UserLoginDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIresponse(http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.LoginUser(input)
	if err != nil {
		response := helper.APIresponse(http.StatusUnprocessableEntity, err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	token, err := h.authService.GenerateToken(loggedinUser.Username)
	if err != nil {
		response := helper.APIresponse(http.StatusUnprocessableEntity, err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var idCabang int
	var namaCabang string

	if loggedinUser.IDCabang != nil {
		fmt.Println("Get cabang for ID:", *loggedinUser.IDCabang)

		cabang, err := h.cabangService.GetByID(*loggedinUser.IDCabang)
		if err != nil {
			fmt.Println("Error get cabang:", err.Error())
		} else if cabang == nil {
			fmt.Println("Cabang not found")
		} else {
			idCabang = cabang.IDCabang
			namaCabang = cabang.NamaCabang
		}
	}

	response := gin.H{
		"message":     "Login successful!",
		"token":       token,
		"id_cabang":   idCabang,
		"nama_cabang": namaCabang,
		"val":         loggedinUser.AccessCode,
	}

	val := helper.APIresponse(http.StatusOK, response)
	c.JSON(http.StatusOK, val)
}
