package models

import "time"

type User struct {
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	AccessCode string    `json:"access_code"`
	IDCabang   *int      `json:"id_cabang"`
	Cabangs    Cabang    `json:"cabang"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserRegisterDTO struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	AccessCode string `json:"access_code" binding:"required"`
	IDCabang   *int   `json:"id_cabang"`
	CabangName string `json:"cabang_name"`
	KodeCabang string `json:"kode_cabang"`
	JamBuka    string `json:"jam_buka"`
	JamTutup   string `json:"jam_tutup"`
}

type UserLoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
