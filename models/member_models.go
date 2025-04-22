package models

import "time"

type Member struct {
	IDMember       int       `json:"id_member"`
	NomorPelanggan string    `json:"nomor_pelanggan"`
	NamaMember     string    `json:"nama_member"`
	NomorTelepon   string    `json:"nomor_telepon"`
	Alamat         string    `json:"alamat"`
	TanggalLahir   time.Time `json:"tanggal_lahir"`
	TanggalDaftar  time.Time `json:"tanggal_daftar"`
	IDCabang       int       `json:"id_cabang"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateMemberDTO struct {
	NomorPelanggan string `json:"nomor_pelanggan" binding:"required"`
	NamaMember     string `json:"nama_member" binding:"required"`
	NomorTelepon   string `json:"nomor_telepon" binding:"required"`
	Alamat         string `json:"alamat" binding:"required"`
	TanggalLahir   string `json:"tanggal_lahir" binding:"required"`
	TanggalDaftar  string `json:"tanggal_daftar" binding:"required"`
	IDCabang       int    `json:"id_cabang" binding:"required"`
}
