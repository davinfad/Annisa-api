package handler

import (
	"annisa-api/helper"
	"annisa-api/models"
	"annisa-api/service"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerTransaksi struct {
	DB      *sql.DB
	Service service.ServiceTransaksi
}

func NewHandlerTransaksi(db *sql.DB, service service.ServiceTransaksi) *HandlerTransaksi {
	return &HandlerTransaksi{
		DB:      db,
		Service: service,
	}
}

func (h *HandlerTransaksi) GetTotalMoneyByDateAndCabang(c *gin.Context) {
	date := c.Param("date")
	idCabangStr := c.Param("id_cabang")

	idCabang, err := strconv.Atoi(idCabangStr)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := h.Service.GetTotalMoneyByDateAndCabang(date, idCabang)
	if err != nil {
		response := helper.APIresponse(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, result)
	c.JSON(http.StatusOK, response)
}

func (h *HandlerTransaksi) GetTotalMoneyByMonthAndYear(c *gin.Context) {
	monthStr := c.Param("month")
	yearStr := c.Param("year")
	idCabangStr := c.Param("id_cabang")

	month, err1 := strconv.Atoi(monthStr)
	year, err2 := strconv.Atoi(yearStr)
	idCabang, err3 := strconv.Atoi(idCabangStr)
	if err1 != nil || err2 != nil || err3 != nil {
		response := helper.APIresponse(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := h.Service.GetTotalMoneyByMonthAndYear(month, year, idCabang)
	if err != nil {
		response := helper.APIresponse(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, result)
	c.JSON(http.StatusOK, response)
}

func (h *HandlerTransaksi) GetTransaksiByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
	}

	transaksi, err := h.Service.GetTransaksiByID(id)
	if err != nil {
		response := helper.APIresponse(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if transaksi == nil {
		response := helper.APIresponse(http.StatusNotFound, gin.H{"error": "Transaksi not found"})
		c.JSON(http.StatusNotFound, response)
		return
	}
	response := helper.APIresponse(http.StatusOK, transaksi)
	c.JSON(http.StatusOK, response)
}

func (h *HandlerTransaksi) GetTransaksiByDateAndCabang(c *gin.Context) {
	date := c.Param("date")
	idCabang := toInt(c.Param("id_cabang"))

	transaksis, err := h.Service.GetTransaksiByDateAndCabang(date, idCabang)
	if err != nil {
		response := helper.APIresponse(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := helper.APIresponse(http.StatusOK, transaksis)
	c.JSON(http.StatusOK, response)
}

func (h *HandlerTransaksi) GetMonthlyTransaksiByCabang(c *gin.Context) {
	month := toInt(c.Param("month"))
	year := toInt(c.Param("year"))
	idCabang := toInt(c.Param("id_cabang"))

	transaksis, err := h.Service.GetMonthlyTransaksiByCabang(month, year, idCabang)
	if err != nil {
		response := helper.APIresponse(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := helper.APIresponse(http.StatusOK, transaksis)
	c.JSON(http.StatusOK, response)
}

func (h *HandlerTransaksi) GetDraftTransaksiByCabang(c *gin.Context) {
	idCabang := toInt(c.Param("id_cabang"))

	transaksis, err := h.Service.GetDraftTransaksiByCabang(idCabang)
	if err != nil {
		response := helper.APIresponse(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := helper.APIresponse(http.StatusOK, transaksis)
	c.JSON(http.StatusOK, response)
}

func (h *HandlerTransaksi) AddTransaksi(c *gin.Context) {
	var req models.TransaksiRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if req.TotalHarga == 0 || req.MetodePembayaran == "" || req.IDCabang == nil || len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields!"})
		return
	}

	status := 0
	if req.IsDraft {
		status = 1
	}

	tx, err := h.DB.Begin()
	if err != nil {
		response := helper.APIresponse(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	transaksi, err := h.Service.CreateTransaksi(tx, req, status)
	if err != nil {
		tx.Rollback()
		response := helper.APIresponse(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if status == 0 {
		err := h.Service.UpdateKomisiKaryawan(tx, req.Items, transaksi.CreatedAt, transaksi.IDCabang)
		if err != nil {
			tx.Rollback()
			response := helper.APIresponse(http.StatusInternalServerError, err.Error())
			c.JSON(http.StatusInternalServerError, response)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		response := helper.APIresponse(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, transaksi)
	c.JSON(http.StatusOK, response)
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func (h *HandlerTransaksi) DeleteTransaksi(c *gin.Context) {
	idParam := c.Param("id_transaksi")
	idTransaksi, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id_transaksi"})
		return
	}

	err = h.Service.DeleteTransaksi(c.Request.Context(), idTransaksi)
	if err != nil {
		if err.Error() == "transaction not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction and its items deleted successfully"})
}
