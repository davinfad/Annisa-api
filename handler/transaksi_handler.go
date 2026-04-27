package handler

import (
	"annisa-api/helper"
	"annisa-api/models"
	"annisa-api/service"
	"database/sql"
	"net/http"
	"strconv"
	"time"

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

type TransaksiDetailResponse struct {
	Transaksi *models.Transaksi            `json:"transaksi"`
	Items     []models.ItemTransaksiDetail `json:"items"`
}

func (h *HandlerTransaksi) GetTransaksiByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	transaksi, err := h.Service.GetTransaksiByID(id)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if transaksi == nil {
		response := helper.APIresponse(http.StatusNotFound, gin.H{"error": "Transaksi not found"})
		c.JSON(http.StatusNotFound, response)
		return
	}

	items, err := h.Service.GetItemTransaksiByTransaksiID(id)
	if err != nil {
		response := helper.APIresponse(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	result := TransaksiDetailResponse{
		Transaksi: transaksi,
		Items:     items,
	}
	response := helper.APIresponse(http.StatusOK, result)
	c.JSON(http.StatusOK, response)
}

// GET /transaksi/:id_cabang?from=2026-04-27&to=2026-04-27
// GET /transaksi/:id_cabang?from=2026-04-01&to=2026-04-30
func (h *HandlerTransaksi) GetTransaksiByDateRange(c *gin.Context) {
	idCabang := toInt(c.Param("id_cabang"))

	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, "from and to query params are required"))
		return
	}

	layout := "2006-01-02"
	loc := time.FixedZone("WIB", 7*3600)

	from, err := time.ParseInLocation(layout, fromStr, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, "invalid from date"))
		return
	}

	to, err := time.ParseInLocation(layout, toStr, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIresponse(http.StatusBadRequest, "invalid to date"))
		return
	}

	// set to to end of day
	to = to.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	transaksis, err := h.Service.GetTransaksiByDateRange(idCabang, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIresponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIresponse(http.StatusOK, transaksis))
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
