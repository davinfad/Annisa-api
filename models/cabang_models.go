package models

import "time"

type Cabang struct {
	IDCabang   int       `json:"id_cabang"`
	NamaCabang string    `json:"nama_cabang"`
	KodeCabang string    `json:"kode_cabang"`
	JamBuka    time.Time `json:"jam_buka"`
	JamTutup   time.Time `json:"jam_tutup"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CabangDTO struct {
	NamaCabang string `json:"nama_cabang" binding:"required"`
	KodeCabang string `json:"kode_cabang" binding:"required"`
	JamBuka    string `json:"jam_buka" binding:"required"`
	JamTutup   string `json:"jam_tutup" binding:"required"`
}
