package models

import "time"

type ItemTransaksi struct {
	IDItemTransaksi int       `json:"id_item_transaksi"`
	IDTransaksi     *int      `json:"id_transaksi"`
	IDKaryawan      *int      `json:"id_karyawan"`
	IDLayanan       *int      `json:"id_layanan"`
	Catatan         *string   `json:"catatan"`
	Harga           float64   `json:"harga"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ItemTransaksiDetail struct {
	IDItemTransaksi int       `json:"id_item_transaksi"`
	IDTransaksi     *int      `json:"id_transaksi"`
	IDKaryawan      *int      `json:"id_karyawan"`
	IDLayanan       *int      `json:"id_layanan"`
	Catatan         *string   `json:"catatan"`
	Harga           float64   `json:"harga"`
	CreatedAt       time.Time `json:"created_at"`
	NamaKaryawan    string    `json:"nama_karyawan"`
	NamaLayanan     string    `json:"nama_layanan"`
}
