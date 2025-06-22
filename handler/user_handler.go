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

	if err := c.ShouldBindJSON(&input); err != nil {
		response := helper.APIresponse(http.StatusUnprocessableEntity, gin.H{"errors": helper.FormatValidationError(err)})
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isAvailable, err := h.userService.IsUsernameAvailability(input.Username)
	if err != nil || !isAvailable {
		status := http.StatusConflict
		if err != nil {
			status = http.StatusInternalServerError
		}
		response := helper.APIresponse(status, err.Error())
		c.JSON(status, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIresponse(http.StatusUnprocessableEntity, err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, newUser))
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

	// var idCabang int
	// var namaCabang string

	// if loggedinUser.IDCabang != nil {
	// 	fmt.Println("Get cabang for ID:", *loggedinUser.IDCabang)

	// 	cabang, err := h.cabangService.GetByID(*loggedinUser.IDCabang)
	// 	if err != nil {
	// 		fmt.Println("Error get cabang:", err.Error())
	// 	} else if cabang == nil {
	// 		fmt.Println("Cabang not found")
	// 	} else {
	// 		idCabang = cabang.IDCabang
	// 		namaCabang = cabang.NamaCabang
	// 	}
	// }

	response := gin.H{
		"message":     "Login successful!",
		"token":       token,
		"id_cabang":   loggedinUser.IDCabang,
		"nama_cabang": loggedinUser.Cabangs.NamaCabang,
		"access_code": loggedinUser.AccessCode,
	}

	val := helper.APIresponse(http.StatusOK, response)
	c.JSON(http.StatusOK, val)
}
