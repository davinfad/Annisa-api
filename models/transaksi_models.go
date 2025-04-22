package models

import "time"

type Transaksi struct {
	IDTransaksi      int       `json:"id_transaksi"`
	IDCabang         *int      `json:"id_cabang"`         //
	IDMember         *int      `json:"id_member"`         //
	NamaPelanggan    string    `json:"nama_pelanggan"`    //
	NomorTelepon     string    `json:"nomor_telepon"`     //
	TotalHarga       float64   `json:"total_harga"`       //
	MetodePembayaran string    `json:"metode_pembayaran"` //
	Status           *int      `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type TransaksiRequest struct {
	NamaPelanggan    string          `json:"nama_pelanggan"`
	NomorTelepon     string          `json:"nomor_telepon"`
	TotalHarga       float64         `json:"total_harga"`
	MetodePembayaran string          `json:"metode_pembayaran"`
	IDMember         *int            `json:"id_member"`
	IDCabang         *int            `json:"id_cabang"`
	Items            []ItemTransaksi `json:"items"`
	IsDraft          bool            `json:"isDraft"`
}

type TotalMoneyResult struct {
	TotalMoney    float64 `json:"total_money"`
	TotalCash     float64 `json:"total_cash"`
	TotalTransfer float64 `json:"total_transfer"`
}
