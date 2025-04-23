package models

import "time"

type User struct {
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	AccessCode string    `json:"access_code"`
	IDCabang   *int      `json:"id_cabang"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserRegisterDTO struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	AccessCode string `json:"access_code" binding:"required"`
	IDCabang   *int   `json:"id_cabang"`
}

type UserLoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
