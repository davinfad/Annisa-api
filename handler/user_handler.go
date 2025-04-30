package handler

import (
	"annisa-api/auth"
	"annisa-api/helper"
	"annisa-api/models"
	"annisa-api/service"
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

	cabang, err := h.cabangService.GetByID(*loggedinUser.IDCabang)
	if err != nil {
		response := helper.APIresponse(http.StatusInternalServerError, "Failed to fetch cabang data")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	responseData := gin.H{
		"message":     "Login successful!",
		"token":       token,
		"id_cabang":   cabang.IDCabang,
		"nama_cabang": cabang.NamaCabang,
		"acces_code":  loggedinUser.AccessCode,
	}
	response := helper.APIresponse(http.StatusOK, responseData)
	c.JSON(http.StatusOK, response)
}
