package models

import "time"

type Layanan struct {
	IDLayanan           int       `json:"id_layanan"`
	NamaLayanan         string    `json:"nama_layanan"`
	PersenKomisi        float64   `json:"persen_komisi"`
	PersenKomisiLuarJam float64   `json:"persen_komisi_luarjam"`
	Kategori            string    `json:"kategori"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type CreateLayananDTO struct {
	NamaLayanan         string  `json:"nama_layanan" binding:"required"`
	PersenKomisi        float64 `json:"persen_komisi" binding:"required"`
	PersenKomisiLuarJam float64 `json:"persen_komisi_luarjam" binding:"required"`
	Kategori            string  `json:"kategori" binding:"required"`
}
