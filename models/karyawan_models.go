package models

import "time"

type Karyawan struct {
	IDKaryawan   int       `json:"id_karyawan"`
	NamaKaryawan string    `json:"nama_karyawan"`
	IDCabang     *int      `json:"id_cabang"`
	NomorTelepon *string   `json:"nomor_telepon"`
	Alamat       *string   `json:"alamat"`
	Komisi       float64   `json:"komisi"`
	KomisiHarian float64   `json:"komisi_harian"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateKaryawanDTO struct {
	NamaKaryawan string  `json:"nama_karyawan" binding:"required"`
	IDCabang     *int    `json:"id_cabang"`
	NomorTelepon *string `json:"nomor_telepon"`
	Alamat       *string `json:"alamat"`
	Komisi       float64 `json:"komisi"`
	KomisiHarian float64 `json:"komisi_harian"`
}
